package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// ServiceController handles the web methods for registering and discovering services
type ServiceController struct {
	Srv *Server
}

// AddController adds the controller routes to the router
func (c *ServiceController) AddController(router *mux.Router, s *Server) {
	c.Srv = s
	router.Methods("POST", "GET").Path("/service/get").
		Handler(Logger(c, http.HandlerFunc(c.handleGet)))
	router.Methods("POST").Path("/service/add").
		Handler(Logger(c, http.HandlerFunc(c.handleAdd)))
	router.Methods("DELETE").Path("/service/remove/{id}").
		Handler(Logger(c, http.HandlerFunc(c.handleRemove)))
}

func (c *ServiceController) handleGet(w http.ResponseWriter, r *http.Request) {
	req := GetRequest{}
	if r.ContentLength != 0 {
		req.ReadFrom(r.Body)
	}
	if resp, err := c.Srv.GetServiceList(req); err != nil {
		http.Error(w, err.Error(), 400)
	} else {
		resp.WriteTo(w)
	}
}

func (c *ServiceController) handleAdd(w http.ResponseWriter, r *http.Request) {
	req := RegisterRequest{}
	req.ReadFrom(r.Body)
	if req.ID == "" {
		http.Error(w, "ID is missing.", 400)
		return
	}
	if req.Name == "" {
		http.Error(w, "Service Name is missing.", 400)
		return
	}
	if req.PortNo <= 0 {
		http.Error(w, "Invalid Port Number.", 400)
		return
	}
	resp := c.Srv.RegisterService(&req)
	resp.WriteTo(w)
}

func (c *ServiceController) handleRemove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid ID", 400)
	} else {
		go c.Srv.DeregisterService(id)
	}
}

// LogInfo is used to log information messages for this controller.
func (c *ServiceController) LogInfo(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Info("ServiceController: [Inf] ", a[1:len(a)-1])
}
