package main

import (
	"log"
	"net/http"

	"github.com/dosedaf/syncup-users-service/database"
	"github.com/dosedaf/syncup-users-service/internal/handler"
	"github.com/dosedaf/syncup-users-service/internal/repository"
	"github.com/dosedaf/syncup-users-service/internal/service"
)

func main() {
	conn, err := database.ConnectDB()
	if err != nil {
		log.Print(err.Error())
	}

	repo := repository.NewUserRepository(conn)
	service := service.NewUserService(repo)
	handler := handler.NewUserHandler(service)

	http.HandleFunc("/api/v1/register", handler.Register)
	http.HandleFunc("/api/v1/login", handler.Login)

	http.ListenAndServe(":3000", nil)
}
