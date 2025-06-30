package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/dosedaf/syncup-users-service/helper"
	"github.com/dosedaf/syncup-users-service/internal/model"
	"github.com/jackc/pgx/v5"
)

type RepositoryInstance interface {
	IsEmailAvailable(ctx context.Context, email string) error
	InsertUser(ctx context.Context, credential model.Credential) error
	GetHashedPassword(ctx context.Context, email string) (string, error)
}

type Repository struct {
	conn   *pgx.Conn
	logger *slog.Logger
}

func NewUserRepository(conn *pgx.Conn, logger *slog.Logger) RepositoryInstance {
	return &Repository{
		conn:   conn,
		logger: logger,
	}
}

func (r *Repository) IsEmailAvailable(ctx context.Context, email string) error {
	query := "SELECT email FROM users WHERE email=@email"
	args := pgx.NamedArgs{
		"email": email,
	}

	var emailDb string

	row := r.conn.QueryRow(ctx, query, args)

	err := row.Scan(&emailDb)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}

		return err
	}

	return helper.ErrEmailAlreadyExists
}

func (r *Repository) InsertUser(ctx context.Context, credential model.Credential) error {
	query := "INSERT INTO users (email, password_hash) VALUES (@email, @password_hash)"
	args := pgx.NamedArgs{
		"email":         credential.Email,
		"password_hash": string(credential.Password),
	}

	_, err := r.conn.Exec(ctx, query, args)
	if err != nil {
		r.logger.Info("failed executing query", "error", err)
		return err
	}

	return nil
}

func (r *Repository) GetHashedPassword(ctx context.Context, email string) (string, error) {
	query := "SELECT password_hash FROM users WHERE email=@email"
	args := pgx.NamedArgs{
		"email": email,
	}

	var passwordDb string

	row := r.conn.QueryRow(ctx, query, args)

	err := row.Scan(&passwordDb)
	if err != nil {
		r.logger.Info("failed scanning the row", "error", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return "", helper.ErrUserNotFound
		}

		return "", err
	}

	return passwordDb, nil
}
