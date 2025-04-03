package router

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"studentgit.kata.academy/Zhodaran/go-kata/controller"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/auth"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/control"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/middle"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/repository"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/service"
)

func Router(bookController *control.BookController, userController *control.UserController, resp controller.Responder, geoService service.GeoProvider, db *sql.DB, books *[]repository.Book, library *controller.Library) http.Handler {
	r := chi.NewRouter()
	auth.GenerateUsers(50)
	r.Use(middleware.Logger)
	r.Use(middle.ProxyMiddleware)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Post("/api/register", auth.Register)
	r.Post("/api/login", auth.Login)
	r.Post("/api/users", userController.CreateUser)        // Создание пользователя
	r.Get("/api/users/{id}", userController.GetUser)       // Получение пользователя по ID
	r.Put("/api/users/{id}", userController.UpdateUser)    // Обновление пользователя
	r.Delete("/api/users/{id}", userController.DeleteUser) // Удаление пользователя
	r.Get("/api/users", userController.ListUsers)

	r.Post("/api/book/take/{index}", controller.TakeBookHandler(resp, db, books, library))
	r.Delete("/api/book/return/{index}", controller.ReturnBook(resp, db, books, library))
	r.Get("/api/users", auth.ListUsersHandler(resp))

	r.Post("/api/authors", controller.AddAuthorHandler(resp, library))

	r.Post("/api/book", controller.AddBookHandler(resp, db, library, books))
	r.Get("/api/books", bookController.ListBook)
	r.Put("/api/book/{index}", controller.UpdateBook(resp, db))
	r.Get("/api/author", controller.ListAuthorsHandler(resp, library))
	r.Get("/api/get-authors", controller.GetAuthorsHandler(resp, library))

	// Используем обработчики с middleware
	r.With(middle.TokenAuthMiddleware(resp)).Post("/api/address/geocode", controller.GeocodeHandler(resp, geoService))
	r.With(middle.TokenAuthMiddleware(resp)).Post("/api/address/search", controller.SearchHandler(resp, geoService))

	return r
}
