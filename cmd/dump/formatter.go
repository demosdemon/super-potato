package dump

import (
	"fmt"
	"io"
	"net/http"
)

type formatter interface {
	format(w io.Writer, key, value string) error
}

type formatTemplate string

const (
	environFormat formatTemplate = "%s=%s\n"
	shellFormat   formatTemplate = "export %s=%q\n"
)

func (f formatTemplate) format(w io.Writer, key, value string) error {
	_, err := fmt.Fprintf(w, string(f), key, value)
	if err != nil {
		return err
	}

	return flush(w)
}

func newFormatter(arg string) formatter {
	switch arg {
	case "environ":
		return environFormat
	case "shell":
		return shellFormat
	default:
		return formatTemplate(arg)
	}
}

type flusher interface {
	Flush() error
}

func flush(w io.Writer) error {
	switch w := w.(type) {
	case http.Flusher:
		w.Flush()
	case flusher:
		return w.Flush()
	}
	return nil
}
