package repository

import (
	"context"

	"github.com/dosedaf/syncup-users-service/internal/model"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	conn *pgx.Conn
}

func NewUserRepository(conn *pgx.Conn) *Repository {
	return &Repository{
		conn: conn,
	}
}

func (r *Repository) InsertUser(credential model.Credential) error {
	query := "INSERT INTO users (email, password_hash) VALUES (@email, @password_hash)"
	args := pgx.NamedArgs{
		"email":         credential.Email,
		"password_hash": string(credential.Password),
	}

	_, err := r.conn.Exec(context.Background(), query, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetHashedPassword(email string) (string, error) {
	query := "SELECT password_hash FROM users WHERE email=@email"
	args := pgx.NamedArgs{
		"email": email,
	}

	var passwordDb string

	row := r.conn.QueryRow(context.Background(), query, args)

	err := row.Scan(&passwordDb)
	if err != nil {
		return "", err
	}

	return passwordDb, nil
}
