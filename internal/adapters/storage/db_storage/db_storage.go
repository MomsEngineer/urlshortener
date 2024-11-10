package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/MomsEngineer/urlshortener/internal/adapters/logger"
	ierrors "github.com/MomsEngineer/urlshortener/internal/errors"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var log = logger.Create()

type Database struct {
	sqlDB *sql.DB
	table string
}

func NewDB(dsn string) (*Database, error) {
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	table := "links"
	exist, err := tableExists(sqlDB, table)
	if err != nil {
		log.Error("Failed to check table existence", err)
		return nil, err
	}

	if !exist {
		if err := createTable(sqlDB, table); err != nil {
			log.Error("Failed to create table", err)
			return nil, err
		}
	}

	return &Database{sqlDB: sqlDB, table: table}, nil
}

func tableExists(sqlDB *sql.DB, tableName string) (bool, error) {
	query := `SELECT to_regclass($1)`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var result sql.NullString
	err := sqlDB.QueryRowContext(ctx, query, tableName).Scan(&result)
	if err != nil {
		return false, err
	}

	if !result.Valid {
		return false, nil
	}

	return true, nil
}

func createTable(sqlDB *sql.DB, table string) error {
	if !isValidTableName(table) {
		return errors.New("invalid table name")
	}

	query := `CREATE TABLE IF NOT EXISTS ` + table + `(
		id SERIAL PRIMARY KEY,
		short_link VARCHAR(255) NOT NULL,
		original_link TEXT NOT NULL UNIQUE
	);`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := sqlDB.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	log.Debug("Table", table, "created successfully")

	return nil
}

func isValidTableName(table string) bool {
	re := regexp.MustCompile(`^[a-zA-Z]+$`)
	return re.MatchString(table)
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

// TODO: use getShortLinkByOriginal
func (db *Database) SaveLink(ctx context.Context, short, original string) (string, error) {
	query := "INSERT INTO " + db.table + " (short_link, original_link) VALUES ($1, $2)"
	stmt, err := db.sqlDB.PrepareContext(ctx, query)
	if err != nil {
		log.Error("Failed to prepare statement", err)
		return "", err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, short, original)
	if err != nil {
		if strings.Contains(err.Error(), "(SQLSTATE 23505)") {
			log.Error("Error: Duplicate link "+original, err)

			oldShort, err := db.getShortLinkByOriginal(ctx, original)
			if err != nil {
				log.Error("Faild to get link "+original, err)
				return "", err
			}
			return oldShort, ierrors.ErrDuplicate
		}
		log.Error("Failed to insert record", err)
		return "", err
	}

	return "", nil
}

func (db *Database) SaveLinksBatch(ctx context.Context, links map[string]string) error {
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

	for short, original := range links {
		_, err := stmt.ExecContext(ctx, short, original)
		if err != nil {
			log.Error("Failed to execute statement", err)
			return err
		}
	}
	return tx.Commit()
}

func (db *Database) GetLink(ctx context.Context, shortLink string) (string, bool, error) {
	query := `SELECT original_link FROM ` + db.table + ` WHERE short_link = $1`
	stmt, err := db.sqlDB.PrepareContext(ctx, query)
	if err != nil {
		log.Error("Failed to prepare statement", err)
		return "", false, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, shortLink)

	var originalLink string
	err = row.Scan(&originalLink)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug("Not found original link for short link", shortLink)
			return "", false, nil
		}
		log.Error("Failed to scan response from DB", err)
		return "", false, err
	}

	return originalLink, true, nil
}

func (db *Database) Ping(ctx context.Context) error {
	return db.sqlDB.PingContext(ctx)
}

func (db *Database) Close() error {
	return db.sqlDB.Close()
}
