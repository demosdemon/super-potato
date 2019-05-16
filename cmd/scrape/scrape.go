package scrape

import (
	"fmt"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/pkg/app"
)

type Config struct {
	*app.App `flag:"-"`
	Output   string `flag:"output o" desc:"Where the output is written"`
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

func Command(app *app.App) *cobra.Command {
	cfg := Config{
		App:    app,
		Output: "-",
	}

	rv := cobra.Command{
		Use:  "scrape URL",
		Args: cobra.ExactArgs(1),
		RunE: cfg.Run,
	}

	err := gpflag.ParseTo(&cfg, rv.Flags())
	if err != nil {
		logrus.WithField("err", err).Fatal("failed to parse config flags")
	}

	return &rv
}

type ResponseError struct {
	StatusCode int
	Status     string
	Body       []byte
}

func (r ResponseError) Error() string {
	return fmt.Sprintf("HTTP %03d: %s - %s", r.StatusCode, r.Status, string(r.Body))
}

func NewResponseError(r *http.Response) error {
	body, _ := ioutil.ReadAll(r.Body)
	return ResponseError{
		StatusCode: r.StatusCode,
		Status:     r.Status,
		Body:       body,
	}
}
