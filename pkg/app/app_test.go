package app

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var token = struct{}{}

func recoverToken() {
	if r := recover(); r != nil && r != token {
		panic(r)
	}
}

type rootConfig struct {
	*App `flag:"-"`
	Name string
}

func (r rootConfig) Use() string {
	return r.Name
}

func (rootConfig) Args(*cobra.Command, []string) error {
	return nil
}

func (r rootConfig) SubCommands() []Config {
	return []Config{
		&config{
			App: r.App,
		},
	}
}

func (rootConfig) PersistentPreRun(*cobra.Command, []string) error {
	return nil
}

func (rootConfig) PersistentPostRun(*cobra.Command, []string) error {
	return nil
}

func (rootConfig) PreRun(*cobra.Command, []string) error {
	return nil
}

func (rootConfig) PostRun(*cobra.Command, []string) error {
	return nil
}

type config struct {
	*App `flag:"-"`
	err  error
}

func (config) Use() string {
	return "test"
}

func (config) Args(cmd *cobra.Command, args []string) error {
	return nil
}

func (c config) Run(cmd *cobra.Command, args []string) error {
	return c.err
}

func Test_setReportCaller(t *testing.T) {
	caller := os.Getenv("PKI_LOG_CALLER")
	if caller == "1" || caller == "true" {
		_ = os.Setenv("PKI_LOG_CALLER", "false")
	} else {
		_ = os.Setenv("PKI_LOG_CALLER", "true")
	}
	setReportCaller()
}

func TestApp_SetPrefix(t *testing.T) {
	app, cancel := New(context.Background())
	cancel()
	assert.Nil(t, app.Environment)
	app.SetPrefix("")
	assert.NotNil(t, app.Environment)
}

func TestApp_Logger(t *testing.T) {
	app, cancel := New(context.Background())
	cancel()

	var buf bytes.Buffer
	app.Stderr = &buf

	l, ok := app.Logger("").(*log.Logger)
	assert.True(t, ok)
	l.SetFlags(0)
	l.Print("test")

	assert.Equal(t, buf.Bytes(), []byte("test\n"))
}

func TestApp_Execute(t *testing.T) {
	tests := []struct {
		name string
		err  error
		exit int
	}{
		{
			name: "NoError",
			err:  nil,
			exit: 0,
		},
		{
			name: "Error",
			err:  assert.AnError,
			exit: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app, cancel := New(context.Background())
			cancel()
			defer recoverToken()

			app.Exit = func(exit int) {
				assert.Equal(t, tt.exit, exit)
				panic(token)
			}
			cfg := config{
				App: app,
				err: tt.err,
			}
			app.Execute(&cfg)
		})
	}
}

func TestApp_GetInput(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		read []byte
	}{
		{
			name: "stdin",
			arg:  "/dev/stdin",
			read: []byte(`this is the stdin`),
		},
		{
			name: "devnull",
			arg:  "/dev/null",
			read: []byte{},
		},
		{
			name: "file",
			arg:  "/tmp/file",
			read: []byte(`this is the file`),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app, cancel := New(context.Background())
			cancel()
			app.Fs = afero.NewMemMapFs()
			err := afero.WriteFile(app, "/tmp/file", []byte(`this is the file`), 0644)
			require.NoError(t, err)
			app.Stdin = bytes.NewReader([]byte(`this is the stdin`))
			fp, err := app.GetInput(tt.arg)
			require.NoError(t, err)
			read, err := ioutil.ReadAll(fp)
			assert.NoError(t, err)
			err = fp.Close()
			assert.NoError(t, err)
			assert.Equal(t, tt.read, read)
		})
	}
}

