package repository

import "context"

// User представляет собой модель пользователя
type User struct {
	ID        int          `json:"id"`
	Name      string       `json:"name"`
	Email     string       `json:"email"`
	DeletedAt *string      `json:"deleted_at"` // Для логического удаления
	Books     map[int]Book `json:"books"`
}

type Book struct {
	Index     int    `json:"index"`
	Book      string `json:"book"`
	Author    string `json:"author"`
	Block     bool   `json:"block"`
	TakeCount int    `json:"take_count"`
}

// UserRepository определяет методы для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user User) error
	GetByID(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]User, error)
}

var Users []User
