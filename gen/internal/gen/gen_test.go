package gen

import (
	"io"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRenderer struct {
	data []byte
	err  error
}

func (t testRenderer) Render(w io.Writer) error {
	if t.data != nil {
		_, err := w.Write(t.data)
		if err != nil {
			return err
		}
	}
	return t.err
}

type invalidFS struct {
	afero.Fs
}

func (invalidFS) Open(name string) (afero.File, error) {
	return nil, assert.AnError
}

func TestRender(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)

	fs := afero.NewMemMapFs()
	err := afero.WriteFile(fs, "/tmp/test", []byte(`package main

func main() {}
`), 0644)
	require.NoError(t, err)

	type args struct {
		r        Renderer
		filename string
		fs       afero.Fs
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid file system",
			args: args{
				r:        testRenderer{},
				filename: "/tmp/test",
				fs:       invalidFS{},
			},
			wantErr: true,
		},
		{
			name: "invalid renderer",
			args: args{
				r:        testRenderer{err: assert.AnError},
				filename: "/tmp/test",
				fs:       afero.NewMemMapFs(),
			},
			wantErr: true,
		},
		{
			name: "invalid go file",
			args: args{
				r: testRenderer{
					data: []byte(`package go`),
				},
				filename: "/tmp/test",
				fs:       afero.NewMemMapFs(),
			},
			wantErr: true,
		},
		{
			name: "change",
			args: args{
				r: testRenderer{
					data: []byte(`package main

func main() {}
`),
				},
				filename: "/tmp/test",
				fs:       afero.NewMemMapFs(),
			},
			wantErr: false,
		},
		{
			name: "no change",
			args: args{
				r: testRenderer{
					data: []byte(`package main

func main() {}
`),
				},
				filename: "/tmp/test",
				fs:       fs,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Render(tt.args.r, tt.args.filename, tt.args.fs); (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRenderMap_Keys(t *testing.T) {
	nop := func(reader io.Reader) (Renderer, error) {
		return nil, nil
	}

	tests := []struct {
		name string
		m    RenderMap
		want []string
	}{
		{
			name: "nil",
			m:    nil,
			want: []string{},
		},
		{
			name: "one",
			m:    RenderMap{"one": nop},
			want: []string{"one"},
		},
		{
			name: "two",
			m:    RenderMap{"one": nop, "two": nop},
			want: []string{"one", "two"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Keys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RenderMap.Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenderMap_Usage(t *testing.T) {
	nop := func(reader io.Reader) (Renderer, error) {
		return nil, nil
	}

	tests := []struct {
		name string
		m    RenderMap
		want string
	}{
		{
			name: "nil",
			m:    nil,
			want: "Specify the template to execute ()",
		},
		{
			name: "one",
			m:    RenderMap{"one": nop},
			want: "Specify the template to execute (\"one\")",
		},
		{
			name: "two",
			m:    RenderMap{"one": nop, "two": nop},
			want: "Specify the template to execute (\"one\", \"two\")",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Usage(); got != tt.want {
				t.Errorf("RenderMap.Usage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenderMap_Register(t *testing.T) {
	nop := func(reader io.Reader) (Renderer, error) {
		return nil, nil
	}

	type args struct {
		name string
		fn   NewRenderer
	}
	tests := []struct {
		name string
		m    RenderMap
		args args
	}{
		{
			name: "nop",
			m:    make(RenderMap, 1),
			args: args{
				name: "nop",
				fn:   nop,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Register(tt.args.name, tt.args.fn)
		})
	}
}
