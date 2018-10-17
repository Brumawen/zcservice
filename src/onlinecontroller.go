package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// OnlineController handles the web methods for registering and discovering services
type OnlineController struct {
	Srv *Server
}

// AddController adds the controller routes to the router
func (c *OnlineController) AddController(router *mux.Router, s *Server) {
	c.Srv = s
	router.Methods("GET", "POST").Path("/online").
		Handler(Logger(c, http.HandlerFunc(c.handleOnline)))
	// router.Methods("POST").Path("/shutdown").
	// 	Handler(Logger(c, http.HandlerFunc(c.handleShutdown)))
}

// handleOnline handles the /online web method call
func (c *OnlineController) handleOnline(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("true"))
}

// handleOnline handles the /shutdown web method call
// func (c *OnlineController) handleShutdown(w http.ResponseWriter, r *http.Request) {
// 	c.Srv.Shutdown()
// }

// LogInfo is used to log information messages for this controller.
func (c *OnlineController) LogInfo(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Info("OnlineController: [Inf] ", a[1:len(a)-1])
}
