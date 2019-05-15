package scrape

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/pkg/app"
)

func Command(app *app.App) *cobra.Command {
	rv := cobra.Command{
		Use:  "scrape URL",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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
			ctx, cancel := context.WithTimeout(app, time.Minute)
			list, err := page.Collect(ctx)
			cancel()

			m := make(map[string]struct{})
			for _, v := range list {
				if _, ok := m[v.Name]; !ok {
					fmt.Fprintln(app.Stdout, v.Name)
				}
				m[v.Name] = struct{}{}
			}

			return nil
		},
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
