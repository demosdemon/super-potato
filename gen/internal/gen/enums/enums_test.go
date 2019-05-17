package enums

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
			want:    Collection{},
			wantErr: true,
		},
		{
			name: "one",
			args: args{
				r: bytes.NewReader([]byte("- name: One\n")),
			},
			want: Collection{
				{
					Name: "One",
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
		e       Collection
		wantW   string
		wantErr bool
	}{
		{
			name: "empty",
			e:    nil,
			wantW: `// This file is generated - do not edit!

package platformsh

type ()

var ()
`,
			wantErr: false,
		},
		{
			name: "no values",
			e: Collection{
				{
					Name: "Invalid",
				},
			},
			wantW:   ``,
			wantErr: true,
		},
		{
			name: "One value",
			e: Collection{
				{
					Name: "Type",
					Values: []EnumValue{
						{
							Name:  "Value",
							Value: "value",
						},
					},
				},
			},
			wantW: `// This file is generated - do not edit!

package platformsh

import "fmt"

type (
	Type uint8
)

const (
	TypeValue Type = iota
	totalTypes
)

var (
	types = [totalTypes]string{
		"value",
	}

	typesMap = map[string]Type{
		"value": TypeValue,
	}
)

func NewType(name string) (Type, error) {
	if v, ok := typesMap[name]; ok {
		return v, nil
	}

	return 0, fmt.Errorf("unknown Type name %q", name)
}

func (v Type) String() string {
	if v < totalTypes {
		return types[v]
	}

	return fmt.Sprintf("unknown Type value %02x", uint8(v))
}

func (v *Type) UnmarshalText(text []byte) (err error) {
	*v, err = NewType(string(text))
	return err
}

func (v Type) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := tt.e.Render(w); (err != nil) != tt.wantErr {
				t.Errorf("Collection.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Collection.Render() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
