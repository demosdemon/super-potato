package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"os"

	. "github.com/dave/jennifer/jen"
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
	file := NewFile("platformsh")
	file.HeaderComment("This file is generated - do not edit!")
	file.Line()

	types := make([]Code, len(enums))
	consts := make([]Code, len(enums))
	slices := make([]Code, len(enums))
	maps := make([]Code, len(enums))
	methods := make([]Code, len(enums))

	for idx, enum := range enums {
		types[idx] = enum.TypeDefinition()
		consts[idx] = enum.ConstDefinition().Line()
		slices[idx] = enum.SliceDefinition().Line()
		maps[idx] = enum.MapDefinition().Line()
		methods[idx] = enum.MethodDefinitions().Line()
	}

	file.Type().Defs(types...).Line()
	file.Add(consts...)

	vars := make([]Code, 0, len(enums)*2)
	vars = append(vars, slices...)
	vars = append(vars, maps...)
	file.Var().Defs(vars...).Line()

	file.Add(methods...)

	return file.Render(w)
}
