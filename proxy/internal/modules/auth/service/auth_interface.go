package service

import (
	"sync"

	"github.com/go-chi/jwtauth"
)

type LoginResponse struct {
	Message string `json:"message"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	BadRequest      string `json:"400"`
	DadataBad       string `json:"500"`
	SuccefulRequest string `json:"200"`
}

var (
	TokenAuth = jwtauth.New("HS256", []byte("your_secret_key"), nil)
	Users     = make(map[string]User) // Хранение пользователей
	mu        sync.Mutex
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
