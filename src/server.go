package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/grandcat/zeroconf"

	"github.com/gorilla/mux"
	"github.com/kardianos/service"
)

// Server defines the web server
type Server struct {
	PortNo   int                  // Port number the server will listen on
	WaitTime int                  // Duration in secs to wait for replies when discovering services
	Debug    bool                 // Indicates whether the server is running in debug
	Config   *Config              // Configuration settings
	exit     chan struct{}        // Exit flag
	shutdown chan struct{}        // Shutdown complete flag
	http     *http.Server         // HTTP server
	router   *mux.Router          // HTTP router
	regList  map[string]*ZCServer // Zeroconf registration server list
	regLock  sync.Mutex           // Mutex lock for appending and removing items from regList
	hostName string               // HostName of computer
}

// AddController adds the specified web service controller to the Router
func (s *Server) addController(c Controller) {
	c.AddController(s.router, s)
}

// Start is called when the service is starting
func (s *Server) Start(v service.Service) error {
	s.logInfo("Service starting")

	s.regList = make(map[string]*ZCServer)

	// Make sure the working directory is the same as the application exe
	ap, err := os.Executable()
	if err != nil {
		s.logError("Error getting the executable path.", err.Error())
	} else {
		wd, err := os.Getwd()
		if err != nil {
			s.logError("Error getting current working directory.", err.Error())
		} else {
			ad := filepath.Dir(ap)
			s.logInfo("Current application path is", ad)
			if ad != wd {
				if err := os.Chdir(ad); err != nil {
					s.logError("Error chaning working directory.", err.Error())
				}
			}
		}
	}

	// Create a channel that will be used to block until the Stop signal is received
	s.exit = make(chan struct{})
	go s.run()
	return nil
}

// Stop is called when the service is stopping
func (s *Server) Stop(v service.Service) error {
	return s.Shutdown()
}

// Shutdown stops the server
func (s *Server) Shutdown() error {
	s.logInfo("Service stopping")
	// Close the channel, this will automatically release the block
	s.shutdown = make(chan struct{})
	close(s.exit)
	// Wait for the shutdown to complete
	_ = <-s.shutdown
	return nil
}

// GetServiceList searches for services based on the search criteria passed in the request
func (s *Server) GetServiceList(r GetRequest) (GetResponse, error) {
	if s.WaitTime <= 0 {
		s.WaitTime = 3
	}
	wt := r.WaitTime
	if wt <= 0 {
		wt = s.WaitTime
	}

	if r.ServiceType == "" {
		r.ServiceType = s.Config.DefaultServiceType
	}
	if r.Domain == "" {
		r.Domain = "local"
	}

	resp := r.CreateResponse()
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		s.logError("Failed to initialize zeroconf resolver.", err.Error())
		return resp, err
	}
	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			resp.Services = append(resp.Services, NewServiceItemFromZeroConf(entry))
		}
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(wt))
	defer cancel()
	err = resolver.Browse(ctx, r.ServiceType, r.Domain, entries)
	if err != nil {
		s.logError(fmt.Sprintf("Failed to browse for services with ServiceType '%s' on Domain '%s'.", r.ServiceType, r.Domain), err.Error())
		return resp, err
	}

	<-ctx.Done()

	return resp, nil
}

// RegisterService registers the service in the specified request
func (s *Server) RegisterService(r *RegisterRequest) RegisterResponse {
	s.regLock.Lock()
	defer s.regLock.Unlock()

	// Check if this service is already registered
	n := NewServerFromRequest(r, s)
	e := s.regList[r.ID]
	addNew := true
	if e != nil {
		if n.ID == e.ID {
			if n.IsDifferentFrom(e) {
				// Remove the existing one
				s.logInfo(fmt.Sprintf("Deregistering existing service %s: %s", e.ID, e.Name))
				e.Stop()
				delete(s.regList, r.ID)
			} else {
				s.logInfo(fmt.Sprintf("Confirming existing service %s: %s", e.ID, e.Name))
				e.LastContact = time.Now()
				addNew = false
			}
		}
	}
	if addNew {
		s.logInfo(fmt.Sprintf("Registering new service %s: %s", n.ID, n.Name))
		s.regList[r.ID] = n
		n.Start()
	}
	return r.CreateResponse()
}

// DeregisterService removes the service registration
func (s *Server) DeregisterService(id string) {
	s.regLock.Lock()
	defer s.regLock.Unlock()

	// Check to see if this service is already registered
	e := s.regList[id]
	if e != nil {
		if e.ID == id {
			s.logInfo(fmt.Sprintf("Deregistering existing service %s: %s", e.ID, e.Name))
			e.Stop()
			delete(s.regList, id)
		}
	}
}

func (s *Server) run() {
	if s.PortNo < 0 {
		s.PortNo = 20404
	}

	// Get the configuration
	if s.Config == nil {
		s.Config = &Config{}
	}
	s.Config.ReadFromFile("config.json")

	// Create a router
	s.router = mux.NewRouter().StrictSlash(true)

	// Add the controllers
	s.addController(new(ServiceController))
	s.addController(new(OnlineController))

	// Create an HTTP server
	// We lock to the loopback so that this service is not visible externally
	s.http = &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", s.PortNo),
		Handler: s.router,
	}

	// Register this service
	if hn, err := os.Hostname(); err != nil {
		s.hostName = s.Config.ID
	} else {
		s.hostName = hn
	}
	s.RegisterService(&RegisterRequest{
		ID:          s.Config.ID,
		Name:        s.Config.Name,
		PortNo:      s.PortNo,
		ServiceType: s.Config.DefaultServiceType,
		Text:        []string{fmt.Sprintf("id=%s", s.Config.ID)},
	})

	// Start the web server
	go func() {
		s.logInfo("Server listening on port", s.PortNo)
		if err := s.http.ListenAndServe(); err != nil {
			msg := err.Error()
			if !strings.Contains(msg, "http: Server closed") {
				s.logError("Error starting Web Server.", err.Error())
			}
		}
	}()

	// Wait for an exit signal
	_ = <-s.exit

	// Shutdown the HTTP server
	s.http.Shutdown(nil)

	// Shutdown the registered services
	s.logDebug("Deregistering service registrations.")
	for _, i := range s.regList {
		s.DeregisterService(i.ID)
	}

	s.logDebug("Shutdown complete")
	close(s.shutdown)
}

// logDebug logs a debug message to the logger
func (s *Server) logDebug(v ...interface{}) {
	if s.Debug {
		a := fmt.Sprint(v)
		logger.Info("Server: [Dbg] ", a[1:len(a)-1])
	}
}

// logInfo logs an information message to the logger
func (s *Server) logInfo(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Info("Server: [Inf] ", a[1:len(a)-1])
}

// logError logs an error message to the logger
func (s *Server) logError(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Error("Server: [Err] ", a[1:len(a)-1])
}
