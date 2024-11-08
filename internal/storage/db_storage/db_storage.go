package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"time"

	"github.com/MomsEngineer/urlshortener/internal/logger"
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
		original_link TEXT NOT NULL
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

func (db *Database) SaveLink(shortLink, originalLink string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO links (short_link, original_link) VALUES ($1, $2)`
	_, err := db.sqlDB.ExecContext(ctx, query, shortLink, originalLink)
	if err != nil {
		log.Error("Failed to insert record", err)
		return err
	}

	return nil
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

func (db *Database) GetLink(shortLink string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT original_link FROM ` + db.table + ` WHERE short_link = $1`
	row := db.sqlDB.QueryRowContext(ctx, query, shortLink)

	var originalLink string
	err := row.Scan(&originalLink)
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

func (db *Database) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return db.sqlDB.PingContext(ctx)
}

func (db *Database) Close() error {
	return db.sqlDB.Close()
}
