package main

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"os"
)

type renderer interface {
	Render(io.Writer) error
}

var (
	ErrNoChange = errors.New("no change detected")
)

const DefaultFilePermissions os.FileMode = 0644

func main() {
	exitCode := flag.Bool("exit-code", false, "If specified, the exit code will be the number of files written plus 1. An exit code of 1 indicates program error.")
	enumOutput := flag.String("enum", "", "The output path of the generated enum code.")
	//envOutput := flag.String("env", "", "The output path of the generated environment code.")
	flag.Parse()

	written := 0

	if *enumOutput != "" {
		err := render(enumData, *enumOutput)
		if err != nil && err != ErrNoChange {
			panic(err)
		}
		if err == nil {
			written += 1
		}
	}

	if *exitCode && written > 0 {
		// add one because panic() always exits with 1, added value differentiates from a panic
		os.Exit(written + 1)
	}
}

func render(r renderer, filename string) error {
	previous, err := ioutil.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	buf := bytes.Buffer{}
	if err := r.Render(&buf); err != nil {
		return err
	}

	current := buf.Bytes()

	if bytes.Compare(previous, current) == 0 {
		return ErrNoChange
	}

	return ioutil.WriteFile(filename, current, DefaultFilePermissions)
}
