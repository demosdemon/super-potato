package dump

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	stdout = os.Stdout
)

type formatter interface {
	format(w io.Writer, key, value string) error
}

type formatTemplate string

const (
	environFormat formatTemplate = "%s=%s\n"
	shellFormat   formatTemplate = "export %s=%q\n"
)

func (f formatTemplate) format(w io.Writer, key, value string) error {
	_, err := fmt.Fprintf(w, string(f), key, value)
	if err != nil {
		return err
	}

	return flush(w)
}

func newFormatter(arg string) formatter {
	switch arg {
	case "environ":
		return environFormat
	case "shell":
		return shellFormat
	default:
		return formatTemplate(arg)
	}
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

func Command() *cobra.Command {
	rv := cobra.Command{
		Use:  "dump",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			format, err := cmd.Flags().GetString("format")
			if err != nil {
				return err
			}

			f := newFormatter(format)
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
