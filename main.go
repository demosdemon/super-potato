//go:generate time go run ./gen

package main

import (
	"context"

	"github.com/demosdemon/super-potato/cmd/secret"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/cmd/dump"
	"github.com/demosdemon/super-potato/cmd/scrape"
	"github.com/demosdemon/super-potato/cmd/serve"
	"github.com/demosdemon/super-potato/pkg/app"
)

func Command(app *app.App) *cobra.Command {
	logrus.SetOutput(app.Stderr)

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
