package models

// User представляет собой модель пользователя

type AddaderBook struct {
	Book   string `json:"book"`
	Author string `json:"author"`
}

// UserRepository определяет методы для работы с пользователями

var Users []User
