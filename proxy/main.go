package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	_ "studentgit.kata.academy/Zhodaran/go-kata/docs"

	"studentgit.kata.academy/Zhodaran/go-kata/controller"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/control"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/repository"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/router"
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

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("DB_USER:", os.Getenv("DB_USER"))
	fmt.Println("DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
	fmt.Println("DB_NAME:", os.Getenv("DB_NAME"))
	fmt.Println("DB_HOST:", os.Getenv("DB_HOST"))
	fmt.Println("DB_PORT:", os.Getenv("DB_PORT"))
	fmt.Println("Запуск задержки")
	time.Sleep(10 * time.Second)

	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to the database:", err)
	}

	control.RunMigrations(db)
	books := control.CreateTableBook(db)
	library := controller.NewLibrary()

	library.AddBooks(books)

	userRepo := repository.NewPostgresUserRepository(db)
	bookController := &control.BookController{DB: db}
	userController := control.NewUserController(userRepo)

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	geoService := service.NewGeoService("d9e0649452a137b73d941aa4fb4fcac859372c8c", "ec99b849ebf21277ec821c63e1a2bc8221900b1d") // Создаем новый экземпляр GeoService
	resp := controller.NewResponder(logger)
	r := router.Router(bookController, userController, resp, geoService, db, &books, library)

	srv := &Server{
		Server: http.Server{
			Addr:         ":8080",
			Handler:      r,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	go srv.Serve()

	waitForShutdown(srv)
}

func initDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	return sql.Open("postgres", connStr)
}

func waitForShutdown(srv *Server) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server: %v\n", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
