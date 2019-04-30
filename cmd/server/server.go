package server

import (
	"github.com/demosdemon/super-potato/pkg/platformsh"
	"net/http"
)

func Execute() {
	l, err := platformsh.NewListener()
	if err != nil {
		panic(err)
	}

	_ = http.Serve(l, nil)
}