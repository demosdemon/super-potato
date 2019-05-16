package serve

import (
	"github.com/demosdemon/super-potato/pkg/app"
	"github.com/demosdemon/super-potato/pkg/server"
)

func New(app *app.App) app.Config {
	return &server.Server{
		App:           app,
		SessionCookie: "super-potato",
	}
}
