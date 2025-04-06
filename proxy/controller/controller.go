package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/repository"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type mErrorResponse struct {
	BadRequest      string `json:"400"`
	DadataBad       string `json:"500"`
	SuccefulRequest string `json:"200"`
}

type GeocodeRequest struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type RequestAddressSearch struct {
	Query string `json:"query"`
}

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{})

	ErrorUnauthorized(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorForbidden(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
}

func (l *Library) AddBook(book repository.Book) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Получаем список книг автора
	booksByAuthor := l.Books[book.Author]

	// Находим первый свободный индекс
	newIndex := 1
	for {
		found := false
		for _, b := range booksByAuthor {
			if b.Index == newIndex {
				found = true
				break
			}
		}
		if !found {
			break
		}
		newIndex++
	}

	// Присваиваем книге новый индекс
	book.Index = newIndex

	// Добавляем книгу в список
	l.Books[book.Author] = append(booksByAuthor, book)

	// Добавляем автора в список, если его там еще нет
	if !contains(l.Authors, book.Author) {
		l.Authors = append(l.Authors, book.Author)
	}
}

func contains(authors []string, author string) bool {
	for _, a := range authors {
		if a == author {
			return true
		}
	}
	return false
}

type Library struct {
	Books   map[string][]repository.Book
	Authors []string
	mu      sync.RWMutex
}

func NewLibrary() *Library {
	return &Library{
		Books: make(map[string][]repository.Book),
	}
}

type Respond struct {
	log *zap.Logger
}

func (l *Library) AddBooks(books []repository.Book) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, book := range books {
		// Добавляем книгу в мапу по автору
		l.Books[book.Author] = append(l.Books[book.Author], book)

		// Добавляем автора в список, если его там еще нет
		if !contains(l.Authors, book.Author) {
			l.Authors = append(l.Authors, book.Author)
		}
	}
}

func NewResponder(logger *zap.Logger) Responder {
	return &Respond{log: logger}
}

func (r *Respond) OutputJSON(w http.ResponseWriter, responseData interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		r.log.Error("responder json encode error", zap.Error(err))
	}
}

func (r *Respond) ErrorBadRequest(w http.ResponseWriter, err error) {
	r.log.Info("http response bad request status code", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}); err != nil {
		r.log.Info("response writer error on write", zap.Error(err))
	}
}

func (r *Respond) ErrorForbidden(w http.ResponseWriter, err error) {
	r.log.Warn("http resposne forbidden", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	if err := json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}); err != nil {
		r.log.Error("response writer error on write", zap.Error(err))
	}
}

func (r *Respond) ErrorUnauthorized(w http.ResponseWriter, err error) {
	r.log.Warn("http resposne Unauthorized", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}); err != nil {
		r.log.Error("response writer error on write", zap.Error(err))
	}
}

func (r *Respond) ErrorInternal(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) {
		return
	}
	r.log.Error("http response internal error", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}); err != nil {
		r.log.Error("response writer error on write", zap.Error(err))
	}
}

type TakeBookRequest struct {
	Username string `json:"username"` // Поле для имени пользователя
}

// @Summary Get Geo Coordinates by Address
// @Description This endpoint allows you to get geo coordinates by address.
// @Tags User
// @Accept json
// @Produce json
// @Param index path int true "Book INDEX"
// @Param Authorization header string true "Bearer Token"
// @Param body body TakeBookRequest true "Request body"
// @Success 200 {object} Response "Успешное выполнение"
// @Failure 400 {object} mErrorResponse "Ошибка запроса"
// @Failure 500 {object} mErrorResponse "Ошибка подключения к серверу"
// @Security BearerAuth
// @Router /api/book/take/{index} [post]
func TakeBookHandler(resp Responder, db *sql.DB, Books *[]repository.Book, library *Library) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexStr := chi.URLParam(r, "index")
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			resp.ErrorBadRequest(w, errors.New("invalid index"))
			return
		}

		// Обновление записи в таблице book
		result, err := db.Exec("UPDATE book SET block = $1, take_count = take_count + 1 WHERE index = $2 AND block = $3", true, index, false)
		if err != nil {
			resp.ErrorInternal(w, err)
			return
		}

		var bookFind repository.Book
		found := false

		// Поиск книги по индексу
		for i, book := range *Books {
			if index == book.Index {
				bookFind = book
				// Удаление книги из массива
				*Books = append((*Books)[:i], (*Books)[i+1:]...)
				found = true
				break
			}
		}

		if !found {
			http.Error(w, fmt.Sprintf("book with index %d not found", index), http.StatusNotFound)
			return
		}

		// Проверка, была ли книга успешно обновлена
		if rowsAffected, err := result.RowsAffected(); err != nil || rowsAffected == 0 {
			resp.ErrorBadRequest(w, errors.New("book not found or already taken"))
			return
		}

		var requestBody TakeBookRequest
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			resp.ErrorBadRequest(w, errors.New("invalid request body"))
			return
		}

		// Проверка, был ли передан username
		if requestBody.Username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}

		// Добавление книги к пользователю
		library.Books[requestBody.Username] = append(library.Books[requestBody.Username], bookFind)
		resp.OutputJSON(w, map[string]string{"message": "Book taken successfully"})
	}
}

