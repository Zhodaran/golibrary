package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit"
	"studentgit.kata.academy/Zhodaran/go-kata/proxy/controller"
)

// @Summary Register a new user
// @Description This endpoint allows you to register a new user with a username and password.
// @Tags auth
// @Accept json
// @Produce json
// @Param user body User true "User registration details"
// @Success 201 {object} TokenResponse "User registered successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 409 {object} ErrorResponse "User already exists"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := Users[user.Username]; exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	Users[user.Username] = User{
		Username: user.Username,
		Password: user.Password,
	}

	// Используем логин пользователя в качестве user_id

}

// @Summary Login a user
// @Description This endpoint allows a user to log in with their username and password.
// @Tags auth
// @Produce json
// @Param user body User true "User login details"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Получаем данные пользователя из мапы Users
	storedUser, exists := Users[user.Username]
	if !exists {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Проверяем совпадение пароля
	if storedUser.Password != user.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Если авторизация успешна, создаем токен
	claims := map[string]interface{}{
		"user_id": user.Username, // Используем username как user_id
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	_, tokenString, err := TokenAuth.Encode(claims)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TokenResponse{Token: tokenString})
	fmt.Println(tokenString)
}

func GenerateUsers(count int) {
	for i := 0; i < count; i++ {
		username := gofakeit.Username()                                   // Генерация случайного имени пользователя
		password := gofakeit.Password(true, true, true, false, false, 10) // Генерация случайного пароля

		Users[username] = User{
			Username: username,
			Password: password,
		}
		fmt.Printf("Created user: %s with password: %s\n", username, password)
	}
}

// @Summary Get List of Registered Users
// @Description This endpoint returns a list of all registered users.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {array} User "List of registered users"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/users [get]
func ListUsersHandler(resp controller.Responder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		// Здесь предполагается, что Users - это глобальная переменная, содержащая всех пользователей
		var users []User
		for _, user := range Users {
			users = append(users, user)
		}

		resp.OutputJSON(w, users) // Возвращаем список пользователей
	}
}
