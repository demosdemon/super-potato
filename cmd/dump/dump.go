package dump

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/pkg/app"
)

type Config struct {
	*app.App `flag:"-"`
	Output   string `flag:"output o" desc:"Where the output is written"`
	Format   string `flag:"format f" desc:"The output format; one of environ, shell"`
}

func New(app *app.App) app.Config {
	return &Config{
		App:    app,
		Output: "-",
		Format: "environ",
	}
}

func (c *Config) Use() string {
	return "dump"
}

func (c *Config) Args(cmd *cobra.Command, args []string) error {
	return cobra.NoArgs(cmd, args)
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
