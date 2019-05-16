package dump

import (
	"os"
	"strings"

	"github.com/octago/sflags/gen/gpflag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/pkg/app"
)

type Config struct {
	*app.App `flag:"-"`
	Output   string `flag:"output o" desc:"Where the output is written"`
	Format   string `flag:"format f" desc:"The output format; one of environ, shell"`
}

func (c *Config) Run(cmd *cobra.Command, args []string) error {
	fp, err := c.GetOutput(c.Output)
	if err != nil {
		return err
	}

	defer fp.Close()

	f := newFormatter(c.Format)
	for _, line := range os.Environ() {
		if idx := strings.Index(line, "="); idx > 0 {
			k := line[:idx]
			v := line[idx+1:]
			if err := f.format(fp, k, v); err != nil {
				return err
			}
		}
	}

	return nil
}

func Command(app *app.App) *cobra.Command {
	logrus.SetOutput(app.Stderr)

	cfg := Config{
		App:    app,
		Output: "-",
		Format: "environ",
	}

	rv := cobra.Command{
		Use:  "dump",
		Args: cobra.NoArgs,
		RunE: cfg.Run,
	}

	err := gpflag.ParseTo(&cfg, rv.Flags())
	if err != nil {
		logrus.WithField("err", err).Fatal("failed to parse config flags")
	}

	return &rv
}
