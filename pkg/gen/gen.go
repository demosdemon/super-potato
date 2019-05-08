package gen

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/spf13/afero"
)

type (
	Renderer interface {
		Render(writer io.Writer) error
	}

	NewRenderer func(reader io.Reader) (Renderer, error)

	RenderMap map[string]NewRenderer
)

var (
	ErrNoChange = errors.New("no change detected")
)

const (
	DefaultFilePermissions os.FileMode = 0644
	ExitCodeUsage                      = "If specified, the exit code will be the number of files written plus 1. An exit code of 1 indicates program error."
)

func Render(r Renderer, filename string, fs afero.Fs) error {
	previous, err := afero.ReadFile(fs, filename)
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

	return afero.WriteFile(fs, filename, current, DefaultFilePermissions)
}

func (m RenderMap) Keys() []string {
	rv := make([]string, 0, len(m))
	for k := range m {
		rv = append(rv, k)
	}
	sort.Strings(rv)
	return rv
}

func (m RenderMap) Usage() string {
	keys := m.Keys()
	keys = Apply(keys, quote)
	keyString := strings.Join(keys, ", ")
	return fmt.Sprintf("Specify the template to execute (%s)", keyString)
}

func Apply(input []string, fn func(string) string) []string {
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

func quote(s string) string {
	return fmt.Sprintf("%q", s)
}
