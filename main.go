//go:generate time go run ./cmd/gen enums ./data/enums.yaml ./pkg/platformsh/enums_gen.go
//go:generate time go run ./cmd/gen variables ./data/variables.yaml ./pkg/platformsh/environment_gen.go
//go:generate time go run ./cmd/gen api ./data/variables.yaml ./cmd/server/generated.go

package main

import (
	"flag"

	"github.com/sirupsen/logrus"

	"github.com/demosdemon/super-potato/cmd/server"
)

func init() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)
}

var (
	logLevel = flag.String("log-level", "trace", "control the logging verbosity")
)

func main() {
	flag.Parse()

	var level logrus.Level

	switch ll := *logLevel; ll {
	case "trace":
		level = logrus.TraceLevel
	case "debug":
		level = logrus.DebugLevel
	case "info":
		level = logrus.InfoLevel
	case "warn", "warning":
		level = logrus.WarnLevel
	case "error":
		level = logrus.ErrorLevel
	case "fatal":
		level = logrus.FatalLevel
	case "panic":
		level = logrus.PanicLevel
	default:
		logrus.Panicf("unknown log level %s", ll)
	}

	logrus.SetLevel(level)

	server.Execute()
}
