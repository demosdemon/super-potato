package deploy

import (
	"sync"

	"bitbucket.org/liamstask/goose/lib/goose"
	_ "github.com/cloudflare/cfssl/certdb/sql"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/pkg/app"
	"github.com/demosdemon/super-potato/pkg/platformsh"
)

const migrationsDir = "./vendor/github.com/cloudflare/cfssl/certdb/pg/migrations"

type Config struct {
	*app.App `flag:"-"`
	Prefix   string `flag:"prefix p" desc:"The environment prefix for Platform.sh"`

	envMu sync.Mutex
	env   platformsh.Environment
}

func (c *Config) Environment() platformsh.Environment {
	c.envMu.Lock()
	if c.env == nil {
		c.env = platformsh.NewEnvironment(c.Prefix)
	}
	env := c.env
	c.envMu.Unlock()
	return env
}

func (c *Config) Run(cmd *cobra.Command, args []string) error {
	rels, err := c.Environment().Relationships()
	if err != nil {
		return errors.Wrap(err, "unable to locate relationships")
	}

	dbOpen, err := rels.Postgresql("database")
	if err != nil {
		return errors.Wrap(err, "unable to get database connection string")
	}

	conf := goose.DBConf{
		MigrationsDir: migrationsDir,
		Env:           "production",
		Driver: goose.DBDriver{
			Name:    "postgres",
			OpenStr: dbOpen,
			Import:  "github.com/lib/pq",
			Dialect: &goose.PostgresDialect{},
		},
	}

	target, err := goose.GetMostRecentDBVersion(conf.MigrationsDir)
	if err != nil {
		return errors.Wrap(err, "unable to determine latest database version")
	}

	err = goose.RunMigrations(&conf, conf.MigrationsDir, target)
	if err != nil {
		return errors.Wrap(err, "unable to run migrations")
	}

	return nil
}

func Command(app *app.App) *cobra.Command {
	cfg := Config{
		App:    app,
		Prefix: "PLATFORM_",
	}

	rv := cobra.Command{
		Use:  "deploy",
		Args: cobra.NoArgs,
		RunE: cfg.Run,
	}

	err := gpflag.ParseTo(&cfg, rv.Flags())
	if err != nil {
		logrus.WithField("err", err).Fatal("failed to parse config flags")
	}

	return &rv
}
