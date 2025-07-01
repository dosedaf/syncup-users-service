package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/dosedaf/syncup-users-service/helper"
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
	jwtSecret  []byte
}

func NewUserService(repo repository.RepositoryInstance, logger *slog.Logger, jwtSecret string) ServiceInstance {
	return &Service{
		repository: repo,
		logger:     logger,
		jwtSecret:  []byte(jwtSecret),
	}
}

func (s *Service) Register(ctx context.Context, credential model.Credential) error {
	err := s.repository.IsEmailAvailable(ctx, credential.Email)
	if err != nil {
		if errors.Is(err, helper.ErrEmailAlreadyExists) {
			s.logger.Info(
				"User registration blocked: email already exists",
				"email", credential.Email,
			)

			return err
		}

		s.logger.Error(
			"Failed to check email availability",
			"email", credential.Email,
			"error", err,
		)

		return fmt.Errorf("failed while checking email availability: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credential.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(
			"Failed while generating hashed password",
			"email", credential.Email,
			"error", err,
		)
		return err
	}

	credential.Password = string(hashedPassword)

	err = s.repository.InsertUser(ctx, credential)
	if err != nil {
		s.logger.Error(
			"Failed while inserting new user",
			"email", credential.Email,
			"error", err,
		)

		return fmt.Errorf("failed while inserting new user %s: %w", credential.Email, err)
	}

	return nil
}

func (s *Service) Login(ctx context.Context, credential model.Credential) (string, error) {
	passwordDb, err := s.repository.GetHashedPassword(ctx, credential.Email)
	if err != nil {
		if errors.Is(err, helper.ErrUserNotFound) {
			s.logger.Info(
				"User login blocked: user with this email does not exist",
				"email", credential.Email,
			)

			return "", err
		}

		s.logger.Error(
			"Failed while getting hashed password",
			"email", credential.Email,
			"error", err,
		)

		return "", fmt.Errorf("failed while getting hashed password from user %s: %w", credential.Email, err)
	}

	// could be wrong pass (mismatched) or an actual error. how do i differ them?
	// https://cs.opensource.google/go/x/crypto/+/refs/tags/v0.37.0:bcrypt/bcrypt.go;l=95
	// i learnt that u can open the source code and look for yourself what are the errors returned from a specific method
	err = bcrypt.CompareHashAndPassword([]byte(passwordDb), []byte(credential.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			s.logger.Info(
				"User login blocked: wrong password",
				"email", credential.Email,
			)

			return "", helper.ErrWrongPassword
		}

		s.logger.Error(
			"Failed while comparing hash and password",
			"email", credential.Email,
			"error", err,
		)

		return "", fmt.Errorf("failed while comparing hash and password from user %s: %w", credential.Email, err)
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": credential.Email,
		"iss": "app",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := claims.SignedString(s.jwtSecret)
	if err != nil {
		s.logger.Error(
			"Failed while getting signed string",
			"email", credential.Email,
			"error", err,
		)

		return "", err
	}

	return tokenString, nil
}
