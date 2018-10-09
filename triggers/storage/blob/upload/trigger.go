package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/faasaf/frameworks/common"
	"github.com/faasaf/frameworks/trigger"
)

const name = "azure-storage-upload"

// Value is injected by the build
var version string

func main() {

	var srvr *server

	trigger.Run(
		name,
		version,
		"A FaaSAF trigger that is tripped when a blob is uploaded to Azure Storage",
		func(cfg trigger.Config) error { // Initializations
			portStr := cfg.GetSetting("eventPort", "")
			if portStr == "" {
				return errors.New(
					"the event listener port (eventPort) was not specified",
				)
			}

			port, err := strconv.Atoi(portStr)
			if err != nil {
				return fmt.Errorf(
					`the specified listener port "%s" could not be pased as an integer`,
					portStr,
				)
			}

			srvr = &server{
				name:                name,
				port:                port,
				blobURLContextKey:   cfg.GetSetting("blobUrlContextKey", ""),
				accountContextKey:   cfg.GetSetting("accountContextKey", ""),
				containerContextKey: cfg.GetSetting("containerContextKey", ""),
				blobPathContextKey:  cfg.GetSetting("blobPathContextKey", ""),
			}

			return nil
		},
		func(ctxCh chan common.Context, errCh chan error) { // Actual functionality
			if err := srvr.Run(ctxCh, errCh); err != nil {
				errCh <- err
			}
		},
	)

}
