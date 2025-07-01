package model

import "time"

type Credential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User represents a user in the database.
type User struct {
	ID           int        `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	PasswordHash *string    `json:"-" db:"password_hash"` // Pointer to handle nullable
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at" db:"updated_at"` // Pointer to handle nullable
}
