package trigger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/faasaf/frameworks/common"
)

type trigger struct {
	runtimePort int
	triggerFn   triggerFn
}

func (t *trigger) run() {
	ctxCh := make(chan ContextWrapper)
	errCh := make(chan error)
	go t.triggerFn(ctxCh, errCh)
	// Handle new contexts
	go func() {
		for ctx := range ctxCh {
			go t.handleTrigger(ctx)
		}
	}()
	// Handle errors
	go func() {
		for err := range errCh {
			log.Error(err)
		}
	}()
}

func (t *trigger) handleTrigger(ctx ContextWrapper) {
	log.Debug("delegating further processing to the faasaf runtime")
	bodyBytes, err := json.Marshal(ctx.GetContext())
	if err != nil {
		ctx.ErrC() <- fmt.Errorf("error marshaling context: %s", err)
		return
	}
	res, err := http.Post(
		fmt.Sprintf("http://localhost:%d", t.runtimePort),
		"application/json",
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		ctx.ErrC() <- fmt.Errorf(
			"error delegating further processing to the faasaf runtime: %s",
			err,
		)
		return
	}
	if res.StatusCode != http.StatusOK {
		ctx.ErrC() <- fmt.Errorf(
			"the faasaf runtime returned status code: %d",
			res.StatusCode,
		)
		return
	}
	defer res.Body.Close()
	bodyBytes, err = ioutil.ReadAll(res.Body)
	if err != nil {
		ctx.ErrC() <- fmt.Errorf(
			"error reading response body: %s",
			err,
		)
		return
	}
	c := common.NewContext()
	if err := json.Unmarshal(bodyBytes, c); err != nil {
		ctx.ErrC() <- fmt.Errorf(
			"error unmarshaling response: %s",
			err,
		)
		return
	}
	ctx.ResC() <- c
	return
}