// @Summary Get Geo Coordinates by Address
// @Description This endpoint allows you to get geo coordinates by address.
// @Tags User
// @Accept json
// @Produce json
// @Param index path int true "Book INDEX"
// @Param Authorization header string true "Bearer Token"
// @Param body body TakeBookRequest true "Request body"
// @Success 200 {object} Response "Успешное выполнение"
// @Failure 400 {object} mErrorResponse "Ошибка запроса"
// @Failure 500 {object} mErrorResponse "Ошибка подключения к серверу"
// @Security BearerAuth
// @Router /api/book/return/{index} [delete]
func ReturnBook(resp Responder, db *sql.DB, Books *[]repository.Book, library *Library) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexStr := chi.URLParam(r, "index")
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			resp.ErrorBadRequest(w, errors.New("invalid index"))
			return
		}

		var requestBody TakeBookRequest
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			resp.ErrorBadRequest(w, errors.New("invalid request body"))
			return
		}

		if requestBody.Username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}

		userBooks, userExists := library.Books[requestBody.Username]
		if !userExists {
			http.Error(w, "User has no books", http.StatusNotFound)
			return
		}

		found := false
		var bookFind repository.Book

		// Поиск книги у пользователя
		for i, book := range userBooks {
			if book.Index == index {
				bookFind = book
				library.Books[requestBody.Username] = append(userBooks[:i], userBooks[i+1:]...) // Удаляем книгу из списка пользователя
				found = true
				break
			}
		}

		if !found {
			http.Error(w, fmt.Sprintf("book with index %d not found for user", index), http.StatusNotFound)
			return
		}

		// Обновление записи в таблице book
		result, err := db.Exec("UPDATE book SET block = $1 WHERE index = $2 AND block = $3", false, index, true)
		if err != nil {
			resp.ErrorInternal(w, err)
			return
		}
		if rowsAffected, err := result.RowsAffected(); err != nil || rowsAffected == 0 {
			resp.ErrorBadRequest(w, errors.New("book not found or already returned"))
			return
		}

		// Добавление книги обратно в общий список книг
		*Books = append(*Books, bookFind) // Добавляем книгу обратно в общий список
		resp.OutputJSON(w, map[string]string{"message": "Book returned successfully"})
	}
}

// @Summary Обновление информации о книге
// @Description Этот эндпоинт позволяет обновить информацию о книге по индексу.
// @Tags Books
// @Accept json
// @Produce json
// @Param index path int true "Индекс книги"
// @Param Authorization header string true "Bearer Token"
// @Param body body repository.Book true "Обновленная информация о книге"
// @Success 200 {object} repository.Book "Успешное обновление книги"
// @Failure 400 {object} mErrorResponse "Ошибка запроса"
// @Failure 404 {object} mErrorResponse "Книга не найдена"
// @Failure 500 {object} mErrorResponse "Ошибка сервера"
// @Router /api/book/{index} [put]
func UpdateBook(resp Responder, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			return
		}

		indexStr := chi.URLParam(r, "index")
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			resp.ErrorBadRequest(w, errors.New("недопустимый индекс"))
			return
		}

		var updatedBook repository.Book
		if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
			resp.ErrorBadRequest(w, errors.New("недопустимый формат данных"))
			return
		}

		// Обновление записи в таблице book
		result, err := db.Exec("UPDATE book SET book = $1, author = $2, block = $3 WHERE index = $4",
			updatedBook.Book, updatedBook.Author, updatedBook.Block, index)
		if err != nil {
			resp.ErrorInternal(w, err)
			return
		}

		// Проверка, была ли книга успешно обновлена
		if rowsAffected, err := result.RowsAffected(); err != nil || rowsAffected == 0 {
			resp.ErrorBadRequest(w, errors.New("книга не найдена или не обновлена"))
			return
		}

		// Возвращаем обновленную книгу
		resp.OutputJSON(w, updatedBook)
	}
}

