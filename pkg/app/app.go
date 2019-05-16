package app

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(os.Stderr)
	log.SetOutput(os.Stderr)

	if caller := os.Getenv("PKI_LOG_CALLER"); caller == "1" || caller == "true" {
		logrus.SetReportCaller(true)
	} else {
		logrus.SetReportCaller(false)
	}
}

type Command func(app *App) *cobra.Command

type LogLogger interface {
	Output(calldepth int, s string) error
}

type App struct {
	context.Context
	afero.Fs
	Exit   func(int)
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (a *App) Logger(prefix string) LogLogger {
	return log.New(a.Stderr, prefix, log.LstdFlags)
}

func New(ctx context.Context) (*App, context.CancelFunc) {
	ctx, cancel := CancelOnSignal(ctx, syscall.SIGINT, syscall.SIGTERM)

	logrus.SetOutput(os.Stderr)

	app := &App{
		Context: ctx,
		Fs:      afero.NewOsFs(),
		Exit:    logrus.Exit,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}

	return app, cancel
}

func (a *App) Execute(command Command) {
	cmd := command(a)

	if err := cmd.Execute(); err != nil {
		a.Exit(1)
	}

	a.Exit(0)
}

func CancelOnSignal(ctx context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan os.Signal, len(signals))
	signal.Notify(ch, signals...)

	go func() {
		select {
		case sig := <-ch:
			logrus.WithField("signal", sig).Debug("received signal")
		case <-ctx.Done():
			logrus.WithField("err", ctx.Err()).Debug("context done")
		}

		signal.Stop(ch)
		close(ch)
		cancel()
	}()

	return ctx, cancel
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
		return a.OpenFile(s, os.O_APPEND|os.O_WRONLY, 0644)
	}
}

func (a *App) ReadYAML(path string, data interface{}) error {
	fp, err := a.GetInput(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	dec := yaml.NewDecoder(fp)
	return dec.Decode(&data)
}
