package control

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/repository"
)

type UserController struct {
	userRepo repository.UserRepository
}

type BookController struct {
	BookRepo repository.BookRepository
}

func NewUserController(userRepo repository.UserRepository) *UserController {
	return &UserController{userRepo: userRepo}
}

func NewBookController(BookRepo repository.BookRepository) *BookController {
	return &BookController{BookRepo: BookRepo}
}

type CreateResponse struct {
	Message string `json:"message"`
}

type rErrorResponse struct {
	BadRequest      string `json:"400"`
	DadataBad       string `json:"500"`
	SuccefulRequest string `json:"200"`
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user repository.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := uc.userRepo.Create(context.Background(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(CreateResponse{Message: "Create successful"})
	w.WriteHeader(http.StatusCreated)
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := uc.userRepo.GetByID(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(CreateResponse{Message: "Great successful"})
	json.NewEncoder(w).Encode(user)
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user repository.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := uc.userRepo.Update(context.Background(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(CreateResponse{Message: "Update successful"})
	w.WriteHeader(http.StatusNoContent)
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := uc.userRepo.Delete(context.Background(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(CreateResponse{Message: "Delete successful"})
	w.WriteHeader(http.StatusNoContent)
}

func (uc *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit := 10 // Установите значение по умолчанию
	offset := 0 // Установите значение по умолчанию
	users, err := uc.userRepo.List(context.Background(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(CreateResponse{Message: "List successful"})
	json.NewEncoder(w).Encode(users)
}

// @Summary List SQL book
// @Description This description created new SQL user
// @Tags Books
// @Accept json
// @Produce json
// @Success 200 {object} CreateResponse "List successful"
// @Failure 400 {object} rErrorResponse "Invalid request"
// @Failure 401 {object} rErrorResponse "Invalid credentials"
// @Failure 500 {object} rErrorResponse "Internal server error"
// @Router /api/books [get]
func (uc *BookController) ListBook(w http.ResponseWriter, r *http.Request) {
	limit := 100 // Установите значение по умолчанию
	offset := 0  // Установите значение по умолчанию
	books, err := uc.BookRepo.MList(context.Background(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(CreateResponse{Message: "List successful"})
	json.NewEncoder(w).Encode(books)
}
