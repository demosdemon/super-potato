package main

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/gen/internal/gen"
	_ "github.com/demosdemon/super-potato/gen/internal/gen/api"
	_ "github.com/demosdemon/super-potato/gen/internal/gen/enums"
	_ "github.com/demosdemon/super-potato/gen/internal/gen/variables"
	"github.com/demosdemon/super-potato/pkg/app"
)

func main() {
	inst, cancel := app.New(context.Background())
	inst.Execute(&Config{App: inst})
	cancel()
}

type Job struct {
	Template string `yaml:"template"`
	Input    string `yaml:"input"`
	Output   string `yaml:"output"`
}

type ConfigFile struct {
	Jobs []Job `yaml:"jobs"`
}

type Config struct {
	*app.App `flag:"-"`
	ExitCode bool `flag:"exit-code" desc:"If specified, the exit code will be the number of files written plus 1. An exit code of 1 indicates program error."`
}

func (c *Config) Use() string {
	return "gen [CONFIG]"
}

func (c *Config) Args(cmd *cobra.Command, args []string) error {
	return cobra.MaximumNArgs(1)(cmd, args)
}

func (c *Config) Run(cmd *cobra.Command, args []string) error {
	path := "./.gen.yaml"
	if len(args) == 1 {
		path = args[0]
	}

	var cfg ConfigFile
	err := c.ReadYAML(path, &cfg)
	if err != nil {
		return err
	}

	flags := cmd.Flags()
	exitCode, err := flags.GetBool("exit-code")
	if err != nil {
		return err
	}

	logrus.WithField("cfg", cfg).WithField("exitCode", exitCode).Trace()

	var hasErr uint32
	var wg sync.WaitGroup
	var written uint32
	for _, j := range cfg.Jobs {
		wg.Add(1)

		go func(j Job) {
			defer wg.Done()
			input, err := c.GetInput(j.Input)
			if err != nil {
				atomic.AddUint32(&hasErr, 1)
				logrus.WithField("j", j).WithField("err", err).Error("unable to open input")
				return
			}

			fn, ok := gen.DefaultRenderMap[j.Template]
			if !ok {
				atomic.AddUint32(&hasErr, 1)
				logrus.WithField("j", j).Error("invalid template")
				return
			}

			renderer, err := fn(input)
			input.Close()
			if err != nil {
				atomic.AddUint32(&hasErr, 1)
				logrus.WithField("j", j).WithField("err", err).Error("unable to parse input")
				return
			}

			err = gen.Render(renderer, j.Output, c)
			logrus.WithField("j", j).WithField("err", err).Trace()

			switch err {
			case nil:
				atomic.AddUint32(&written, 1)
			case gen.ErrNoChange:
			default:
				atomic.AddUint32(&hasErr, 1)
				logrus.WithField("j", j).WithField("err", err).Error("unable to render template")
			}
		}(j)
	}

	wg.Wait()

	if hasErr > 0 {
		c.Exit(1)
	}

	if written > 0 && exitCode {
		c.Exit(int(written) + 1)
	}

	return nil
}
