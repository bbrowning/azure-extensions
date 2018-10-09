package trigger

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func (s *server) healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(emptyJSONBytes); err != nil {
		log.WithField(
			"error", err,
		).Error("error writing /healthz response")
	}
}
