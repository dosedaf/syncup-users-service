package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/dosedaf/syncup-users-service/internal/model"
	"github.com/dosedaf/syncup-users-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type ServiceInstance interface {
	Register(ctx context.Context, credential model.Credential) error
	Login(ctx context.Context, credential model.Credential) (string, error)
}

type Service struct {
	repository repository.RepositoryInstance
	logger     *slog.Logger
}

func NewUserService(repo repository.RepositoryInstance, logger *slog.Logger) ServiceInstance {
	return &Service{
		repository: repo,
		logger:     logger,
	}
}

func (s *Service) Register(ctx context.Context, credential model.Credential) error {
	err := s.repository.IsEmailAvailable(ctx, credential.Email)
	if err != nil {
		s.logger.Info("failed to register", "error", err)
		return fmt.Errorf("failed while checking email availability: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credential.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Info("failed to register", "error", err)
		return err
	}

	credential.Password = string(hashedPassword)

	err = s.repository.InsertUser(ctx, credential)
	if err != nil {
		s.logger.Info("failed to register", "error", err)
		return fmt.Errorf("failed while inserting new user %s: %w", credential.Email, err)
	}

	return nil
}

func (s *Service) Login(ctx context.Context, credential model.Credential) (string, error) {
	passwordDb, err := s.repository.GetHashedPassword(ctx, credential.Email)
	if err != nil {
		s.logger.Info("failed to login", "error", err)
		return "", fmt.Errorf("failed while getting hashed password from user %s: %w", credential.Email, err)
	}

	// mungkin salah pass
	err = bcrypt.CompareHashAndPassword([]byte(passwordDb), []byte(credential.Password))
	if err != nil {
		s.logger.Info("failed to login", "error", err)
		return "", err
	}

	// aman
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": credential.Email,
		"iss": "app",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := claims.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		s.logger.Info("failed to login", "error", err)
		return "", err
	}

	return tokenString, nil
}
