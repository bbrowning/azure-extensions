package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *server) router() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/", s.handleEvent).Methods(http.MethodPost)
	return r
}
