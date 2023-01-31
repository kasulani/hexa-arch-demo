package app

import (
	"context"
	"log"
	"time"

	"github.com/goava/di"
	"github.com/kelseyhightower/envconfig"

	"github.com/kasulani/hexa-arch-demo/internal/url"
)

type (
	config struct {
		LogLevel        string        `envconfig:"LOG_LEVEL" default:"debug"`
		Address         string        `envconfig:"SERVER_ADDRESS" default:"localhost:80"`
		Host            string        `envconfig:"URL_SHORTENER_HOST" default:"localhost"`
		Proto           string        `envconfig:"URL_SHORTENER_PROTO" default:"http"`
		Port            string        `envconfig:"URL_SHORTENER_PORT" default:"8080"`
		DSN             string        `envconfig:"DATABASE_DSN" required:"true"`
		MaxOpenConns    int           `envconfig:"DATABASE_MAX_OPEN_CONNS" default:"15"`
		MaxIdleConns    int           `envconfig:"DATABASE_MAX_IDLE_CONNS" default:"15"`
		ConnMaxLifetime time.Duration `envconfig:"DATABASE_CONN_MAX_LIFETIME" default:"10m"`
	}
)

func newConfig() *config {
	cfg := new(config)
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatalf("failed to load configuration: %q", err)
	}

	return cfg
}

// Container returns an instance of a DI container.
func Container() *di.Container {
	c, err := di.New(
		di.Provide(context.Background),
		di.Provide(newConfig),
		di.Provide(newDatabase),
		di.Provide(newRepository, di.As(new(url.Persister))),
		di.Provide(newRootCommand),
		di.Provide(newShortenCommand, di.As(new(SubCommand))),
		di.Invoke(registerSubCommands),
		di.Provide(newGenerator, di.As(new(url.CodeGenerator))),
		di.Provide(url.NewShortenURLUseCase, di.As(new(url.Shortener))),
	)
	if err != nil {
		log.Fatalf("failed to create DI container: %q", err)
	}

	return c
}

// Run is our application entry point.
func Run(root *rootCommand) error {
	return root.Execute()
}

// TerminateConnections will close connections to app dependencies.
func TerminateConnections(db *database) {
	err := db.Close()
	if err != nil {
		log.Fatalf("failed to close database connection: %q", err)
	}
}
