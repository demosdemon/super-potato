package main

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/pkg/app"
	"github.com/demosdemon/super-potato/pkg/gen"
	"github.com/demosdemon/super-potato/pkg/gen/api"
	"github.com/demosdemon/super-potato/pkg/gen/enums"
	"github.com/demosdemon/super-potato/pkg/gen/variables"
)

func main() {
	inst, cancel := app.New(context.Background())
	inst.Execute(Command)
	cancel()
}

var renderMap = gen.RenderMap{
	"api":       api.NewCollection,
	"enums":     enums.NewCollection,
	"variables": variables.NewCollection,
}

func Command(app *app.App) *cobra.Command {
	rv := cobra.Command{
		Use: "gen TEMPLATE INPUT OUTPUT",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return errors.New("expected exactly three arguments")
			}
			if _, ok := renderMap[args[0]]; !ok {
				return errors.New("invalid template name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			template := args[0]
			logrus.WithField("template", template).Trace()

			input, err := app.GetInput(args[1])
			if err != nil {
				return err
			}
			defer input.Close()

			output := args[2]
			logrus.WithField("output", output).Trace()

			flags := cmd.Flags()

			exitCode, err := flags.GetBool("exit-code")
			if err != nil {
				return err
			}
			logrus.WithField("exitCode", exitCode).Trace()

			renderer, err := renderMap[template](input)
			if err != nil {
				return err
			}
			logrus.WithField("renderer", renderer).Trace()

			err = gen.Render(renderer, output, app)
			logrus.WithField("err", err).Trace()

			if err != nil && err == gen.ErrNoChange {
				return nil
			}

			if err == nil && exitCode {
				app.Exit(2)
			}

			return err
		},
	}

	flags := rv.Flags()
	flags.Bool("exit-code", false, gen.ExitCodeUsage)

	return &rv
}
