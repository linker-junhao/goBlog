package container

import "database/sql"

type MyContainer struct {
	MysqlDB *sql.DB
}