func TestApp_GetOutput(t *testing.T) {
	tests := []struct {
		name      string
		arg       string
		nopWriter bool
		stdout    []byte
		stderr    []byte
		file      []byte
	}{
		{
			name:      "stdout",
			arg:       "/dev/stdout",
			nopWriter: true,
			stdout:    []byte(`this is the outputthis is the output`),
			stderr:    nil,
			file:      nil,
		},
		{
			name:      "stderr",
			arg:       "/dev/stderr",
			nopWriter: true,
			stdout:    nil,
			stderr:    []byte(`this is the outputthis is the output`),
			file:      nil,
		},
		{
			name:      "devnull",
			arg:       "/dev/null",
			nopWriter: true,
			stdout:    nil,
			stderr:    nil,
			file:      nil,
		},
		{
			name:      "file",
			arg:       "/tmp/file",
			nopWriter: false,
			stdout:    nil,
			stderr:    nil,
			file:      []byte(`this is the output`),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app, cancel := New(context.Background())
			cancel()
			stdout := new(bytes.Buffer)
			stderr := new(bytes.Buffer)
			app.Stdout, app.Stderr = stdout, stderr
			app.Fs = afero.NewMemMapFs()

			fp, err := app.GetOutput(tt.arg)
			require.NoError(t, err)
			_, ok := fp.(NopWriterCloser)
			assert.Equal(t, tt.nopWriter, ok)
			_, err = fp.Write([]byte(`this is the output`))
			assert.NoError(t, err)
			err = fp.Close()
			assert.NoError(t, err)

			fp, err = app.GetOutput(tt.arg)
			require.NoError(t, err)
			_, ok = fp.(NopWriterCloser)
			assert.Equal(t, tt.nopWriter, ok)
			_, err = fp.Write([]byte(`this is the output`))
			assert.NoError(t, err)
			err = fp.Close()
			assert.NoError(t, err)

			file, err := afero.ReadFile(app, "/tmp/file")
			if err != nil && !os.IsNotExist(err) {
				require.FailNow(t, err.Error())
			}
			assert.Equal(t, tt.stdout, stdout.Bytes())
			assert.Equal(t, tt.stderr, stderr.Bytes())
			assert.Equal(t, tt.file, file)

		})
	}
}

func TestApp_Append(t *testing.T) {
	tests := []struct {
		name      string
		arg       string
		nopWriter bool
		stdout    []byte
		stderr    []byte
		file      []byte
	}{
		{
			name:      "stdout",
			arg:       "/dev/stdout",
			nopWriter: true,
			stdout:    []byte(`this is the outputthis is the output`),
			stderr:    nil,
			file:      nil,
		},
		{
			name:      "stderr",
			arg:       "/dev/stderr",
			nopWriter: true,
			stdout:    nil,
			stderr:    []byte(`this is the outputthis is the output`),
			file:      nil,
		},
		{
			name:      "devnull",
			arg:       "/dev/null",
			nopWriter: true,
			stdout:    nil,
			stderr:    nil,
			file:      nil,
		},
		{
			name:      "file",
			arg:       "/tmp/file",
			nopWriter: false,
			stdout:    nil,
			stderr:    nil,
			file:      []byte(`this is the outputthis is the output`),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app, cancel := New(context.Background())
			cancel()
			stdout := new(bytes.Buffer)
			stderr := new(bytes.Buffer)
			app.Stdout, app.Stderr = stdout, stderr
			app.Fs = afero.NewMemMapFs()

			fp, err := app.Append(tt.arg)
			require.NoError(t, err)
			_, ok := fp.(NopWriterCloser)
			assert.Equal(t, tt.nopWriter, ok)
			_, err = fp.Write([]byte(`this is the output`))
			assert.NoError(t, err)
			err = fp.Close()
			assert.NoError(t, err)

			fp, err = app.Append(tt.arg)
			require.NoError(t, err)
			_, ok = fp.(NopWriterCloser)
			assert.Equal(t, tt.nopWriter, ok)
			_, err = fp.Write([]byte(`this is the output`))
			assert.NoError(t, err)
			err = fp.Close()
			assert.NoError(t, err)

			file, err := afero.ReadFile(app, "/tmp/file")
			if err != nil && !os.IsNotExist(err) {
				require.FailNow(t, err.Error())
			}
			assert.Equal(t, tt.stdout, stdout.Bytes())
			assert.Equal(t, tt.stderr, stderr.Bytes())
			assert.Equal(t, tt.file, file)

		})
	}
}

func TestApp_ReadYAML(t *testing.T) {
	type data struct {
		Line  string `yaml:"line"`
		Block string `yaml:"block"`
	}

	app, cancel := New(context.Background())
	cancel()
	app.Fs = afero.NewMemMapFs()
	fp, err := app.GetOutput("/tmp/file")
	require.NoError(t, err)
	_, err = fp.Write([]byte(`---
line: This is a line
block: |
    This is a block
`))
	require.NoError(t, err)
	err = fp.Close()
	require.NoError(t, err)

	var d data
	err = app.ReadYAML("/tmp/badfile", &d)
	assert.Error(t, err)
	err = app.ReadYAML("/tmp/file", &d)
	assert.NoError(t, err)
	assert.Equal(t, "This is a line", d.Line)
	assert.Equal(t, "This is a block\n", d.Block)
}

func TestApp_command(t *testing.T) {
	app, cancel := New(context.Background())
	cancel()
	cfg := rootConfig{App: app}
	assert.Panics(t, func() {
		_ = app.command(cfg)
	})
	cmd := app.command(&cfg)
	assert.NotNil(t, cmd)
}
