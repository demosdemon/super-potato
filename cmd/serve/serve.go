package serve

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/pkg/platformsh"
)

var env = platformsh.DefaultEnvironment

func Command(fs afero.Fs) *cobra.Command {
	rv := cobra.Command{
		Use:  "serve",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := env.Listener()
			if err != nil {
				return err
			}

			env.SetFileSystem(fs)

			engine := gin.New()
			engine.Use(
				gin.Logger(),
				gin.Recovery(),
			)

			_ = New(engine.Group("/api"), env)

			return http.Serve(l, engine)
		},
	}

	return &rv
}
