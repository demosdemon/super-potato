package secret

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	_ "github.com/octago/sflags/gen/gflag"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/pkg/app"
)

type Config struct {
	*app.App `flag:"-"`
	Format   string `flag:"format f" desc:"The output format; one of base16, base32, base64, PEM, blob" valid:"in(base16,base32,base64,PEM,blob)"`
	Output   string `flag:"output o" desc:"Where the output is written"`
	Bytes    int    `flag:"bytes b" desc:"The number of bytes to generate (n bits / 8 bits per byte)"`
}

func (c *Config) Run(cmd *cobra.Command, args []string) error {
	rv, err := RandBytes(c.Bytes)
	if err != nil {
		return err
	}

	n, err := c.Write(rv)
	if n < c.Bytes {
		return io.ErrShortWrite
	}

	return err
}

func (c *Config) Write(data []byte) (int, error) {
	fp, err := c.GetOutput(c.Output)
	if err != nil {
		return 0, err
	}
	defer fp.Close()

	var output string
	switch c.Format {
	case "base16":
		ss := make([]string, len(data))
		for idx, b := range data {
			ss[idx] = fmt.Sprintf("%02x", b)
		}
		output = strings.Join(ss, ":")
	case "base32":
		output = base32.StdEncoding.EncodeToString(data)
	case "base64":
		output = base64.StdEncoding.EncodeToString(data)
	case "PEM":
		block := pem.EncodeToMemory(
			&pem.Block{
				Type:  "RANDOM DATA",
				Bytes: data,
			},
		)
		output = string(block)
	case "blob":
		output = string(data)
	}

	n, err := fp.Write([]byte(output))
	if err != nil {
		return n, err
	}

	if istty(fp) {
		n2, err := fp.Write([]byte("\n"))
		n += n2
		return n, err
	}

	return n, nil
}

func RandBytes(count int) ([]byte, error) {
	rv := make([]byte, count)
	for read := 0; read < count; {
		n, err := rand.Read(rv[read:])
		logrus.WithFields(logrus.Fields{
			"read": read,
			"n":    n,
			"err":  err,
		}).Trace("rand.Read")
		read += n
		if err != nil {
			return nil, err
		}
	}
	return rv, nil
}

func Command(app *app.App) *cobra.Command {
	cfg := Config{
		App:    app,
		Format: "base32",
		Output: "-",
		Bytes:  40,
	}

	rv := cobra.Command{
		Use:  "secret",
		Args: cobra.NoArgs,
		RunE: cfg.Run,
	}

	err := gpflag.ParseTo(&cfg, rv.Flags())
	if err != nil {
		logrus.WithField("err", err).Fatal("failed to parse config flags")
	}

	return &rv
}

func istty(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		return isatty.IsTerminal(f.Fd())
	}
	return false
}
