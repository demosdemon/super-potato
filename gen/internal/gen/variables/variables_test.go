package variables

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
			name:    "empty",
			args:    args{new(bytes.Buffer)},
			want:    Collection{},
			wantErr: true,
		},
		{
			name:    "one",
			args:    args{bytes.NewReader([]byte("- name: One\n"))},
			want:    Collection{{Name: "One"}},
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

package platformsh

type EnvironmentAPI interface{}
`,
		},
		{
			name: "one",
			l:    Collection{{Name: "one"}},
			wantW: `// This file is generated - do not edit!

package platformsh

type EnvironmentAPI interface {
	One() (string, error)
}

func (e *environment) One() (string, error) {
	name := e.Prefix() + "ONE"
	value, ok := e.Lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}
`,
		},
		{
			name: "one no prefix",
			l:    Collection{{Name: "one", NoPrefix: true}},
			wantW: `// This file is generated - do not edit!

package platformsh

type EnvironmentAPI interface {
	One() (string, error)
}

func (e *environment) One() (string, error) {
	name := "ONE"
	value, ok := e.Lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}
`,
		},
		{
			name: "one alias",
			l:    Collection{{Name: "one", Aliases: []string{"two", "three"}}},
			wantW: `// This file is generated - do not edit!

package platformsh

type EnvironmentAPI interface {
	One() (string, error)
	Two() (string, error)
	Three() (string, error)
}

func (e *environment) One() (string, error) {
	name := e.Prefix() + "ONE"
	value, ok := e.Lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) Two() (string, error) {
	name := e.Prefix() + "ONE"
	value, ok := e.Lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) Three() (string, error) {
	name := e.Prefix() + "ONE"
	value, ok := e.Lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}
`,
		},
		{
			name: "one decoded type",
			l:    Collection{{Name: "one", DecodedType: "Decoded"}},
			wantW: `// This file is generated - do not edit!

package platformsh

import (
	"encoding/base64"
	"encoding/json"
)

type EnvironmentAPI interface {
	One() (Decoded, error)
}

func (e *environment) One() (Decoded, error) {
	name := e.Prefix() + "ONE"
	value, ok := e.Lookup(name)
	if !ok {
		return nil, missingEnvironment(name)
	}

	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	obj := Decoded{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	return obj, nil
}
`,
		},
		{
			name: "one decoded pointer",
			l:    Collection{{Name: "one", DecodedType: "Decoded", DecodedPointer: true}},
			wantW: `// This file is generated - do not edit!

package platformsh

import (
	"encoding/base64"
	"encoding/json"
)

type EnvironmentAPI interface {
	One() (*Decoded, error)
}

func (e *environment) One() (*Decoded, error) {
	name := e.Prefix() + "ONE"
	value, ok := e.Lookup(name)
	if !ok {
		return nil, missingEnvironment(name)
	}

	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	obj := Decoded{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	return &obj, nil
}
`,
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
