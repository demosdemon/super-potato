package api

import (
	"io"

	. "github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v2"

	"github.com/demosdemon/super-potato/pkg/gen"
	"github.com/demosdemon/super-potato/pkg/gen/variables"
)

type Collection []variables.WellKnownVariable

func NewCollection(r io.Reader) (gen.Renderer, error) {
	rv := Collection{}
	decoder := yaml.NewDecoder(r)
	err := decoder.Decode(&rv)
	return rv, err
}

func (l Collection) Render(w io.Writer) error {
	file := NewFile("serve")
	file.HeaderComment("This file is generated - do not edit!")
	file.Line()

	file.Func().Params(
		Id("x").Op("*").Id("API"),
	).Id("addGeneratedRoutes").Params().BlockFunc(func(g *Group) {
		for _, v := range l {
			name := "get" + v.Name
			file.Add(definition(v, name))
			g.Id("x").Dot("routes").Dot("GET").Call(
				Lit(strcase.ToSnake(v.Name)),
				Id("x").Dot(name),
			)
			for _, a := range v.Aliases {
				g.Id("x").Dot("routes").Dot("GET").Call(
					Lit(strcase.ToSnake(a)),
					Id("x").Dot(name),
				)
			}
		}
	}).Line()

	return file.Render(w)
}

func definition(v variables.WellKnownVariable, name string) Code {
	/*
		func (x *API) getApplication(c *gin.Context) {
			logrus.Trace("getApplication")
			obj, err := x.env.Application()
			_, ok := err.(platformsh.MissingEnvironment)
			switch {
			case err == nil:
				negotiate(c, http.StatusOK, obj)
			case ok:
				negotiate(c, http.StatusNotFound, err)
			default:
				negotiate(c, http.StatusInternalServerError, err)
			}
		}
	*/
	return Func().Params(
		Id("x").Op("*").Id("API"),
	).Id(name).Params(
		Id("c").Op("*").Qual("github.com/gin-gonic/gin", "Context"),
	).Block(
		Qual("github.com/sirupsen/logrus", "Trace").Call(
			Lit(name),
		),
		List(
			Id("obj"),
			Err(),
		).Op(":=").Id("x").Dot("env").Dot(v.Name).Call(),
		List(
			Id("_"),
			Id("ok"),
		).Op(":=").Err().Assert(Qual(
			"github.com/demosdemon/super-potato/pkg/platformsh",
			"MissingEnvironment",
		)),
		Switch().Block(
			Case(Err().Op("==").Nil()).Block(
				Id("negotiate").Call(
					Id("c"),
					Qual("net/http", "StatusOK"),
					Id("obj"),
				),
			),
			Case(Id("ok")).Block(
				Id("negotiate").Call(
					Id("c"),
					Qual("net/http", "StatusNotFound"),
					Err(),
				),
			),
			Default().Block(
				Id("negotiate").Call(
					Id("c"),
					Qual("net/http", "StatusInternalServerError"),
					Err(),
				),
			),
		),
	).Line()
}
