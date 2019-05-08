package platformsh

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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

func (e Environment) Application() (*Application, error) {
	name := e.Prefix + "APPLICATION"
	value, ok := e.lookup(name)
	if !ok {
		return nil, missingEnvironment(name)
	}

	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	app := Application{}
	err = json.Unmarshal(data, &app)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	return &app, nil
}

func (e Environment) Routes() (Routes, error) {
	name := e.Prefix + "ROUTES"
	value, ok := e.lookup(name)
	if !ok {
		return nil, missingEnvironment(name)
	}

	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	routes := Routes{}
	err = json.Unmarshal(data, &routes)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	return routes, nil
}

func (e Environment) Relationships() (Relationships, error) {
	name := e.Prefix + "RELATIONSHIPS"
	value, ok := e.lookup(name)
	if !ok {
		return nil, missingEnvironment(name)
	}

	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	rels := Relationships{}
	err = json.Unmarshal(data, &rels)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	return rels, nil
}

func (e Environment) Variables() (JSONObject, error) {
	name := e.Prefix + "VARIABLES"
	value, ok := e.lookup(name)
	if !ok {
		return nil, missingEnvironment(name)
	}

	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	obj := JSONObject{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	return obj, nil
}

func (e Environment) Listener() (net.Listener, error) {
	logrus.Trace("NewListener")

	if socket, ok := e.lookup("SOCKET"); ok {
		logrus.Debugf("found SOCKET=%q", socket)
		return net.Listen("unix", socket)
	}

	if port, ok := e.lookup("PORT"); ok {
		logrus.Debugf("found PORT=%q", port)
		addr := fmt.Sprintf("127.0.0.1:%s", port)
		logrus.Debugf("listening on %q", addr)
		return net.Listen("tcp", addr)
	}

	logrus.Debug("found neither SOCKET nor PORT")
	return nil, missingEnvironment("SOCKET", "PORT")
}
