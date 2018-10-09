package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/faasaf/frameworks/common"
)

// For reference:
// {
//   "topic": "/subscriptions/319a9601-1ec0-0000-aebc-8fe82724c81e/resourceGroups/testrg/providers/Microsoft.Storage/storageAccounts/myaccount",
//   "subject": "/blobServices/default/containers/testcontainer/blobs/file1.txt",
//   "eventType": "Microsoft.Storage.BlobCreated",
//   "eventTime": "2017-08-16T01:57:26.005121Z",
//   "id": "602a88ef-0001-00e6-1233-1646070610ea",
//   "data": {
//     "api": "PutBlockList",
//     "clientRequestId": "799304a4-bbc5-45b6-9849-ec2c66be800a",
//     "requestId": "602a88ef-0001-00e6-1233-164607000000",
//     "eTag": "0x8D4E44A24ABE7F1",
//     "contentType": "text/plain",
//     "contentLength": 447,
//     "blobType": "BlockBlob",
//     "url": "https://myaccount.blob.core.windows.net/testcontainer/file1.txt",
//     "sequencer": "00000000000000EB000000000000C65A",
//   },
//   "dataVersion": "",
//   "metadataVersion": "1"
// }

type event struct {
	Data eventData `json:"data"`
}

type eventData struct {
	URL string `json:"url"`
}

func (s *server) handleEvent(w http.ResponseWriter, r *http.Request) {
	evt := event{}

	defer r.Body.Close()
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.errCh <- err
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bodyBytes, &evt)
	if err != nil {
		s.errCh <- err
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.WithField(
		"blobUrl", evt.Data.URL,
	).Debug("received event")

	matches := blobURLRegex.FindStringSubmatch(evt.Data.URL)
	if len(matches) == 0 {
		s.errCh <- fmt.Errorf(
			"blob URL is not a valid Azure Storage URL: %s",
			evt.Data.URL,
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := common.NewContext()

	if s.blobURLContextKey != "" {
		log.WithField(
			"key", s.blobURLContextKey,
		).WithField(
			"value", evt.Data.URL,
		).Debug("updating context")
		ctx.Set(s.blobURLContextKey, evt.Data.URL)
	}

	if s.accountContextKey != "" {
		log.WithField(
			"key", s.accountContextKey,
		).WithField(
			"value", matches[1],
		).Debug("updating context")
		ctx.Set(s.accountContextKey, matches[1])
	}

	if s.containerContextKey != "" {
		log.WithField(
			"key", s.containerContextKey,
		).WithField(
			"value", matches[2],
		).Debug("updating context")
		ctx.Set(s.containerContextKey, matches[2])
	}

	if s.blobPathContextKey != "" {
		log.WithField(
			"key", s.blobPathContextKey,
		).WithField(
			"value", matches[3],
		).Debug("updating context")
		ctx.Set(s.blobPathContextKey, matches[3])
	}

	s.ctxCh <- ctx

	// TODO: Is it premature to respond now?
	// TODO: I think we should probably wait a response-- can use a one-off
	// channel to do this. i.e. Wrap context and response channel in some object
	// and put THAT on ctxCh.
	w.WriteHeader(http.StatusOK)
}
