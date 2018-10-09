package binding

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/faasaf/frameworks/common"
)

func (s *server) bind(w http.ResponseWriter, r *http.Request) {
	ctx, err := s.getBindContext(r)
	if err != nil {
		log.WithField(
			"error", err,
		).Debug("received bad binding request")
		s.writeBindResponse(w, http.StatusBadRequest, ctx)
		return
	}

	log.WithFields(
		ctx.GetRawMap(),
	).Debug("received binding request from faasaf-runtime")

	if err = s.bindFn(ctx); err != nil {
		log.WithField(
			"error", err,
		).Error("error executing binding")
		s.writeBindResponse(w, http.StatusInternalServerError, ctx)
		return
	}

	log.Debug("completed binding request")

	s.writeBindResponse(w, http.StatusOK, ctx)
}

func (s *server) getBindContext(r *http.Request) (common.Context, error) {
	ctx := common.NewContext()

	defer r.Body.Close()
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ctx, err
	}

	err = json.Unmarshal(bodyBytes, ctx)
	return ctx, err
}

func (s *server) writeBindResponse(
	w http.ResponseWriter,
	statusCode int,
	ctx common.Context,
) {
	w.Header().Set("Content-Type", "application/json")

	bodyBytes, err := json.Marshal(ctx)
	if err != nil {
		log.WithField(
			"error", err,
		).Error("error marshaling response context")
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write(emptyJSONBytes); err != nil {
			log.WithField(
				"error", err,
			).Error("error writing err response")
		}
		return
	}

	w.WriteHeader(statusCode)
	if _, err := w.Write(bodyBytes); err != nil {
		log.WithField(
			"error", err,
		).Error("error writing response context")
	}
}
