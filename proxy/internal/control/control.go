package control

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/brianvoe/gofakeit"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/repository"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/service"
)

type Controller struct {
	geoService service.GeoServicer
}

func NewController(geoService service.GeoServicer) *Controller {
	return &Controller{geoService: geoService}
}

type ErrorResponse struct {
	Message string `json:"message"` // Сообщение об ошибке
	Code    int    `json:"code"`    // Код ошибки
}

func (c *Controller) GetGeoCoordinatesAddress(w http.ResponseWriter, r *http.Request) {
	var req service.RequestAddressSearch
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	geo, err := c.geoService.GetGeoCoordinatesAddress(req.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(geo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (c *Controller) GetGeoCoordinatesGeocode(w http.ResponseWriter, r *http.Request) {
	var req service.GeocodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	geo, err := c.geoService.GetGeoCoordinatesGeocode(req.Lat, req.Lng)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(geo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func RunMigrations(db *sql.DB) {
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

func CreateTableBook(db *sql.DB) []repository.Book {
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
	for i := 1; i < 101; i++ {
		block := false
		book := repository.Book{
			Index:     i,
			Book:      gofakeit.Sentence(1),                        // Генерация названия книги
			Author:    authors[gofakeit.Number(0, len(authors)-1)], // Случайный автор
			Block:     &block,                                      // Устанавливаем значение блокировки
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
