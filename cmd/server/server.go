package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/demosdemon/super-potato/pkg/platformsh"
)

var env = platformsh.DefaultEnvironment

func Execute() {
	l, err := env.Listener()
	if err != nil {
		panic(err)
	}

	engine := gin.New()
	engine.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	api(engine.Group("/api"))

	_ = http.Serve(l, engine)
}
