package enums

import (
	"io"

	. "github.com/dave/jennifer/jen"
	"github.com/go-openapi/inflect"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v2"

	"github.com/demosdemon/super-potato/gen/internal/gen"
)

type Collection []Enum

type Enum struct {
	Name   string      `yaml:"name"`
	Values []EnumValue `yaml:"values"`
}

type EnumValue struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func NewCollection(r io.Reader) (gen.Renderer, error) {
	rv := Collection{}
	decoder := yaml.NewDecoder(r)
	err := decoder.Decode(&rv)
	return rv, err
}

func (e Collection) Render(w io.Writer) error {
	file := NewFile("platformsh")
	file.HeaderComment("This file is generated - do not edit!")
	file.Line()

	types := make([]Code, len(e))
	consts := make([]Code, len(e))
	slices := make([]Code, len(e))
	maps := make([]Code, len(e))
	methods := make([]Code, len(e))

	for idx, enum := range e {
		types[idx] = enum.typeDefinition()
		consts[idx] = enum.constDefinition().Line()
		slices[idx] = enum.sliceDefinition().Line()
		maps[idx] = enum.mapDefinition().Line()
		methods[idx] = enum.methodDefinitions().Line()
	}

	file.Type().Defs(types...).Line()
	file.Add(consts...)

	vars := make([]Code, 0, len(e)*2)
	vars = append(vars, slices...)
	vars = append(vars, maps...)
	file.Var().Defs(vars...).Line()

	file.Add(methods...)

	return file.Render(w)
}

func (e Enum) len() int {
	return len(e.Values)
}

func (e Enum) totalName() string {
	return "total" + inflect.Pluralize(e.Name)
}

func (e Enum) sliceName() string {
	return strcase.ToLowerCamel(inflect.Pluralize(e.Name))
}

func (e Enum) mapName() string {
	return e.sliceName() + "Map"
}

func (e Enum) typeDefinition() *Statement {
	return Id(e.Name).Uint8()
}

func (e Enum) constDefinition() *Statement {
	defs := make([]Code, e.len()+1)

	for idx, v := range e.Values {
		name := e.Name + v.Name
		stmt := Id(name)

		if idx == 0 {
			stmt = stmt.Id(e.Name).Op("=").Iota()
		}

		defs[idx] = stmt
	}

	defs[e.len()] = Id(e.totalName())

	return Const().Defs(defs...).Line()
}

func (e Enum) sliceDefinition() *Statement {
	defs := make([]Code, e.len()+1)

	for idx, v := range e.Values {
		defs[idx] = Line().Lit(v.Value)
	}

	defs[e.len()] = Line()

	return Id(e.sliceName()).Op("=").Index(Id(e.totalName())).String().Values(defs...)
}

func (e Enum) mapDefinition() *Statement {
	defs := make([]Code, e.len()+1)

	for idx, v := range e.Values {
		defs[idx] = Line().Lit(v.Value).Op(":").Id(e.Name + v.Name)
	}

	defs[e.len()] = Line()

	return Id(e.mapName()).Op("=").Map(String()).Id(e.Name).Values(defs...)
}

func (e Enum) methodDefinitions() *Statement {
	rv := new(Statement)
	rv.Add(
		e.newEnumMethodDefinition().Line(),
		e.stringMethodDefinition().Line(),
		e.unmarshalTextMethodDefinition().Line(),
		e.marshalTextMethodDefinition().Line(),
	)
	return rv
}

func (e Enum) newEnumMethodDefinition() *Statement {
	/*
		func NewEnum(name string) (Enum, error) {
			if v, ok := enumMap[name]; ok {
				return v, nil
			}

			return 0, fmt.Errorf("unknown Enum name %q", name)
		}
	*/

	return Func().Id("New"+e.Name).Params(
		Id("name").String(),
	).Params(
		Id(e.Name),
		Error(),
	).Block(
		If(
			List(
				Id("v"),
				Id("ok"),
			).Op(":=").Id(e.mapName()).Index(Id("name")),
			Id("ok"),
		).Block(
			Return(
				Id("v"),
				Nil(),
			),
		),
		Line(),
		Return(
			Lit(0),
			Qual("fmt", "Errorf").Call(
				Lit("unknown "+e.Name+" name %q"),
				Id("name"),
			),
		),
	).Line()
}

func (e Enum) stringMethodDefinition() *Statement {
	/*
		func (v Enum) String() string {
			if v < totalEnums {
				return enums[v]
			}

			return fmt.Sprintf("unknown Enum value %02x", uint8(v))
		}
	*/

	return Func().Params(
		Id("v").Id(e.Name),
	).Id("String").Params().String().Block(
		If(
			Id("v").Op("<").Id(e.totalName()),
		).Block(
			Return(
				Id(e.sliceName()).Index(Id("v")),
			),
		),
		Line(),
		Return(
			Qual("fmt", "Sprintf").Call(
				Lit("unknown "+e.Name+" value %02x"),
				Uint8().Call(Id("v")),
			),
		),
	).Line()
}

func (e Enum) unmarshalTextMethodDefinition() *Statement {
	/*
		func (v *Enum) UnmarshalText(text []byte) (err error) {
			*v, err = NewEnum(string(text))
			return err
		}
	*/

	return Func().Params(
		Id("v").Op("*").Id(e.Name),
	).Id("UnmarshalText").Params(
		Id("text").Index().Byte(),
	).Params(
		Err().Error(),
	).Block(
		List(
			Op("*").Id("v"),
			Err(),
		).Op("=").Id("New"+e.Name).Call(
			String().Call(Id("text")),
		),
		Return(
			Err(),
		),
	).Line()
}

func (e Enum) marshalTextMethodDefinition() *Statement {
	/*
		func (v Enum) MarshalText() ([]byte, error) {
			return []byte(v.String()), nil
		}
	*/

	return Func().Params(
		Id("v").Id(e.Name),
	).Id("MarshalText").Params().Params(
		Index().Byte(),
		Error(),
	).Block(
		Return(
			Index().Byte().Call(Id("v").Dot("String").Call()),
			Nil(),
		),
	).Line()
}
