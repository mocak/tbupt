package controllers

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	r  *mux.Router
	dc *Decks
}

// NewServer returns new server instance
func NewServer(dc *Decks) *Server {
	return &Server{dc: dc}
}

// ServeHTTP dispatches the handler registered in the matched route.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r = mux.NewRouter()
	s.r.HandleFunc("/deck", s.dc.Create).Methods("POST")
	s.r.HandleFunc("/deck/{uuid}/open", s.dc.Open).Methods("PUT")
	s.r.HandleFunc("/deck/{uuid}/draw", s.dc.Draw).Methods("POST")

	s.r.ServeHTTP(w, r)
}
