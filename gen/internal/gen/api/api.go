package api

import (
	"io"

	. "github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v2"

	"github.com/demosdemon/super-potato/gen/internal/gen"
	"github.com/demosdemon/super-potato/gen/internal/gen/variables"
)

const (
	logrusPath     = "github.com/sirupsen/logrus"
	platformshPath = "github.com/demosdemon/super-potato/pkg/platformsh"
	ginPath        = "github.com/gin-gonic/gin"
	httpPath       = "net/http"
)

func init() {
	gen.DefaultRenderMap.Register("api", NewCollection)
}

type Collection []variables.WellKnownVariable

func NewCollection(r io.Reader) (gen.Renderer, error) {
	rv := Collection{}
	decoder := yaml.NewDecoder(r)
	err := decoder.Decode(&rv)
	return rv, err
}

func (l Collection) Render(w io.Writer) error {
	file := NewFile("server")
	file.HeaderComment("This file is generated - do not edit!")
	file.Line()

	/*
		func (s *Server) registerGeneratedRoutes(r gin.IRoutes) {
			r.GET("application", s.getApplication)
			r.GET("app", s.getApplication)
		}
	*/

	file.Func().Params(receiver()).Id("registerGeneratedRoutes").Params(
		Id("r").Qual(ginPath, "IRoutes"),
	).BlockFunc(func(g *Group) {
		for _, v := range l {
			name := "get" + v.Name
			file.Add(getterDefinition(v, name))
			g.Id("r").Dot("GET").Call(
				Lit(strcase.ToSnake(v.Name)),
				Id("s").Dot(name),
			)
			for _, a := range v.Aliases {
				g.Id("r").Dot("GET").Call(
					Lit(strcase.ToSnake(a)),
					Id("s").Dot(name),
				)
			}
		}
	}).Line()

	return file.Render(w)
}

func getterDefinition(v variables.WellKnownVariable, name string) Code {
	/*
		func (s *Server) getApplication(c *gin.Context) {
			logrus.Trace("getApplication")
			obj, err := s.Environment.Application()
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
	*/
	return Func().Params(receiver()).Id(name).Params(contextParam()).Block(
		Qual(logrusPath, "Trace").Call(Lit(name)),
		List(
			Id("obj"),
			Err(),
		).Op(":=").Id("s").Dot("Environment").Dot(v.Name).Call(),
		List(
			Id("_"),
			Id("ok"),
		).Op(":=").Err().Assert(Qual(platformshPath, "MissingEnvironment")),
		Switch().Block(
			Case(Err().Op("==").Nil()).Block(negotiate("StatusOK", Id("obj"))),
			Case(Id("ok")).Block(negotiate("StatusNotFound", Err())),
			Default().Block(negotiate("StatusInternalServerError", Err())),
		),
	).Line()
}

func negotiate(status string, result Code) Code {
	return Id("s").Dot("negotiate").Call(
		Id("c"),
		Qual(httpPath, status),
		result,
	)
}

func receiver() Code {
	return Id("s").Op("*").Id("Server")
}

func contextParam() Code {
	return Id("c").Op("*").Qual(ginPath, "Context")
}
