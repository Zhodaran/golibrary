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
	"strings"
	"syscall"
	"time"

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

	userRepo := repository.NewPostgresUserRepository(db)
	userController := control.NewUserController(userRepo)

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	geoService := service.NewGeoService("d9e0649452a137b73d941aa4fb4fcac859372c8c", "ec99b849ebf21277ec821c63e1a2bc8221900b1d") // Создаем новый экземпляр GeoService
	resp := controller.NewResponder(logger)
	r := router(userController, resp, geoService, db)

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

func createTableBook(db *sql.DB) {
	table := ` 
	CREATE TABLE IF NOT EXISTS book (
		index SERIAL PRIMARY KEY,
		book VARCHAR(50) NOT NULL,
		author VARCHAR(255) NOT NULL,
		block BOOLEAN,
		take_count INT DEFAULT 0
	);`
	_, err := db.Exec(table)
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
	books := []repository.Book{
		{Book: "Книга 1", Author: "Автор 1", Block: false},
		{Book: "Книга 2", Author: "Автор 2", Block: false},
		{Book: "Книга 3", Author: "Автор 3", Block: false},
		{Book: "Книга 4", Author: "Автор 4", Block: false},
		{Book: "Книга 5", Author: "Автор 5", Block: false},
		{Book: "Книга 6", Author: "Автор 6", Block: false},
		{Book: "Книга 7", Author: "Автор 7", Block: false},
		{Book: "Книга 8", Author: "Автор 8", Block: false},
		{Book: "Книга 9", Author: "Автор 9", Block: false},
		{Book: "Книга 10", Author: "Автор 10", Block: false},
	}

	for _, b := range books {
		_, err := db.Exec("INSERT INTO book (book, author, block) VALUES ($1, $2, $3)", b.Book, b.Author, b.Block)
		if err != nil {
			log.Fatalf("Error inserting book: %v", err)
		}
	}
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

// @Summary Get Geo Coordinates by Address
// @Description This endpoint allows you to get geo coordinates by address.
// @Tags TakeBook
// @Accept json
// @Produce json
// @Param index path int true "Book INDEX"
// @Success 200 {object} service.ResponseAddress "Успешное выполнение"
// @Failure 400 {object} mErrorResponse "Ошибка запроса"
// @Failure 500 {object} mErrorResponse "Ошибка подключения к серверу"
// @Security BearerAuth
// @Router /api/book/take [post]
func takeBookHandler(resp controller.Responder, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Index int `json:"index"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp.ErrorBadRequest(w, err)
			return
		}

		// Обновление записи в таблице book
		result, err := db.Exec("UPDATE book SET block = $1, take_count = take_count + 1 WHERE index = $2 AND block = $3", true, req.Index, false)
		if err != nil {
			resp.ErrorInternal(w, err)
			return
		}

		// Проверка, была ли книга успешно обновлена
		if rowsAffected, err := result.RowsAffected(); err != nil || rowsAffected == 0 {
			resp.ErrorBadRequest(w, errors.New("book not found or already taken"))
			return
		}

		resp.OutputJSON(w, map[string]string{"message": "Book taken successfully"})
	}
}

func router(userController *control.UserController, resp controller.Responder, geoService service.GeoProvider, db *sql.DB) http.Handler {
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

	r.Post("/api/book/take", takeBookHandler(resp, db))

	// Используем обработчики с middleware
	r.With(TokenAuthMiddleware(resp)).Post("/api/address/geocode", geocodeHandler(resp, geoService))
	r.With(TokenAuthMiddleware(resp)).Post("/api/address/search", searchHandler(resp, geoService))

	return r
}
