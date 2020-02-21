package mysqlORM

import (
	"database/sql"
)

type ORMConfig struct {
	table       string
	db          *sql.DB
}

func NewORMConfig(table string, db *sql.DB) ORMConfig {
	return ORMConfig{
		table:                 table,
		db:                    db,
	}
}