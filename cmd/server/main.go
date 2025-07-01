package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/dosedaf/syncup-users-service/database"
	"github.com/dosedaf/syncup-users-service/internal/handler"
	"github.com/dosedaf/syncup-users-service/internal/repository"
	"github.com/dosedaf/syncup-users-service/internal/service"
	"github.com/dosedaf/syncup-users-service/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print(err.Error())
	}

	conn, err := database.ConnectDB()
	if err != nil {
		log.Print(err.Error())
	}

	jwtSecret := os.Getenv("SECRET")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	repo := repository.NewUserRepository(conn, logger)
	service := service.NewUserService(repo, logger, jwtSecret)
	handler := handler.NewUserHandler(service, logger)

	middleware := middleware.NewMiddleware(repo, logger, jwtSecret)

	mux := http.NewServeMux()

	mux.Handle("POST /api/v1/register", http.HandlerFunc(handler.Register))
	mux.Handle("POST /api/v1/login", http.HandlerFunc(handler.Login))
	mux.Handle("GET /api/v1/me", middleware.JWTMiddleware(http.HandlerFunc(handler.Me)))

	// http.HandleFunc("/api/v1/register", handler.Register)
	//http.HandleFunc("/api/v1/login", handler.Login)

	err = http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