func ListAuthorsHandler(resp Responder, library *Library) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		library.mu.RLock()         // Блокируем чтение
		defer library.mu.RUnlock() // Разблокируем чтение после завершения

		authorsSet := make(map[string]struct{}) // Используем множество для уникальных авторов

		// Проходим по всем книгам в библиотеке и собираем авторов
		for _, books := range library.Books {
			for _, book := range books {
				authorsSet[book.Author] = struct{}{} // Добавляем автора в множество
			}
		}

		// Преобразуем множество в срез
		var authors []string
		for author := range authorsSet {
			authors = append(authors, author)
		}

		resp.OutputJSON(w, authors) // Возвращаем список авторов в формате JSON
	}
}

// @Summary Add a new book to the library
// @Description This endpoint allows you to add a new book to the library.
// @Tags Books
// @Accept json
// @Produce json
// @Param book body repository.AddaderBook false "Book details"
// @Success 201 {object} repository.Book "Book added successfully"
// @Failure 400 {object} mErrorResponse "Invalid request"
// @Failure 500 {object} mErrorResponse "Internal server error"
// @Router /api/book [post]
func AddBookHandler(resp Responder, db *sql.DB, library *Library, Books *[]repository.Book) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var addaderBook repository.AddaderBook
		if err := json.NewDecoder(r.Body).Decode(&addaderBook); err != nil {
			resp.ErrorBadRequest(w, errors.New("invalid request body"))
			return
		}

		// Проверка на существование книги
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM book WHERE book = $1 AND author = $2)", addaderBook.Book, addaderBook.Author).Scan(&exists)
		if err != nil {
			resp.ErrorInternal(w, err)
			return
		}
		if exists {
			resp.ErrorBadRequest(w, errors.New("book already exists"))
			return
		}

		var newBook repository.Book
		newBook.Book = addaderBook.Book
		newBook.Author = addaderBook.Author

		bloc := false
		newBook.Block = &bloc

		// Вставка новой книги в базу данных
		_, err = db.Exec("INSERT INTO book (book, author, block) VALUES ($1, $2, $3)", newBook.Book, newBook.Author, newBook.Block)
		if err != nil {
			resp.ErrorInternal(w, err)
			return
		}
		bookPtr := &Books

		// Получаем последний элемент
		lastElement := (**bookPtr)[len(**bookPtr)-1].Index
		newBook.Index = lastElement + 1

		library.AddBook(newBook)
		*Books = append(*Books, newBook)
		resp.OutputJSON(w, newBook) // Возвращаем добавленную книгу
	}
}

type AuthorRequest struct {
	Name string `json:"name"`
}

// @Summary Add a new author to the library
// @Description This endpoint allows you to add a new author to the library.
// @Tags Authors
// @Accept json
// @Produce json
// @Param author body AuthorRequest true "Author name"
// @Success 201 {object} string "Author added successfully"
// @Failure 400 {object} mErrorResponse "Invalid request"
// @Failure 500 {object} mErrorResponse "Internal server error"
// @Router /api/authors [post]
func AddAuthorHandler(resp Responder, library *Library) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var authorRequest AuthorRequest
		if err := json.NewDecoder(r.Body).Decode(&authorRequest); err != nil {
			resp.ErrorBadRequest(w, errors.New("invalid request body"))
			return
		}

		if authorRequest.Name == "" {
			resp.ErrorBadRequest(w, errors.New("author name is required"))
			return
		}

		library.mu.Lock()         // Блокируем запись
		defer library.mu.Unlock() // Разблокируем запись после завершения

		// Добавление автора в библиотеку (можно добавить логику для проверки уникальности)
		// Здесь предполагается, что у вас есть структура для хранения авторов
		library.Authors = append(library.Authors, authorRequest.Name)

		resp.OutputJSON(w, map[string]string{"message": "Author added successfully"})
	}
}

// getAuthorsHandler godoc
// @Summary Get all authors
// @Description Get a list of all authors in the library
// @Tags Authors
// @Produce json
// @Success 200 {array} string "List of authors"
// @Failure 404 {object} mErrorResponse "No authors found"
// @Router /api/get-authors [get]
func GetAuthorsHandler(resp Responder, library *Library) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		library.mu.RLock()         // Блокируем чтение
		defer library.mu.RUnlock() // Разблокируем чтение после завершения

		// Проверяем, есть ли авторы
		if len(library.Authors) == 0 {
			http.Error(w, "No authors found", http.StatusNotFound)
			return
		}

		// Возвращаем список авторов в формате JSON
		resp.OutputJSON(w, library.Authors)
	}
}
