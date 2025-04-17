package config

import (
	"context"
	"encoding/json"
	"net/http"
)

type Book struct {
	Index     int    `json:"index"`
	Book      string `json:"book"`
	Author    string `json:"author"`
	Block     *bool  `json:"block"`
	TakeCount int    `json:"take_count"`
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

func (uc *BookController) getBooksFromDB() ([]Book, error) {
	query := "SELECT index, book, author, block, take_count FROM book"
	rows, err := uc.DB.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
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
