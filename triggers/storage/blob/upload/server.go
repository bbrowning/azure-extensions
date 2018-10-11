package main

import (
	"fmt"
	"net/http"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/faasaf/frameworks/trigger"
)

var blobURLRegex = regexp.MustCompile(
	`^https:\/\/([a-zA-Z]+).blob.core.windows.net\/(\w+)\/(.+)$`,
)

type server struct {
	name                string
	port                int
	blobURLContextKey   string
	accountContextKey   string
	containerContextKey string
	blobPathContextKey  string
	ctxCh               chan trigger.ContextWrapper
	errCh               chan error
}

func (s *server) Run(
	ctxCh chan trigger.ContextWrapper,
	errCh chan error,
) error {
	s.ctxCh = ctxCh
	s.errCh = errCh
	bindAddress := fmt.Sprintf(":%d", s.port)
	log.WithField(
		"bindAddress", bindAddress,
	).Info("listening for events only")
	return http.ListenAndServe(bindAddress, s.router())
}
