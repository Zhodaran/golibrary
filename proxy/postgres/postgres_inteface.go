package postgres

import (
	"context"

	"studentgit.kata.academy/Zhodaran/go-kata/proxy/internal/models"
)

type ErrorResponse struct {
	Message string `json:"message"` // Сообщение об ошибке
	Code    int    `json:"code"`    // Код ошибки
}

type UserController struct {
	userRepo UserRepository
}

type UserRepository interface {
	Create(ctx context.Context, user models.User) error
	GetByID(ctx context.Context, id string) (models.User, error)
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]models.User, error)
}
