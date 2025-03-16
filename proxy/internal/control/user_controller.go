package control

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/repository"
)

type UserController struct {
	userRepo repository.UserRepository
}

type BookController struct {
	DB *sql.DB
}

func NewUserController(userRepo repository.UserRepository) *UserController {
	return &UserController{userRepo: userRepo}
}

type CreateResponse struct {
	Message string            `json:"message"`
	Books   []repository.Book `json:"books"` // Добавляем поле для списка книг
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
	// Получаем список книг из базы данных
	books, err := uc.getBooksFromDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Устанавливаем статус 200 OK

	// Кодируем и отправляем список книг
	if err := json.NewEncoder(w).Encode(books); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (uc *BookController) getBooksFromDB() ([]repository.Book, error) {
	query := "SELECT index, book, author, block, take_count FROM book"
	rows, err := uc.DB.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []repository.Book
	for rows.Next() {
		var book repository.Book
		if err := rows.Scan(&book.Index, &book.Book, &book.Author, &book.Block, &book.TakeCount); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
