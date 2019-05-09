package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/pkg/gen"
	"github.com/demosdemon/super-potato/pkg/gen/api"
	"github.com/demosdemon/super-potato/pkg/gen/enums"
	"github.com/demosdemon/super-potato/pkg/gen/variables"
)

var (
	fs        = afero.NewOsFs()
	exit      = os.Exit
	renderMap = gen.RenderMap{
		"api":       api.NewCollection,
		"enums":     enums.NewCollection,
		"variables": variables.NewCollection,
	}
)

func Command() *cobra.Command {
	rv := cobra.Command{
		Use: "gen TEMPLATE INPUT [OUTPUT]",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("expected at least two arguments")
			}
			if len(args) > 3 {
				return errors.New("expected at most three arguments")
			}
			if _, ok := renderMap[args[0]]; !ok {
				return errors.New("invalid template name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			template := args[0]
			logrus.Tracef("template = %q", template)

			input, err := getInput(args[1])
			if err != nil {
				return err
			}
			defer input.Close()

			output := "/dev/stdout"
			if len(args) > 2 {
				output = args[2]
			}
			if output == "-" {
				output = "/dev/stdout"
			}
			logrus.Tracef("output = %v", output)

			flags := cmd.Flags()

			exitCode, err := flags.GetBool("exit-code")
			if err != nil {
				return err
			}
			logrus.Tracef("exitCode = %v", exitCode)

			renderer, err := renderMap[template](input)
			if err != nil {
				return err
			}
			logrus.Tracef("renderer = %v", renderer)

			err = gen.Render(renderer, output, fs)
			logrus.Tracef("render err = %v", err)

			if err != nil && err == gen.ErrNoChange {
				return nil
			}

			if err == nil && exitCode {
				exit(2)
			}

			return err
		},
	}

	flags := rv.Flags()
	flags.Bool("exit-code", false, gen.ExitCodeUsage)
	flags.String("output", "/dev/stdout", "Specify the output path.")

	return &rv
}

func getInput(s string) (io.ReadCloser, error) {
	logrus.Tracef("getInput(%q)", s)
	switch s {
	case "-", "/dev/stdin":
		return ioutil.NopCloser(os.Stdin), nil
	case "/dev/null":
		return ioutil.NopCloser(new(bytes.Buffer)), nil
	default:
		return fs.Open(s)
	}
}

func main() {
	err := Command().Execute()
	if err != nil {
		exit(1)
	}
}
