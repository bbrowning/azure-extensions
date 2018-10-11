package trigger

import (
	"os"

	log "github.com/Sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v1"
)

var emptyJSONBytes = []byte("{}")

type initFn func(Config) error
type triggerFn func(chan ContextWrapper, chan error)

func Run(
	name string,
	version string,
	description string,
	initFn initFn,
	triggerFn triggerFn,
) {
	app := cli.NewApp()
	app.Name = name
	app.Version = version
	app.Usage = ""
	app.Description = description
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "runtime-port, r",
			Value: 8080,
			Usage: "port number that the faasaf runtime is listening on",
		},
		cli.IntFlag{
			Name:  "port, p",
			Value: 8081,
			Usage: "port number to listen on (health checks only)",
		},
		cli.StringSliceFlag{
			Name: "set, s",
			Usage: "configure the trigger using key=value pairs; this flag may be " +
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
		trigger := &trigger{
			runtimePort: c.Int("runtime-port"),
			triggerFn:   triggerFn,
		}
		go trigger.run()
		return (&server{
			port: c.Int("port"),
		}).run()
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
