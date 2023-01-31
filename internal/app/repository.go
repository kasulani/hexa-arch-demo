package app

import (
	"context"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Blank import to load and register the PostgreSQL driver.

	"github.com/kasulani/hexa-arch-demo/internal/url"
)

type (
	database struct {
		*sqlx.DB
	}

	shortURLRow struct {
		RawURL string `db:"raw_url"`
		Code   string `db:"code"`
	}

	repository struct {
		db *database
	}
)

const driver = "postgres"

func newDatabase(cfg *config) *database {
	db, err := sql.Open(driver, cfg.DSN)
	if err != nil {
		log.Fatalf("failed to open newDatabase connection: %q", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping newDatabase: %q", err)
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return &database{sqlx.NewDb(db, driver)}
}

func newRepository(db *database) *repository {
	return &repository{db: db}
}

// Persist will save short url to the database.
func (r *repository) Persist(ctx context.Context, url *url.ShortURL) (err error) {
	const query = `INSERT INTO url (raw_url, code) VALUES (:raw_url, :code)`

	data := shortURLRow{RawURL: url.Raw(), Code: url.Code()}

	_, err = r.db.NamedExecContext(ctx, query, data)

	return err
}
