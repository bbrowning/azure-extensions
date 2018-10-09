package binding

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/faasaf/frameworks/common"
	cli "gopkg.in/urfave/cli.v1"
)

var emptyJSONBytes = []byte("{}")

type initFn func(Config) error
type bindFn func(common.Context) error

func Run(
	name string,
	version string,
	description string,
	initFn initFn,
	bindFn bindFn,
) {
	app := cli.NewApp()
	app.Name = name
	app.Version = version
	app.Usage = ""
	app.Description = description
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Value: 8082,
			Usage: "port number to listen on",
		},
		cli.StringSliceFlag{
			Name: "set, s",
			Usage: "configure the binding using key=value pairs; this flag may be " +
				"applied multiple times",
		},
		cli.StringFlag{
			Name:  "log-level, ll",
			Value: "info",
		},
	}
	app.Action = func(c *cli.Context) error {
		logLevelStr := c.String("log-level")
		logLevel, err := log.ParseLevel(logLevelStr)
		if err != nil {
			return err
		}
		log.SetLevel(logLevel)

		cfg, err := newConfig(c.StringSlice("set"))
		if err != nil {
			return err
		}
		if initFn != nil {
			if err := initFn(cfg); err != nil {
				return err
			}
		}
		return (&server{
			port:   c.Int("port"),
			bindFn: bindFn,
		}).run()
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
