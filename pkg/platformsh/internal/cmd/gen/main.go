package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
)

const DefaultFilePermissions os.FileMode = 0644

func main() {
	enumOutput := flag.String("enum", "", "The output path of the generated enum code.")
	//envOutput := flag.String("env", "", "The output path of the generated environment code.")
	flag.Parse()

	if *enumOutput != "" {
		buf := bytes.Buffer{}
		if err := enumData.Render(&buf); err != nil {
			panic(err)
		}

		if err := ioutil.WriteFile(*enumOutput, buf.Bytes(), DefaultFilePermissions); err != nil {
			panic(err)
		}
	}
}
