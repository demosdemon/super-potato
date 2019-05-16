package deploy

import (
	"bitbucket.org/liamstask/goose/lib/goose"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/pkg/app"
)

const migrationsDir = "./vendor/github.com/cloudflare/cfssl/certdb/pg/migrations"

type Config struct {
	*app.App `flag:"-"`
}

func New(app *app.App) app.Config {
	return &Config{
		App: app,
	}
}

func (c *Config) Use() string {
	return "deploy"
}

func (c *Config) Args(cmd *cobra.Command, args []string) error {
	return cobra.NoArgs(cmd, args)
}

func (c *Config) Run(cmd *cobra.Command, args []string) error {
	rels, err := c.Relationships()
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
