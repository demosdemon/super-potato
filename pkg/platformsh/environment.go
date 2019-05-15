package platformsh

import (
	"bufio"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Environment interface {
	EnvironmentAPI

	Prefix() string

	FileSystem() afero.Fs
	SetFileSystem(afero.Fs)

	SetLookupFunc(LookupFunc)
	Lookup(string) (string, bool)

	ReadDotEnv()

	Listener() (net.Listener, error)

	Variable(key string) (interface{}, bool)
}

type LookupFunc func(string) (string, bool)

type environment struct {
	prefix string

	fsMu sync.Mutex
	fs   afero.Fs

	lookupMu sync.Mutex
	lookup   LookupFunc

	readDotEnvOnce sync.Once
	dotEnv         map[string]string
}

func DefaultFileSystem(cwd string) afero.Fs {
	if cwd == "" {
		cwd, _ = os.Getwd()
	}

	fs := afero.NewOsFs()
	if cwd != "" {
		fs = afero.NewBasePathFs(fs, cwd)
	}

	return fs
}

func NewEnvironment(prefix string) Environment {
	return &environment{prefix: prefix}
}

func (e *environment) Prefix() string {
	return e.prefix
}

func (e *environment) FileSystem() afero.Fs {
	cwd, _ := os.Getwd()

	e.fsMu.Lock()
	if e.fs == nil {
		e.fs = DefaultFileSystem(cwd)
	}
	fs := e.fs
	e.fsMu.Unlock()
	return fs
}

func (e *environment) SetFileSystem(fs afero.Fs) {
	e.fsMu.Lock()
	e.fs = fs
	e.fsMu.Unlock()
}

func (e *environment) ReadDotEnv() {
	e.readDotEnvOnce.Do(e.reallyReadDotEnv)
}

func (e *environment) reallyReadDotEnv() {
	fs := e.FileSystem()
	fp, err := fs.Open("/.env")
	if os.IsNotExist(err) {
		logrus.Debug("no .env file found")
		return
	}
	e.dotEnv = make(map[string]string)
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.Index(line, "="); idx > 0 {
			k := line[:idx]
			v := line[idx+1:]
			logrus.WithFields(logrus.Fields{
				"key":   k,
				"value": v,
			}).Infof("read %s from .env", k)
			e.dotEnv[k] = v
		}
	}
	_ = fp.Close()
}

func (e *environment) Lookup(name string) (string, bool) {
	e.ReadDotEnv()
	if v, ok := e.dotEnv[name]; ok {
		return v, true
	}

	e.lookupMu.Lock()
	fn := e.lookup
	e.lookupMu.Unlock()

	if fn == nil {
		fn = os.LookupEnv
	}
	return fn(name)
}

func (e *environment) SetLookupFunc(fn LookupFunc) {
	e.lookupMu.Lock()
	e.lookup = fn
	e.lookupMu.Unlock()
}

func (e *environment) Listener() (net.Listener, error) {
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

func (e *environment) Variable(key string) (interface{}, bool) {
	vars, err := e.Variables()
	if err != nil {
		return nil, false
	}
	v, ok := vars[key]
	return v, ok
}
