package trigger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/faasaf/frameworks/common"
)

type trigger struct {
	runtimePort int
	triggerFn   triggerFn
}

func (t *trigger) run() {
	ctxCh := make(chan common.Context)
	errCh := make(chan error)
	go t.triggerFn(ctxCh, errCh)
	// Handle new contexts
	go func() {
		for ctx := range ctxCh {
			go t.handleTrigger(ctx, errCh)
		}
	}()
	// Handle errors
	go func() {
		for err := range errCh {
			log.Error(err)
		}
	}()
}

func (t *trigger) handleTrigger(ctx common.Context, errCh chan error) {
	log.Debug("delegating further processing to the faasaf runtime")
	bodyBytes, err := json.Marshal(ctx)
	if err != nil {
		errCh <- fmt.Errorf("error marshaling context: %s", err)
		return
	}
	res, err := http.Post(
		fmt.Sprintf("http://localhost:%d", t.runtimePort),
		"application/json",
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		errCh <- fmt.Errorf(
			"error delegating further processing to the faasaf runtime: %s",
			err,
		)
		return
	}
	if res.StatusCode != http.StatusOK {
		errCh <- fmt.Errorf(
			"the faasaf runtime returned status code: %d",
			res.StatusCode,
		)
		return
	}
	return
}
