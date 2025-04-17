package config

import (
	"database/sql"
	"net/http"
)

type BookController struct {
	DB *sql.DB
}

type Server struct {
	http.Server
}
