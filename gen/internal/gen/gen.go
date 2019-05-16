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
	"github.com/sqs/goreturns/returns"
	"golang.org/x/tools/imports"
)

type (
	Renderer interface {
		Render(writer io.Writer) error
	}

	NewRenderer func(reader io.Reader) (Renderer, error)

	RenderMap map[string]NewRenderer
)

var (
	ErrNoChange      = errors.New("no change detected")
	DefaultRenderMap = make(RenderMap)
)

const (
	DefaultFilePermissions os.FileMode = 0644
	ExitCodeUsage                      = ""
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

	imports.LocalPrefix = "github.com/demosdemon"
	current, err = imports.Process(filename, current, &imports.Options{
		Comments:  true,
		TabIndent: true,
		TabWidth:  8,
	})
	if err != nil {
		return err
	}

	current, err = returns.Process("", filename, current, &returns.Options{
		RemoveBareReturns: true,
	})
	if err != nil {
		return err
	}

	if bytes.Equal(previous, current) {
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

func (m RenderMap) Register(name string, fn NewRenderer) {
	m[name] = fn
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
