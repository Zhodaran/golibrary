package router

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"studentgit.kata.academy/Zhodaran/go-kata/config"
	"studentgit.kata.academy/Zhodaran/go-kata/proxy/controller"
	"studentgit.kata.academy/Zhodaran/go-kata/proxy/internal/db/adapter"
	"studentgit.kata.academy/Zhodaran/go-kata/proxy/internal/modules/auth/service"
	"studentgit.kata.academy/Zhodaran/go-kata/proxy/postgres"
)

func Router(resp controller.Responder, db *sql.DB) http.Handler {

	postgres.RunMigrations(db)
	books := postgres.CreateTableBook(db)
	library := controller.NewLibrary()

	library.AddBooks(books)
	userRepo := adapter.NewPostgresUserRepository(db)
	bookController := &config.BookController{DB: db}

	userController := adapter.UserController{userRepo}

	r := chi.NewRouter()
	service.GenerateUsers(50)

	r.Use(middleware.Logger)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Post("/api/register", service.Register)
	r.Post("/api/login", service.Login)
	r.Post("/api/users", userController.CreateUser)        // Создание пользователя
	r.Get("/api/users/{id}", userController.GetUser)       // Получение пользователя по ID
	r.Put("/api/users/{id}", userController.UpdateUser)    // Обновление пользователя
	r.Delete("/api/users/{id}", userController.DeleteUser) // Удаление пользователя
	r.Get("/api/users", userController.ListUsers)

	r.Post("/api/book/take/{index}", controller.TakeBookHandler(resp, db, &books, library))
	r.Delete("/api/book/return/{index}", controller.ReturnBook(resp, db, &books, library))
	r.Get("/api/users", service.ListUsersHandler(resp))

	r.Post("/api/authors", controller.AddAuthorHandler(resp, library))

	r.Post("/api/book", controller.AddBookHandler(resp, db, library, &books))
	r.Get("/api/books", bookController.ListBook)
	r.Put("/api/book/{index}", controller.UpdateBook(resp, db))
	r.Get("/api/author", controller.ListAuthorsHandler(resp, library))
	r.Get("/api/get-authors", controller.GetAuthorsHandler(resp, library))

	// Используем обработчики с middleware

	return r
}
