package platformsh

import (
	"net"
	"os"

	"github.com/sirupsen/logrus"
)

var DefaultEnvironment = Environment{Prefix: "PLATFORM_"}

type LookupFunc func(string) (string, bool)

type Environment struct {
	Prefix string
	Lookup LookupFunc
}

// TODO: use jennifer to generate methods for all well known variables

func (e Environment) lookup(name string) (string, bool) {
	fn := e.Lookup
	if fn == nil {
		fn = os.LookupEnv
	}
	return fn(name)
}

func (e Environment) Listener() (net.Listener, error) {
	logrus.Trace("NewListener")

	var agg = make(AggregateError, 0, 2)

	socket, err := e.Socket()
	if err == nil {
		logrus.Debugf("found SOCKET=%q", socket)
		return net.Listen("unix", socket)
	}
	agg = agg.Append(err)

	port, err := e.Port()
	if err == nil {
		logrus.Debugf("found PORT=%q", port)
		addr := "127.0.0.1:" + port
		logrus.Debugf("listening on %q", addr)
		return net.Listen("tcp", addr)
	}
	agg = agg.Append(err)

	logrus.Debug("found neither SOCKET nor PORT")
	return nil, agg
}
