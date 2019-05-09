package dump

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	stdout = os.Stdout
)

type formatter interface {
	format(w io.Writer, key, value string) error
}

type environFormatter struct{}
type shellFormatter struct{}

func newFormatter(arg string) formatter {
	switch arg {
	case "environ":
		return environFormatter{}
	case "shell":
		return shellFormatter{}
	default:
		logrus.WithField("arg", arg).Panic("invalid formatter")
	}
	panic(nil)
}

type flusher interface {
	Flush() error
}

func flush(w io.Writer) error {
	switch w := w.(type) {
	case http.Flusher:
		w.Flush()
	case flusher:
		return w.Flush()
	}
	return nil
}

func (environFormatter) format(w io.Writer, key, value string) error {
	_, err := fmt.Fprintf(w, "%s=%s\n", key, value)
	if err != nil {
		return err
	}

	return flush(w)
}

func (shellFormatter) format(w io.Writer, key, value string) error {
	_, err := fmt.Fprintf(w, "export %s=%q\n", key, value)
	if err != nil {
		return err
	}

	return flush(w)
}

func Command() *cobra.Command {
	rv := cobra.Command{
		Use:  "dump",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			formatFlag := cmd.Flag("format")
			if formatFlag == nil {
				return errors.New("format flag not defined")
			}

			f := newFormatter(formatFlag.Value.String())
			for _, line := range os.Environ() {
				if idx := strings.Index(line, "="); idx > 0 {
					k := line[:idx]
					v := line[idx+1:]
					if err := f.format(stdout, k, v); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}

	flags := rv.Flags()
	flags.StringP("format", "f", "environ", "Specify the output format. (environ: bare environ format as read by the API, shell: shell export commands)")

	return &rv
}
