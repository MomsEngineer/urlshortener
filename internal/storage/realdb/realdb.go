package realdb

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewRealDB(dsn string) (*sql.DB, error) {
	return sql.Open("pgx", dsn)
}
