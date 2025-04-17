package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	_ "studentgit.kata.academy/Zhodaran/go-kata/proxy/docs"

	"studentgit.kata.academy/Zhodaran/go-kata/config"
	"studentgit.kata.academy/Zhodaran/go-kata/proxy/controller"

	"studentgit.kata.academy/Zhodaran/go-kata/proxy/router"
)

// @title Address API
// @version 1.0
// @description API для поиска
// @host localhost:8080
// @BasePath
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @RequestAddressSearch представляет запрNewUserControllerос для поиска
// @Description Этот эндпоинт позволяет получить адрес по наименованию
// @Param address body ResponseAddress true "Географические координаты"

// TokenResponse представляет ответ с токеном

// LoginResponse представляет ответ при успешном входе

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Error connecting to the database:", err) // Обработка ошибки
	}
	defer db.Close()
	config.PullSQL()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	resp := controller.NewResponder(logger)
	r := router.Router(resp, db)

	srv := &config.Server{
		Server: http.Server{
			Addr:         ":8080",
			Handler:      r,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	go srv.Serve()

	config.WaitForShutdown(srv)
}
