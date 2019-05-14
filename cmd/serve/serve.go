package serve

import (
	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/pkg/app"
	"github.com/demosdemon/super-potato/pkg/server"
)

func Command(app *app.App) *cobra.Command {
	rv := cobra.Command{
		Use:  "serve",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.New(app).Serve(nil)
		},
	}

	return &rv
}
