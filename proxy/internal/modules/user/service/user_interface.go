package service

import "database/sql"

type PostgresUserRepository struct {
	Db *sql.DB
}
