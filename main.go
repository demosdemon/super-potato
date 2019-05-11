//go:generate time go run . gen enums /data/enums.yaml /pkg/platformsh/enums_gen.go
//go:generate time go run . gen variables /data/variables.yaml /pkg/platformsh/environment_gen.go
//go:generate time go run . gen api /data/variables.yaml /cmd/serve/generated.go

package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/cmd/dump"
	"github.com/demosdemon/super-potato/cmd/gen"
	"github.com/demosdemon/super-potato/cmd/serve"
)

var (
	exit = logrus.Exit
	fs   afero.Fs
)

func init() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)

	if cwd, err := os.Getwd(); err == nil {
		fs = afero.NewBasePathFs(afero.NewOsFs(), cwd)
	} else {
		logrus.WithField("err", err).Warning("error getting CWD")
	}
}

func Command() *cobra.Command {
	rv := cobra.Command{
		Use: "super-potato",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ll, err := cmd.Flags().GetString("log-level")
			if err != nil {
				return err
			}

			level, err := logrus.ParseLevel(ll)
			if err != nil {
				return err
			}

			logrus.SetLevel(level)
			return nil
		},
	}

	flags := rv.PersistentFlags()
	flags.StringP("log-level", "l", "trace", "control the logging verbosity")

	rv.AddCommand(dump.Command())
	rv.AddCommand(gen.Command(fs, exit))
	rv.AddCommand(serve.Command(fs))

	return &rv
}

func main() {
	if err := Command().Execute(); err != nil {
		exit(1)
	}
	exit(0)
}
