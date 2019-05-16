package scrape

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/pkg/app"
)

type Config struct {
	*app.App `flag:"-"`
	Output   string `flag:"output o" desc:"Where the output is written"`
}

func New(app *app.App) app.Config {
	return &Config{
		App:    app,
		Output: "-",
	}
}

func (c *Config) Use() string {
	return "scrape URL"
}

func (c *Config) Args(cmd *cobra.Command, args []string) error {
	return cobra.ExactArgs(1)(cmd, args)
}

func (c *Config) Run(cmd *cobra.Command, args []string) error {
	rootURL := args[0]

	resp, err := http.Get(rootURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return NewResponseError(resp)
	}

	page := NewCharacterPage(resp.Body)
	list, err := page.Collect(c)

	fp, err := c.GetOutput(c.Output)
	if err != nil {
		return err
	}

	defer fp.Close()

	m := make(map[string]struct{})
	for _, v := range list {
		if _, ok := m[v.Name]; !ok {
			fmt.Fprintln(fp, v.Name)
		}
		m[v.Name] = struct{}{}
	}

	return nil
}
