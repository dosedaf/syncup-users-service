package service

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/dosedaf/syncup-users-service/helper"
	"github.com/dosedaf/syncup-users-service/internal/model"
)

type mockRepo struct {
	MockIsEmailAvailable  func(ctx context.Context, email string) error
	MockInsertUser        func(ctx context.Context, credential model.Credential) error
	MockGetHashedPassword func(ctx context.Context, email string) (string, error)
}

func (m *mockRepo) IsEmailAvailable(ctx context.Context, email string) error {
	return m.MockIsEmailAvailable(ctx, email)
}

func (m *mockRepo) InsertUser(ctx context.Context, credential model.Credential) error {
	return m.MockInsertUser(ctx, credential)
}

func (m *mockRepo) GetHashedPassword(ctx context.Context, email string) (string, error) {
	return m.MockGetHashedPassword(ctx, email)
}

func TestRegisterNoError(t *testing.T) {
	credential := model.Credential{
		Email:    "newemail@gmail.com",
		Password: "thisisapassword",
	}

	mock := &mockRepo{MockInsertUser: func(context.Context, model.Credential) error { return nil }}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	service := NewUserService(mock, logger)
	err := service.Register(context.Background(), credential)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestRegisterErrorEmailAlreadyExists(t *testing.T) {
	credential := model.Credential{
		Email:    "newemail@gmail.com",
		Password: "thisisapassword",
	}

	mock := &mockRepo{MockInsertUser: func(context.Context, model.Credential) error { return helper.ErrEmailAlreadyExists }}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	service := NewUserService(mock, logger)
	err := service.Register(context.Background(), credential)
	if err == helper.ErrEmailAlreadyExists {
		t.Errorf(err.Error())
	}
}
