package trigger

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type server struct {
	port int
}

func (s *server) run() error {
	bindAddress := fmt.Sprintf("localhost:%d", s.port)
	log.WithField(
		"bindAddress", bindAddress,
	).Info("listening for health checks only")
	return http.ListenAndServe(bindAddress, s.router())
}
