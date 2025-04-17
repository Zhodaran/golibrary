package models

import (
	"studentgit.kata.academy/Zhodaran/go-kata/config"
)

type User struct {
	ID        int                 `json:"id"`
	Name      string              `json:"name"`
	Email     string              `json:"email"`
	DeletedAt *string             `json:"deleted_at"` // Для логического удаления
	Books     map[int]config.Book `json:"books"`
}
