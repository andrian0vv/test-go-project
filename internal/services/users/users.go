package users

import (
	"context"
	"fmt"

	"github.com/andrian0vv/test-go-project/internal/models"
)

type repository interface {
	CreateUser(ctx context.Context, user models.User) error
}

type hasher interface {
	HashPassword(string) (string, error)
}

type Service struct {
	repository repository
	hasher     hasher
}

func New(repository repository, hasher hasher) *Service {
	return &Service{
		repository: repository,
		hasher:     hasher,
	}
}

func (s *Service) CreateUser(ctx context.Context, user User) error {
	model, err := s.convertUser(user)
	if err != nil {
		return fmt.Errorf("convert user: %w", err)
	}

	if err := s.repository.CreateUser(ctx, model); err != nil {
		return err
	}

	return nil
}

func (s *Service) convertUser(user User) (models.User, error) {
	passwordHash, err := s.hasher.HashPassword(user.Password)
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		FullName:  fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Age:       user.Age,
		IsMarried: user.IsMarried,
		Password:  passwordHash,
	}, nil
}
