//go:generate time go run ./gen

package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/cmd/deploy"
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
	Prefix    string `flag:"prefix" desc:"The prefix for Platform.sh environment variables."`
}

func (c *Config) Use() string {
	return "super-potato"
}

func (c *Config) Args(cmd *cobra.Command, args []string) error {
	return nil
}

func (c *Config) PersistentPreRun(cmd *cobra.Command, args []string) error {
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
	c.SetPrefix(c.Prefix)

	logrus.Trace("program beginning")
	return nil
}

func (c *Config) PersistentPostRun(cmd *cobra.Command, args []string) error {
	logrus.Trace("program ending")
	return nil
}

func (c *Config) SubCommands() []app.Config {
	return []app.Config{
		deploy.New(c.App),
		dump.New(c.App),
		scrape.New(c.App),
		secret.New(c.App),
		serve.New(c.App),
	}
}

func main() {
	inst, cancel := app.New(context.Background())
	inst.Execute(&Config{
		App:       inst,
		LogLevel:  "trace",
		LogOutput: "/dev/stderr",
		Prefix:    "PLATFORM_",
	})
	cancel()
}
