package env

import (
	"io"
	"sync"

	. "github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

type WellKnownVariables []WellKnownVariable

type WellKnownVariable struct {
	Name           string
	NoPrefix       bool
	Aliases        []string
	DecodedType    string
	DecodedPointer bool
}

func (l WellKnownVariables) Render(w io.Writer) error {
	file := NewFile("platformsh")
	file.HeaderComment("This file is generated - do not edit!")
	file.Line()

	for _, v := range l {
		file.Add(v.definition())
	}

	return file.Render(w)
}

func (v WellKnownVariable) names() []string {
	rv := make([]string, len(v.Aliases)+1)
	rv[0] = v.Name
	copy(rv[1:], v.Aliases)
	return apply(rv, strcase.ToCamel)
}

func (v WellKnownVariable) definition() *Statement {
	names := v.names()
	stmts := new(Statement)
	for _, name := range names {
		stmts.Add(v.function(name)).Line()
	}

	return stmts
}

func (v WellKnownVariable) function(name string) Code {
	/*
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
	*/

	return Func().Params(
		Id("e").Id("Environment"),
	).Id(name).Params().ParamsFunc(v.returnParams).BlockFunc(v.functionBlock).Line()
}

func (v WellKnownVariable) returnParams(g *Group) {
	// must be done one at a time even though the method is variadic
	g.Add(v.returnType())
	g.Add(Error())
}

func (v WellKnownVariable) functionBlock(g *Group) {
	// must be done one at a time even though the method is variadic
	g.Add(v.initName())
	g.Add(valueEquals())
	g.Add(v.ifNotOk())
	g.Add(Line())
	v.returnValue(g)
}

func (v WellKnownVariable) returnValue(g *Group) {
	if v.DecodedType == "" {
		g.Add(Return(Id("value"), Nil()))
		return
	}

	// must be done one at a time even though the method is variadic
	g.Add(decodeData())
	g.Add(v.ifErrNotNil())
	g.Add(Line())
	g.Add(v.initObj())
	g.Add(unmarshalObj())
	g.Add(v.ifErrNotNil())
	g.Add(Line())
	g.Add(Return(v.returnValueStmt(), Nil()))
}

func (v WellKnownVariable) returnType() Code {
	if v.DecodedType == "" {
		return String()
	}

	rv := new(Statement)
	if v.DecodedPointer {
		rv.Op("*")
	}

	return rv.Id(v.DecodedType)
}

func (v WellKnownVariable) returnValueStmt() Code {
	var rvStmt = new(Statement)
	if v.DecodedPointer {
		rvStmt.Op("&")
	}
	return rvStmt.Id("obj")
}

func (v WellKnownVariable) initName() Code {
	rv := Id("name").Op(":=")
	if !v.NoPrefix {
		rv.Id("e").Dot("Prefix").Op("+")
	}
	return rv.Lit(strcase.ToScreamingSnake(v.Name))
}

func (v WellKnownVariable) initObj() Code {
	return Id("obj").Op(":=").Id(v.DecodedType).Values()
}

func (v WellKnownVariable) ifNotOk() Code {
	return If(
		Op("!").Id("ok"),
	).Block(
		Return(
			v.zeroValue(),
			Id("missingEnvironment").Call(
				Id("name"),
			),
		),
	)
}

func (v WellKnownVariable) ifErrNotNil() Code {
	return If(
		Err().Op("!=").Nil(),
	).Block(
		Return(
			v.zeroValue(),
			Id("MissingEnvironment").Values(
				Id("name"),
				Err(),
			),
		),
	)
}

func (v WellKnownVariable) zeroValue() Code {
	if v.DecodedType == "" {
		return Lit("")
	}
	return Nil()
}

func unmarshalObj() Code {
	return Err().Op("=").Qual("encoding/json", "Unmarshal").Call(
		Id("data"),
		Op("&").Id("obj"),
	)
}

func decodeData() Code {
	return List(
		Id("data"),
		Err(),
	).Op(":=").Qual("encoding/base64", "StdEncoding").Dot("DecodeString").Call(
		Id("value"),
	)
}

func valueEquals() Code {
	return List(
		Id("value"),
		Id("ok"),
	).Op(":=").Id("e").Dot("lookup").Call(
		Id("name"),
	)
}

func apply(input []string, fn func(string) string) []string {
	output := make([]string, len(input))
	wg := sync.WaitGroup{}
	wg.Add(len(input))
	for idx, in := range input {
		go func(idx int, in string) {
			output[idx] = fn(in)
			wg.Done()
		}(idx, in)
	}
	wg.Wait()
	return output
}
