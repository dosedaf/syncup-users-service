package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dosedaf/syncup-users-service/helper"
	"github.com/dosedaf/syncup-users-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey = contextKey("user")

type Middleware struct {
	repo      repository.RepositoryInstance
	logger    *slog.Logger
	jwtSecret []byte
}

func NewMiddleware(repo repository.RepositoryInstance, logger *slog.Logger, jwtSecret string) *Middleware {
	return &Middleware{
		repo:      repo,
		logger:    logger,
		jwtSecret: []byte(jwtSecret),
	}
}

func (m *Middleware) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helper.JSONError(w, http.StatusUnauthorized, "authorization header required")
			return
		}

		tokenStr, found := strings.CutPrefix(authHeader, "Bearer ")
		if !found {
			helper.JSONError(w, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return m.jwtSecret, nil
		})

		if err != nil {
			m.logger.Info("invalid JWT", "error", err)
			helper.JSONError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			email, _ := claims.GetSubject()
			user, err := m.repo.GetUserByEmail(r.Context(), email)
			if err != nil {
				if errors.Is(err, helper.ErrUserNotFound) {
					helper.JSONError(w, http.StatusUnauthorized, "user not found")
					return
				}

				helper.JSONError(w, http.StatusInternalServerError, "internal server error")
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))

		} else {
			helper.JSONError(w, http.StatusUnauthorized, "invalid token claims")
		}
	})
}
