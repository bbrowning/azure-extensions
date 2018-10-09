package binding

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type server struct {
	port   int
	bindFn bindFn
}

func (s *server) run() error {
	bindAddress := fmt.Sprintf("localhost:%d", s.port)
	log.WithField(
		"bindAddress", bindAddress,
	).Info("listening")
	return http.ListenAndServe(bindAddress, s.router())
}
