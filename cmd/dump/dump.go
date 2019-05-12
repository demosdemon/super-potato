package dump

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/pkg/app"
)

func Command(app *app.App) *cobra.Command {
	logrus.SetOutput(app.Stderr)

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
					if err := f.format(app.Stdout, k, v); err != nil {
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
