package data

import "database/sql"

type ServerContext struct {
	DB *sql.DB
}
