package trigger

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *server) router() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/healthz", s.healthz).Methods(http.MethodGet)
	return r
}
