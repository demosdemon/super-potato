package app

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type Command func(app *App) *cobra.Command

type App struct {
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

func New() *App {
	return &App{
		Fs:     afero.NewOsFs(),
		Exit:   logrus.Exit,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

func (a *App) Execute(command Command) {
	cmd := command(a)

	if err := cmd.Execute(); err != nil {
		a.Exit(1)
	}

	a.Exit(0)
}
