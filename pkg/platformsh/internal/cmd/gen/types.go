package main

import (
	. "github.com/dave/jennifer/jen"
	"github.com/go-openapi/inflect"
	"github.com/iancoleman/strcase"
)

type Enum struct {
	Name   string
	Values []EnumValue
}

type EnumValue struct {
	Name  string
	Value string
}

func (e Enum) Len() int {
	return len(e.Values)
}

func (e Enum) TotalName() string {
	return "total" + inflect.Pluralize(e.Name)
}

func (e Enum) SliceName() string {
	return strcase.ToLowerCamel(inflect.Pluralize(e.Name))
}

func (e Enum) MapName() string {
	return e.SliceName() + "Map"
}

func (e Enum) TypeDefinition() *Statement {
	return Id(e.Name).Uint8()
}

func (e Enum) ConstDefinition() *Statement {
	defs := make([]Code, e.Len()+1)

	for idx, v := range e.Values {
		name := e.Name + v.Name
		stmt := Id(name)

		if idx == 0 {
			stmt = stmt.Id(e.Name).Op("=").Iota()
		}

		defs[idx] = stmt
	}

	defs[e.Len()] = Id(e.TotalName())

	return Const().Defs(defs...).Line()
}

func (e Enum) SliceDefinition() *Statement {
	defs := make([]Code, e.Len()+1)

	for idx, v := range e.Values {
		defs[idx] = Line().Lit(v.Value)
	}

	defs[e.Len()] = Line()

	return Id(e.SliceName()).Op("=").Index(Id(e.TotalName())).String().Values(defs...)
}

func (e Enum) MapDefinition() *Statement {
	defs := make([]Code, e.Len()+1)

	for idx, v := range e.Values {
		defs[idx] = Line().Lit(v.Value).Op(":").Id(e.Name + v.Name)
	}

	defs[e.Len()] = Line()

	return Id(e.MapName()).Op("=").Map(String()).Id(e.Name).Values(defs...)
}

func (e Enum) MethodDefinitions() *Statement {
	rv := make(Statement, 0, 4)
	rv = append(
		rv,
		e.newEnumMethodDefinition().Line(),
		e.stringMethodDefinition().Line(),
		e.unmarshalTextMethodDefinition().Line(),
		e.marshalTextMethodDefinition().Line(),
	)
	return &rv
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
			).Op(":=").Id(e.MapName()).Index(Id("name")),
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
			Id("v").Op("<").Id(e.TotalName()),
		).Block(
			Return(
				Id(e.SliceName()).Index(Id("v")),
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
