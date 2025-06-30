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

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	repo := repository.NewUserRepository(conn, logger)
	service := service.NewUserService(repo, logger)
	handler := handler.NewUserHandler(service, logger)

	http.HandleFunc("/api/v1/register", handler.Register)
	http.HandleFunc("/api/v1/login", handler.Login)

	http.ListenAndServe(":3000", nil)
}
