package server

import "net/http"

type healthStatus struct {
	Status bool
}

func (s *Server) CheckLive(rw http.ResponseWriter, req *http.Request) {
	if s.IsHealthy == nil || !s.IsHealthy.Load().(bool) {
		s.encodeJSON(rw, req, http.StatusInternalServerError, healthStatus{Status: false})
	}
	s.encodeJSON(rw, req, http.StatusOK, healthStatus{Status: true})
}

func (s *Server) CheckReady(rw http.ResponseWriter, req *http.Request) {
	if s.IsReady == nil || !s.IsReady.Load().(bool) {
		s.encodeJSON(rw, req, http.StatusInternalServerError, healthStatus{Status: false})
	}
	s.encodeJSON(rw, req, http.StatusOK, healthStatus{Status: true})
}

func (s *Server) GetTest(rw http.ResponseWriter, req *http.Request) {
	if s.IsReady == nil || !s.IsReady.Load().(bool) {
		s.encodeJSON(rw, req, http.StatusInternalServerError, healthStatus{Status: false})
	}
	s.encodeJSON(rw, req, http.StatusOK, healthStatus{Status: true})
}
