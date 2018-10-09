package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"

	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/2018-03-28/azblob"
	log "github.com/Sirupsen/logrus"
	"github.com/faasaf/frameworks/binding"
	"github.com/faasaf/frameworks/common"
)

// Value is injected by the build
var version string

func main() {

	var blobURLRegex = regexp.MustCompile(
		`^https:\/\/([a-zA-Z]+).blob.core.windows.net\/(\w+)\/(.+)$`,
	)
	var accountName string
	var containerName string
	var blobURLContextKey string
	var containerContextKey string
	var blobPathContextKey string
	var localFilePathContextKey string
	var pipeline pipeline.Pipeline

	binding.Run(
		"azure-storage-blob-download",
		version,
		"A FaaSAF binding to download a blob from Azure Storage",
		func(cfg binding.Config) error { // Initialization
			accountName = cfg.GetSetting("account", "")
			if accountName == "" {
				return errors.New(
					"the azure storage account (account) was not specified",
				)
			}

			accessKey := cfg.GetSetting("accessKey", "")
			if accessKey == "" {
				return errors.New(
					"the azure storage access key (accessKey) was not specified",
				)
			}

			containerName = cfg.GetSetting("container", "")

			log.WithField(
				"account", accountName,
			).WithField(
				"container", containerName,
			).Info("Azure storage configured")

			blobURLContextKey = cfg.GetSetting("blobUrlContextKey", "")
			containerContextKey = cfg.GetSetting("containerContextKey", "")
			blobPathContextKey = cfg.GetSetting("blobPathContextKey", "")
			localFilePathContextKey = cfg.GetSetting("localFilePathContextKey", "")

			log.WithField(
				"blobUrlContextKey", blobURLContextKey,
			).WithField(
				"containerContextKey", containerContextKey,
			).WithField(
				"blobPathContextKey", blobPathContextKey,
			).WithField(
				"localFilePathContextKey", localFilePathContextKey,
			).Debug("context keys configured")

			if containerName == "" &&
				blobURLContextKey == "" &&
				(containerContextKey == "" ||
					blobPathContextKey == "") {
				return errors.New(
					"the container name (container) OR the blob URL context key " +
						"(blobUrlContextKey) OR BOTH OF container context key " +
						"(containerContextKey) and blob path context key " +
						"(blobPathContextKey) must be specified",
				)
			}

			if localFilePathContextKey == "" {
				return errors.New(
					"the local file path context key (localFilePathContextKey) was not " +
						"specified",
				)
			}

			credential, err := azblob.NewSharedKeyCredential(accountName, accessKey)
			if err != nil {
				return err
			}

			pipeline = azblob.NewPipeline(credential, azblob.PipelineOptions{})

			return nil
		},
		func(ctx common.Context) error { // Actual functionality
			var blobURLStr string

			if blobURLContextKey != "" {
				var err error
				blobURLStr, err = ctx.GetString(blobURLContextKey, "")
				if err != nil {
					return err
				}
				if blobURLStr == "" {
					return fmt.Errorf(
						"did not receive blob URL (%s) in context",
						blobURLContextKey,
					)
				}
				matches := blobURLRegex.FindStringSubmatch(blobURLStr)
				if len(matches) == 0 {
					return fmt.Errorf(
						"blob URL is not a valid Azure Storage URL: %s",
						blobURLStr,
					)
				}
				if matches[1] != accountName {
					return fmt.Errorf(
						`blob URL "%s" is not valid for configured storage account "%s"`,
						blobURLStr,
						accountName,
					)
				}
			} else {
				var container string
				if containerContextKey != "" {
					var err error
					container, err = ctx.GetString(containerContextKey, "")
					if err != nil {
						return err
					}
					if container == "" {
						return fmt.Errorf(
							"did not receive container (%s) in context",
							containerContextKey,
						)
					}
				} else {
					container = containerName
				}

				blobPath, err := ctx.GetString(blobPathContextKey, "")
				if err != nil {
					return err
				}
				if blobPath == "" {
					return fmt.Errorf(
						"did not receive blob path (%s) in context",
						blobPathContextKey,
					)
				}

				blobURLStr = fmt.Sprintf(
					"https://%s.blob.core.windows.net/%s/%s",
					accountName,
					container,
					blobPath,
				)
			}

			blobURL, err := url.Parse(blobURLStr)
			if err != nil {
				return err
			}

			file, err := ioutil.TempFile("", "")
			if err != nil {
				return err
			}
			defer file.Close()

			log.WithField(
				"url", blobURLStr,
			).WithField(
				"localfilePath", file.Name(),
			).Debug("downloading blob")
			if err := azblob.DownloadBlobToFile(
				context.Background(),
				azblob.NewBlobURL(*blobURL, pipeline),
				0,
				0,
				file,
				azblob.DownloadFromBlobOptions{},
			); err != nil {
				return fmt.Errorf(
					"error downloading blob \"%s\" to \"%s\": %s",
					blobURLStr,
					file.Name(),
					err,
				)
			}
			log.WithField(
				"url", blobURLStr,
			).WithField(
				"localfilePath", file.Name(),
			).Debug("downloaded blob")

			log.WithField(
				"key", localFilePathContextKey,
			).WithField(
				"value", file.Name(),
			).Debug("updating context")
			ctx.Set(localFilePathContextKey, file.Name())

			return nil
		},
	)

}
