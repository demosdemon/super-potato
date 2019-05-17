package app

import (
	"bytes"
	"context"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/demosdemon/super-potato/pkg/platformsh"
)

func init() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(os.Stderr)
	log.SetOutput(os.Stderr)
	setReportCaller()
	logrus.Trace("init app")
}

func setReportCaller() {
	if caller := os.Getenv("PKI_LOG_CALLER"); caller == "1" || caller == "true" {
		logrus.SetReportCaller(true)
	} else {
		logrus.SetReportCaller(false)
	}
}

type LogLogger interface {
	Output(calldepth int, s string) error
}

type App struct {
	context.Context
	afero.Fs
	platformsh.Environment

	Exit   func(int)
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func New(ctx context.Context) (*App, context.CancelFunc) {
	ctx, cancel := CancelOnSignal(ctx, syscall.SIGINT, syscall.SIGTERM)

	cwd, _ := os.Getwd()
	fs := afero.NewBasePathFs(afero.NewOsFs(), cwd)

	app := &App{
		Context: ctx,
		Fs:      fs,
		Exit:    logrus.Exit,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}
	go app.logMemoryTick()

	return app, cancel
}

func (a *App) SetPrefix(s string) {
	a.Environment = platformsh.NewEnvironment(s)
}

func (a *App) Logger(prefix string) LogLogger {
	return log.New(a.Stderr, prefix, log.LstdFlags)
}

func (a *App) Execute(cfg Config) {
	cmd := a.command(cfg)

	if err := cmd.Execute(); err != nil {
		a.Exit(1)
	}

	a.Exit(0)
}

func (a *App) LogMemory() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	logrus.WithFields(logrus.Fields{
		"Alloc":        m.Alloc,
		"AllocMB":      m.Alloc >> 20,
		"TotalAlloc":   m.TotalAlloc,
		"TotalAllocMB": m.TotalAlloc >> 20,
		"Sys":          m.Sys,
		"SysMB":        m.Sys >> 20,
		"NumGC":        m.NumGC,
	}).Debug("memory stats")
}

func (a *App) logMemoryTick() {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			a.LogMemory()
		case <-a.Done():
			ticker.Stop()
			return
		}
	}
}

func (a *App) GetInput(s string) (io.ReadCloser, error) {
	switch s {
	case "-", "/dev/stdin":
		return ioutil.NopCloser(a.Stdin), nil
	case "/dev/null":
		return ioutil.NopCloser(new(bytes.Buffer)), nil
	default:
		return a.Open(s)
	}
}

func (a *App) GetOutput(s string) (io.WriteCloser, error) {
	switch s {
	case "-", "/dev/stdout":
		return NewNopWriterCloser(a.Stdout)
	case "/dev/stderr":
		return NewNopWriterCloser(a.Stderr)
	case "/dev/null":
		return NewNopWriterCloser(ioutil.Discard)
	default:
		return a.Create(s)
	}
}

func (a *App) Append(s string) (io.WriteCloser, error) {
	switch s {
	case "-", "/dev/stdout":
		return NewNopWriterCloser(a.Stdout)
	case "/dev/stderr":
		return NewNopWriterCloser(a.Stderr)
	case "/dev/null":
		return NewNopWriterCloser(ioutil.Discard)
	default:
		return a.OpenFile(s, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	}
}

func (a *App) ReadYAML(path string, data interface{}) error {
	fp, err := a.GetInput(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	dec := yaml.NewDecoder(fp)
	return dec.Decode(data)
}

type cobraRunner func(cmd *cobra.Command, args []string) error

func (a *App) command(cfg Config) *cobra.Command {
	var persistentPreRun, preRun, run, postRun, persistentPostRun cobraRunner
	var subCommands []Config

	if c, ok := cfg.(RootRunner); ok {
		persistentPreRun = c.PersistentPreRun
		persistentPostRun = c.PersistentPostRun
	}
	if c, ok := cfg.(PreRunner); ok {
		preRun = c.PreRun
	}
	if c, ok := cfg.(Runner); ok {
		run = c.Run
	}
	if c, ok := cfg.(PostRunner); ok {
		postRun = c.PostRun
	}
	if c, ok := cfg.(MasterRunner); ok {
		subCommands = c.SubCommands()
	}

	cmd := cobra.Command{
		Use:                cfg.Use(),
		Args:               cfg.Args,
		PersistentPreRunE:  persistentPreRun,
		PreRunE:            preRun,
		RunE:               run,
		PostRunE:           postRun,
		PersistentPostRunE: persistentPostRun,
	}
	cmd.SetOutput(a.Stdout)

	var flags *pflag.FlagSet
	if _, ok := cfg.(RootRunner); ok {
		flags = cmd.PersistentFlags()
	} else {
		flags = cmd.Flags()
	}

	err := gpflag.ParseTo(cfg, flags)
	if err != nil {
		logrus.WithError(err).Panic("failed to parse config flags")
	}

	for _, subCfg := range subCommands {
		cmd.AddCommand(a.command(subCfg))
	}

	return &cmd
}
