package app

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type Command func(app *App) *cobra.Command

type App struct {
	context.Context
	afero.Fs
	Exit   func(int)
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func init() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)
}

func New(ctx context.Context) (*App, context.CancelFunc) {
	ctx, cancel := CancelOnSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
	return &App{
		Context: ctx,
		Fs:      afero.NewOsFs(),
		Exit:    logrus.Exit,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}, cancel
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
