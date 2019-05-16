package serve

import (
	"github.com/octago/sflags/gen/gpflag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/pkg/app"
	"github.com/demosdemon/super-potato/pkg/server"
)

type Config struct {
	*app.App      `flag:"-"`
	Prefix        string `flag:"prefix p" desc:"The Platform.sh environment prefix"`
	SessionCookie string `flag:"session-cookie" desc:"The cookie name used for session storage"`
}

func (c *Config) Run(cmd *cobra.Command, args []string) error {
	return server.New(c, c.Prefix, c.SessionCookie).Serve(nil)
}

func Command(app *app.App) *cobra.Command {
	cfg := Config{
		App:           app,
		Prefix:        "PLATFORM_",
		SessionCookie: "super-potato",
	}

	rv := cobra.Command{
		Use:  "serve",
		Args: cobra.NoArgs,
		RunE: cfg.Run,
	}

	err := gpflag.ParseTo(&cfg, rv.Flags())
	if err != nil {
		logrus.WithField("err", err).Fatal("failed to parse config flags")
	}

	return &rv
}
