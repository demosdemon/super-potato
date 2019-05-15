//go:generate time go run ./gen

package main

import (
	"context"

	"github.com/octago/sflags/gen/gpflag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/cmd/dump"
	"github.com/demosdemon/super-potato/cmd/scrape"
	"github.com/demosdemon/super-potato/cmd/secret"
	"github.com/demosdemon/super-potato/cmd/serve"
	"github.com/demosdemon/super-potato/pkg/app"
)

type Config struct {
	*app.App  `flag:"-"`
	LogLevel  string `flag:"log-level l" desc:"The logging verbosity"`
	LogOutput string `flag:"log-output" desc:"Where logging is written"`
}

func (c *Config) Run(cmd *cobra.Command, args []string) error {
	level, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		return err
	}

	fp, err := c.GetOutput(c.LogOutput)
	if err != nil {
		return err
	}

	logrus.SetLevel(level)
	logrus.SetOutput(fp)
	return nil
}

func Command(app *app.App) *cobra.Command {
	logrus.SetOutput(app.Stderr)

	cfg := Config{
		App:       app,
		LogLevel:  "trace",
		LogOutput: "/dev/stderr",
	}

	rv := cobra.Command{
		Use:               "super-potato",
		PersistentPreRunE: cfg.Run,
	}

	err := gpflag.ParseTo(&cfg, rv.PersistentFlags())
	if err != nil {
		logrus.WithField("err", err).Fatal("failed to parse config flags")
	}

	rv.AddCommand(dump.Command(app))
	rv.AddCommand(serve.Command(app))
	rv.AddCommand(scrape.Command(app))
	rv.AddCommand(secret.Command(app))

	return &rv
}

func main() {
	inst, cancel := app.New(context.Background())
	defer cancel()
	inst.Execute(Command)
}
