package dbstorage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/MomsEngineer/urlshortener/internal/adapters/logger"
	"github.com/MomsEngineer/urlshortener/internal/entities/link"
	ierror "github.com/MomsEngineer/urlshortener/internal/errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var log = logger.Create()

type Database struct {
	sqlDB *sql.DB
	table string
}

func NewDB(dsn string) (*Database, error) {
	table := "links"

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate driver, %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migration", table, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("failed to do migrate %w", err)
	}

	return &Database{sqlDB: sqlDB, table: table}, nil
}

func (db *Database) getShortLinkByOriginal(ctx context.Context, original string) (string, error) {
	query := `SELECT short_link FROM ` + db.table + ` WHERE original_link = $1`
	stmt, err := db.sqlDB.PrepareContext(ctx, query)
	if err != nil {
		log.Error("Failed to prepare statement", err)
		return "", err
	}
	defer stmt.Close()

	var short string
	row := stmt.QueryRowContext(ctx, original)
	err = row.Scan(&short)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug("Not found original link", original)
			return "", err
		}
		log.Error("Failed to scan response from DB", err)
		return "", err
	}

	return short, nil
}

func (db *Database) SaveLinksBatch(ctx context.Context, ls []*link.Link) error {
	tx, err := db.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		log.Error("Failed to create transaction", err)
		return err
	}
	defer tx.Rollback()

	query := "INSERT INTO " + db.table + " (short_link, original_link) VALUES($1, $2)"
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Error("Failed to prepare statement", err)
		return err
	}
	defer stmt.Close()

	for _, l := range ls {
		_, err := stmt.ExecContext(ctx, l.ShortURL, l.OriginalURL)
		if err != nil {
			log.Error("Failed to execute statement", err)
			return err
		}
	}
	return tx.Commit()
}

func (db *Database) SaveLink(ctx context.Context, l *link.Link) error {
	query := "INSERT INTO " + db.table + " (short_link, original_link) VALUES ($1, $2)"
	stmt, err := db.sqlDB.PrepareContext(ctx, query)
	if err != nil {
		log.Error("Failed to prepare statement", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, l.ShortURL, l.OriginalURL)
	if err != nil {
		if strings.Contains(err.Error(), "(SQLSTATE 23505)") {
			log.Error("Error: Duplicate link "+l.OriginalURL, err)

			oldShort, err := db.getShortLinkByOriginal(ctx, l.OriginalURL)
			if err != nil {
				log.Error("Faild to get link "+l.OriginalURL, err)
				return err
			}
			l.ShortURL = oldShort
			return ierror.ErrDuplicate
		}
		log.Error("Failed to insert record", err)
		return err
	}

	return nil
}

func (db *Database) GetLink(ctx context.Context, l *link.Link) error {
	query := `SELECT original_link FROM ` + db.table + ` WHERE short_link = $1`
	stmt, err := db.sqlDB.PrepareContext(ctx, query)
	if err != nil {
		log.Error("Failed to prepare statement", err)
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, l.ShortURL)

	err = row.Scan(&l.OriginalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug("Not found original link for short link", l.ShortURL)
			return nil
		}
		log.Error("Failed to scan response from DB", err)
		return err
	}

	return nil
}

func (db *Database) Ping(ctx context.Context) error {
	return db.sqlDB.PingContext(ctx)
}

func (db *Database) Close() error {
	return db.sqlDB.Close()
}
