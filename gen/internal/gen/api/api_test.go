package api

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/demosdemon/super-potato/gen/internal/gen"
)

func TestNewCollection(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    gen.Renderer
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				r: new(bytes.Buffer),
			},
			want:    make(Collection, 0),
			wantErr: true,
		},
		{
			name: "one",
			args: args{
				r: bytes.NewReader([]byte("- name: one\n")),
			},
			want: Collection{
				{
					Name: "one",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCollection(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollection_Render(t *testing.T) {
	tests := []struct {
		name    string
		l       Collection
		wantW   string
		wantErr bool
	}{
		{
			name: "empty",
			l:    nil,
			wantW: `// This file is generated - do not edit!

package server

import gin "github.com/gin-gonic/gin"

func (s *Server) registerGeneratedRoutes(r gin.IRoutes) {}
`,
			wantErr: false,
		},
		{
			name: "one",
			l: Collection{
				{
					Name: "one",
				},
			},
			wantW: `// This file is generated - do not edit!

package server

import (
	platformsh "github.com/demosdemon/super-potato/pkg/platformsh"
	gin "github.com/gin-gonic/gin"
	logrus "github.com/sirupsen/logrus"
	"net/http"
)

func (s *Server) registerGeneratedRoutes(r gin.IRoutes) {
	r.GET("one", s.getone)
}

func (s *Server) getone(c *gin.Context) {
	logrus.Trace("getone")
	obj, err := s.Environment.one()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}
`,
			wantErr: false,
		},
		{
			name: "one alias",
			l: Collection{
				{
					Name:    "One",
					Aliases: []string{"Two", "Three"},
				},
			},
			wantW: `// This file is generated - do not edit!

package server

import (
	platformsh "github.com/demosdemon/super-potato/pkg/platformsh"
	gin "github.com/gin-gonic/gin"
	logrus "github.com/sirupsen/logrus"
	"net/http"
)

func (s *Server) registerGeneratedRoutes(r gin.IRoutes) {
	r.GET("one", s.getOne)
	r.GET("two", s.getOne)
	r.GET("three", s.getOne)
}

func (s *Server) getOne(c *gin.Context) {
	logrus.Trace("getOne")
	obj, err := s.Environment.One()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := tt.l.Render(w); (err != nil) != tt.wantErr {
				t.Errorf("Collection.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Collection.Render() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
