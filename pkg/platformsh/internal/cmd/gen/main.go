package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"os"

	"github.com/dave/jennifer/jen"
	"github.com/go-openapi/inflect"
	"github.com/iancoleman/strcase"
)

const DefaultFilePermissions os.FileMode = 0644

func main() {
	output := flag.String("output", "/dev/stdout", "The output path of the generated file.")
	flag.Parse()

	buf := bytes.Buffer{}
	if err := render(&buf); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(*output, buf.Bytes(), DefaultFilePermissions); err != nil {
		panic(err)
	}
}

func render(w io.Writer) error {
	file := jen.NewFile("platformsh")
	file.HeaderComment("This file is generated - do not edit!")
	file.Line()

	for _, enum := range enums {
		var (
			name        = enum.Name
			pluralName  = inflect.Pluralize(name)
			totalName   = "total" + strcase.ToCamel(pluralName)
			sliceName   = strcase.ToLowerCamel(pluralName)
			mapName     = sliceName + "Map"
			newEnumName = "New" + name

			num         = len(enum.Values)
			defs        = make([]jen.Code, num+1)
			sliceValues = make([]jen.Code, num+1)
			mapValues   = make([]jen.Code, num+1)
		)

		file.Type().Id(name).Uint8().Line()

		for idx, v := range enum.Values {
			var (
				enumName = name + v.Name
				stmt     = jen.Id(enumName)
			)

			if idx == 0 {
				stmt = stmt.Id(name).Op("=").Iota()
			}

			defs[idx] = stmt
			sliceValues[idx] = jen.Line().Lit(v.Value)
			mapValues[idx] = jen.Line().Lit(v.Value).Op(":").Id(name + v.Name)
		}

		defs[num] = jen.Id(totalName).Line()
		file.Const().Defs(defs...).Line()

		sliceValues[num] = jen.Line()
		file.Var().Id(sliceName).Op("=").Index(jen.Id(totalName)).String().Values(sliceValues...).Line()

		mapValues[num] = jen.Line()
		file.Var().Id(mapName).Op("=").Map(jen.String()).Id(name).Values(mapValues...).Line()

		/*
			func NewEnum(name string) (Enum, error) {
				if v, ok := enumMap[name]; ok {
					return v, nil
				}

				return 0, fmt.Errorf("unknown Enum name %q", name)
			}
		*/

		file.Func().Id(newEnumName).Params(
			jen.Id("name").String(),
		).Params(
			jen.Id(name),
			jen.Error(),
		).Block(
			jen.If(
				jen.List(jen.Id("v"), jen.Id("ok")).Op(":=").Id(mapName).Index(jen.Id("name")),
				jen.Id("ok"),
			).Block(
				jen.Return(
					jen.Id("v"),
					jen.Nil(),
				),
			),
			jen.Line(),
			jen.Return(
				jen.Lit(0),
				jen.Qual("fmt", "Errorf").Call(jen.Lit("unknown "+name+" name %q"), jen.Id("name")),
			),
		).Line()

		/*
			func (v Enum) String() string {
				if v < totalEnums {
					return enums[v]
				}

				return fmt.Sprintf("unknown Enum value %02x", uint8(v))
			}
		*/

		file.Func().Params(
			jen.Id("v").Id(name),
		).Id("String").Params().String().Block(
			jen.If(jen.Id("v").Op("<").Id(totalName)).Block(
				jen.Return(jen.Id(sliceName).Index(jen.Id("v"))),
			),
			jen.Line(),
			jen.Return(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("unknown "+name+" value %02x"),
					jen.Uint8().Call(jen.Id("v")),
				),
			),
		).Line()

		/*
			func (v *Enum) UnmarshalText(text []byte) (err error) {
				*v, err = NewEnum(string(text))
				return err
			}
		*/

		file.Func().Params(
			jen.Id("v").Op("*").Id(name),
		).Id("UnmarshalText").Params(
			jen.Id("text").Index().Byte(),
		).Params(
			jen.Err().Error(),
		).Block(
			jen.List(
				jen.Op("*").Id("v"),
				jen.Err(),
			).Op("=").Id(newEnumName).Call(
				jen.String().Call(jen.Id("text")),
			),
			jen.Return(jen.Err()),
		).Line()

		/*
			func (v Enum) MarshalText() ([]byte, error) {
				return []byte(v.String()), nil
			}
		*/

		file.Func().Params(
			jen.Id("v").Id(name),
		).Id("MarshalText").Params().Params(
			jen.Index().Byte(),
			jen.Error(),
		).Block(
			jen.Return(jen.Index().Byte().Call(jen.Id("v").Dot("String").Call()), jen.Nil()),
		)
	}

	return file.Render(w)
}
