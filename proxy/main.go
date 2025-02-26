package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	_ "studentgit.kata.academy/Zhodaran/go-kata/docs"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"studentgit.kata.academy/Zhodaran/go-kata/controller"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/auth"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/control"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/repository"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/service"
)

// @title Address API
// @version 1.0
// @description API для поиска
// @host localhost:8080
// @BasePath
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @RequestAddressSearch представляет запрос для поиска
// @Description Этот эндпоинт позволяет получить адрес по наименованию
// @Param address body ResponseAddress true "Географические координаты"

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

// TokenResponse представляет ответ с токеном

// LoginResponse представляет ответ при успешном входе

type Server struct {
	http.Server
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

func (s *Server) Serve() {
	log.Println("Starting server...")
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: &v", err)
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	fmt.Println("DB_USER:", os.Getenv("DB_USER"))
	fmt.Println("DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
	fmt.Println("DB_NAME:", os.Getenv("DB_NAME"))
	fmt.Println("DB_HOST:", os.Getenv("DB_HOST"))
	fmt.Println("DB_PORT:", os.Getenv("DB_PORT"))
	fmt.Println("Запуск задержки")
	time.Sleep(10 * time.Second)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to the database:", err)
	}

	runMigrations(db)
	books := createTableBook(db)
	library := NewLibrary()

	userRepo := repository.NewPostgresUserRepository(db)
	BookRepo := repository.NewPostgresBookRepository(db)
	userController := control.NewUserController(userRepo)
	bookController := control.NewBookController(BookRepo)

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	geoService := service.NewGeoService("d9e0649452a137b73d941aa4fb4fcac859372c8c", "ec99b849ebf21277ec821c63e1a2bc8221900b1d") // Создаем новый экземпляр GeoService
	resp := controller.NewResponder(logger)
	r := router(userController, resp, geoService, db, bookController, &books, library)

	srv := &Server{
		Server: http.Server{
			Addr:         ":8080",
			Handler:      r,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	go srv.Serve()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Ошибка при завершении работы: %v\n", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}

func runMigrations(db *sql.DB) {
	migrationSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		email VARCHAR(255) NOT NULL,
		deleted_at TIMESTAMP NULL
	);`

	_, err := db.Exec(migrationSQL)
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
}

func createTableBook(db *sql.DB) []repository.Book {
	table := ` 
	CREATE TABLE IF NOT EXISTS book (
		index SERIAL PRIMARY KEY,
		book VARCHAR(50) NOT NULL,
		author VARCHAR(255) NOT NULL,
		block BOOLEAN,
		take_count INT DEFAULT 0
	);`

	var authors []string
	for i := 0; i < 10; i++ {
		author := gofakeit.Name()
		authors = append(authors, author)
	}
	_, err := db.Exec(table)
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
	var books []repository.Book
	for i := 0; i < 100; i++ {
		book := repository.Book{
			Index:     i,
			Book:      gofakeit.Sentence(1),                        // Генерация названия книги
			Author:    authors[gofakeit.Number(0, len(authors)-1)], // Случайный автор
			Block:     false,                                       // Устанавливаем значение блокировки
			TakeCount: 0,                                           // Начальное значение take_count
		}
		books = append(books, book) // Добавляем книгу в массив
	}

	// Вставка книг в базу данных
	for _, b := range books {
		_, err := db.Exec("INSERT INTO book (book, author, block) VALUES ($1, $2, $3)", b.Book, b.Author, b.Block)
		if err != nil {
			log.Fatalf("Error inserting book: %v", err)
		}
	}

	return books

}

func proxyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			next.ServeHTTP(w, r)
			return
		}
		proxyURL, _ := url.Parse("http://hugo:1313")
		proxy := httputil.NewSingleHostReverseProxy(proxyURL)
		proxy.ServeHTTP(w, r)
	})
}

func TokenAuthMiddleware(resp controller.Responder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				resp.ErrorUnauthorized(w, errors.New("missing authorization token"))
				return
			}

			token = strings.TrimPrefix(token, "Bearer ")

			_, err := auth.TokenAuth.Decode(token)
			if err != nil {
				resp.ErrorUnauthorized(w, err)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func geocodeHandler(resp controller.Responder, geoService service.GeoProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req GeocodeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp.ErrorBadRequest(w, err)
			return
		}

		geo, err := geoService.GetGeoCoordinatesGeocode(req.Lat, req.Lng)
		if err != nil {
			resp.ErrorInternal(w, err)
			return
		}

		resp.OutputJSON(w, geo)
	}
}

func searchHandler(resp controller.Responder, geoService service.GeoProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestAddressSearch
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp.ErrorBadRequest(w, err)
			return
		}

		geo, err := geoService.GetGeoCoordinatesAddress(req.Query)
		if err != nil {
			resp.ErrorInternal(w, err)
			return
		}

		resp.OutputJSON(w, geo)
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
// @Success 200 {object} service.ResponseAddress "Успешное выполнение"
// @Failure 400 {object} mErrorResponse "Ошибка запроса"
// @Failure 500 {object} mErrorResponse "Ошибка подключения к серверу"
// @Security BearerAuth
// @Router /api/book/take/{index} [post]
func takeBookHandler(resp controller.Responder, db *sql.DB, Books *[]repository.Book, library *Library) http.HandlerFunc {
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
// @Success 200 {object} service.ResponseAddress "Успешное выполнение"
// @Failure 400 {object} mErrorResponse "Ошибка запроса"
// @Failure 500 {object} mErrorResponse "Ошибка подключения к серверу"
// @Security BearerAuth
// @Router /api/book/return/{index} [delete]
func ReturnBook(resp controller.Responder, db *sql.DB, Books *[]repository.Book, library *Library) http.HandlerFunc {
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
func updateBook(resp controller.Responder, db *sql.DB) http.HandlerFunc {
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

// @Summary Get List of Authors
// @Description This endpoint returns a list of all authors from the library.
// @Tags Authors
// @Accept json
// @Produce json
// @Success 200 {array} string "List of authors"
// @Failure 500 {object} mErrorResponse "Internal server error"
// @Router /api/authors [get]
func listAuthorsHandler(resp controller.Responder, library *Library) http.HandlerFunc {
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
// @Param book body repository.Book true "Book details"
// @Success 201 {object} repository.Book "Book added successfully"
// @Failure 400 {object} mErrorResponse "Invalid request"
// @Failure 500 {object} mErrorResponse "Internal server error"
// @Router /api/book [post]
func addBookHandler(resp controller.Responder, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newBook repository.Book
		if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
			resp.ErrorBadRequest(w, errors.New("invalid request body"))
			return
		}

		// Вставка новой книги в базу данных
		_, err := db.Exec("INSERT INTO book (book, author, block) VALUES ($1, $2, $3)", newBook.Book, newBook.Author, newBook.Block)
		if err != nil {
			resp.ErrorInternal(w, err)
			return
		}

		resp.OutputJSON(w, newBook) // Возвращаем добавленную книгу
	}
}

// @Summary Get List of Registered Users
// @Description This endpoint returns a list of all registered users.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {array} auth.User "List of registered users"
// @Failure 500 {object} mErrorResponse "Internal server error"
// @Router /api/users [get]
func listUsersHandler(resp controller.Responder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Здесь предполагается, что Users - это глобальная переменная, содержащая всех пользователей
		var users []auth.User
		for _, user := range auth.Users {
			users = append(users, user)
		}

		resp.OutputJSON(w, users) // Возвращаем список пользователей
	}
}

// @Summary Add a new author to the library
// @Description This endpoint allows you to add a new author to the library.
// @Tags Authors
// @Accept json
// @Produce json
// @Param author body string true "Author name"
// @Success 201 {object} string "Author added successfully"
// @Failure 400 {object} mErrorResponse "Invalid request"
// @Failure 500 {object} mErrorResponse "Internal server error"
// @Router /api/authors [post]
func addAuthorHandler(resp controller.Responder, library *Library) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var authorName string
		if err := json.NewDecoder(r.Body).Decode(&authorName); err != nil {
			resp.ErrorBadRequest(w, errors.New("invalid request body"))
			return
		}

		if authorName == "" {
			resp.ErrorBadRequest(w, errors.New("author name is required"))
			return
		}

		library.mu.Lock()         // Блокируем запись
		defer library.mu.Unlock() // Разблокируем запись после завершения

		// Добавление автора в библиотеку (можно добавить логику для проверки уникальности)
		// Здесь предполагается, что у вас есть структура для хранения авторов
		library.Authors = append(library.Authors, authorName)

		resp.OutputJSON(w, map[string]string{"message": "Author added successfully"})
	}
}

func router(userController *control.UserController, resp controller.Responder, geoService service.GeoProvider, db *sql.DB, bookController *control.BookController, books *[]repository.Book, library *Library) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(proxyMiddleware)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Post("/api/register", auth.Register)
	r.Post("/api/login", auth.Login)
	r.Post("/api/users", userController.CreateUser)        // Создание пользователя
	r.Get("/api/users/{id}", userController.GetUser)       // Получение пользователя по ID
	r.Put("/api/users/{id}", userController.UpdateUser)    // Обновление пользователя
	r.Delete("/api/users/{id}", userController.DeleteUser) // Удаление пользователя
	r.Get("/api/users", userController.ListUsers)

	r.Post("/api/book/take/{index}", takeBookHandler(resp, db, books, library))
	r.Delete("/api/book/return/{index}", ReturnBook(resp, db, books, library))
	r.Get("/api/users", listUsersHandler(resp))

	r.Post("/api/authors", addAuthorHandler(resp, library))

	r.Post("/api/book", addBookHandler(resp, db))
	r.Get("/api/book", bookController.ListBook)
	r.Put("/api/book/{index}", updateBook(resp, db))
	r.Get("/api/authors", listAuthorsHandler(resp, library))

	// Используем обработчики с middleware
	r.With(TokenAuthMiddleware(resp)).Post("/api/address/geocode", geocodeHandler(resp, geoService))
	r.With(TokenAuthMiddleware(resp)).Post("/api/address/search", searchHandler(resp, geoService))

	return r
}
