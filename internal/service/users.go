package service

import (
	"log"
	"os"
	"time"

	"github.com/dosedaf/syncup-users-service/internal/model"
	"github.com/dosedaf/syncup-users-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *repository.Repository
}

func NewUserService(repo *repository.Repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) Register(credential model.Credential) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credential.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	credential.Password = string(hashedPassword)

	err = s.repository.InsertUser(credential)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Login(credential model.Credential) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	passwordDb, err := s.repository.GetHashedPassword(credential.Email)
	if err != nil {
		return "", nil
	}

	// mungkin salah pass
	err = bcrypt.CompareHashAndPassword([]byte(passwordDb), []byte(credential.Password))
	if err != nil {
		log.Print(err.Error())
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
		log.Print(err.Error())
		return "", nil
	}

	return tokenString, nil
}
