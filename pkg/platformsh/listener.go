package platformsh

import (
	"fmt"
	"net"
	"os"

	"github.com/sirupsen/logrus"
)

func NewListener() (net.Listener, error) {
	logrus.Trace("NewListener")

	if socket, ok := os.LookupEnv("SOCKET"); ok {
		logrus.Debugf("found SOCKET=%q", socket)
		return net.Listen("unix", socket)
	}

	if port, ok := os.LookupEnv("PORT"); ok {
		logrus.Debugf("found PORT=%q", port)
		addr := fmt.Sprintf("127.0.0.1:%s", port)
		logrus.Tracef("listening on %q", addr)
		return net.Listen("tcp", addr)
	}

	logrus.Debug("found neither SOCKET nor PORT")
	return nil, missingEnvironment("SOCKET", "PORT")
}
