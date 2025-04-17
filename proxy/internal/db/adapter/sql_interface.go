package adapter

import (
	"database/sql"

	"studentgit.kata.academy/Zhodaran/go-kata/config"
	"studentgit.kata.academy/Zhodaran/go-kata/proxy/postgres"
)

type CreateResponse struct {
	Message string        `json:"message"`
	Books   []config.Book `json:"books"` // Добавляем поле для списка книг
}

type rErrorResponse struct {
	BadRequest      string `json:"400"`
	DadataBad       string `json:"500"`
	SuccefulRequest string `json:"200"`
}

type UserController struct {
	UserRepo postgres.UserRepository
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{Db: db}
}

type PostgresUserRepository struct {
	Db *sql.DB
}
