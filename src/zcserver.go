package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/satori/go.uuid"

	"github.com/grandcat/zeroconf"
)

// ZCServer defines a Zeroconf service registration
type ZCServer struct {
	ID          string    // ID of the service
	Name        string    // Service Instance Name
	PortNo      int       // Port number service is available on
	ServiceType string    // Service Type
	Domain      string    // Domain name
	Text        []string  // Associated Text
	LastContact time.Time // Date and time of last contact
	Srv         *Server   // Web Server
	shutdown    chan bool // Registration shutdown signal
	isRunning   bool      // Indicate whether currently running
}

// NewServerFromRequest creates a new server from the specified registration request
func NewServerFromRequest(r *RegisterRequest, srv *Server) *ZCServer {
	r.SetDefaults()
	if r.ID == "" {
		if uuid, err := uuid.NewV4(); err == nil {
			r.ID = strings.Replace(uuid.String(), "-", "", -1)
		}
	}
	if r.ServiceType == "" {
		r.ServiceType = srv.Config.DefaultServiceType
	}
	if r.Domain == "" {
		r.Domain = "local."
	}
	s := ZCServer{
		ID:          r.ID,
		Name:        r.Name,
		PortNo:      r.PortNo,
		ServiceType: r.ServiceType,
		Text:        r.Text,
		Domain:      r.Domain,
		LastContact: time.Now(),
	}
	return &s
}

// IsDifferentFrom returns whether or not the servers differ
func (s *ZCServer) IsDifferentFrom(i *ZCServer) bool {
	if s.ID != i.ID || s.PortNo != i.PortNo || s.Name != i.Name || s.ServiceType != i.ServiceType {
		return true
	}
	if len(s.Text) != len(i.Text) {
		return true
	}
	for x := range s.Text {
		if s.Text[x] != i.Text[x] {
			return true
		}
	}
	return false
}

// Start registers the service so that it is discoverable
func (s *ZCServer) Start() {
	if s.isRunning {
		return
	}
	s.shutdown = make(chan bool, 1)
	s.isRunning = true
	go s.register()
}

// Stop deregisters the service so that it is no longer discoverable
func (s *ZCServer) Stop() {
	if !s.isRunning {
		return
	}
	s.shutdown <- true
}

func (s *ZCServer) register() {
	if s.Domain == "" {
		s.Domain = "local."
	}
	s.logInfo("Registering service '" + s.Name + "'.")
	zsrv, err := zeroconf.Register(s.Name, s.ServiceType, s.Domain, s.PortNo, s.Text, nil)
	if err != nil {
		s.logError("Failed to register service '"+s.Name+"'. ", err.Error())
	}
	defer zsrv.Shutdown()

	// Ensure a clean exit
	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig, os.Interrupt, syscall.SIGTERM)

	select {
	case <-osSig:
		// Exit by user
	case <-s.shutdown:
		// Service shutdown
	}
}

// logDebug logs a debug message to the logger
func (s *ZCServer) logDebug(v ...interface{}) {
	if s.Srv.Debug {
		a := fmt.Sprint(v)
		logger.Info("ZCServer: [Dbg] ", a[1:len(a)-1])
	}
}

// logInfo logs an information message to the logger
func (s *ZCServer) logInfo(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Info("ZCServer: [Inf] ", a[1:len(a)-1])
}

// logError logs an error message to the logger
func (s *ZCServer) logError(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Error("ZCServer: [Err] ", a[1:len(a)-1])
}
