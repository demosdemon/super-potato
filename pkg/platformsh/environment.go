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

var DefaultEnvironment = &Environment{Prefix: "PLATFORM_"}

type LookupFunc func(string) (string, bool)

type Environment struct {
	Prefix string
	Lookup LookupFunc

	fsMu sync.Mutex
	fs   afero.Fs

	readDotEnvOnce sync.Once
	dotEnv         map[string]string
}

func (e *Environment) FileSystem() afero.Fs {
	cwd, _ := os.Getwd()

	e.fsMu.Lock()
	if e.fs == nil {
		fs := afero.NewOsFs()
		if cwd != "" {
			fs = afero.NewBasePathFs(fs, cwd)
		}
		e.fs = fs
	}
	fs := e.fs
	e.fsMu.Unlock()
	return fs
}

func (e *Environment) SetFileSystem(fs afero.Fs) {
	e.fsMu.Lock()
	e.fs = fs
	e.fsMu.Unlock()
}

func (e *Environment) ReadDotEnv() {
	e.readDotEnvOnce.Do(e.reallyReadDotEnv)
}

func (e *Environment) reallyReadDotEnv() {
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

func (e *Environment) lookup(name string) (string, bool) {
	e.ReadDotEnv()
	if v, ok := e.dotEnv[name]; ok {
		return v, true
	}

	fn := e.Lookup
	if fn == nil {
		fn = os.LookupEnv
	}
	return fn(name)
}

func (e *Environment) Listener() (net.Listener, error) {
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
