package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/dosedaf/syncup-users-service/database"
	"github.com/dosedaf/syncup-users-service/internal/handler"
	"github.com/dosedaf/syncup-users-service/internal/repository"
	"github.com/dosedaf/syncup-users-service/internal/service"
	"github.com/dosedaf/syncup-users-service/middleware"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file", "error", err)
	}

	runMigrate(os.Getenv("DATABASE_URL"), logger)

	conn, err := database.ConnectDB()
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	jwtSecret := os.Getenv("SECRET")
	repo := repository.NewUserRepository(conn, logger)
	svc := service.NewUserService(repo, logger, jwtSecret)
	h := handler.NewUserHandler(svc, logger)
	authMiddleware := middleware.NewMiddleware(repo, logger, jwtSecret)

	mux := http.NewServeMux()
	mux.Handle("POST /api/v1/register", http.HandlerFunc(h.Register))
	mux.Handle("POST /api/v1/login", http.HandlerFunc(h.Login))
	mux.Handle("GET /api/v1/me", authMiddleware.JWTMiddleware(http.HandlerFunc(h.Me)))

	logger.Info("Starting server on port 3000")
	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		logger.Error("Failed to start new server", "error", err)
		os.Exit(1)
	}
}

func runMigrate(databaseURL string, logger *slog.Logger) {
	migration, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		logger.Error("Unable to create new migrate instance", "error", err)
		os.Exit(1)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("Failed to run migrate up", "error", err)
		os.Exit(1)
	}

	logger.Info("DB migrated successfully")
}
