package service

import (
	"context"

	"alvintanoto.id/blog-htmx-templ/internal/dto"
	"alvintanoto.id/blog-htmx-templ/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct {
	repository *repository.Repository
}

func NewAuthenticationService(repository *repository.Repository) *AuthenticationService {
	return &AuthenticationService{
		repository: repository,
	}
}

func (s *AuthenticationService) SignIn(ctx context.Context, username, password string) (user *dto.UserDTO, err error) {
	entity, err := s.repository.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(entity.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return &dto.UserDTO{
		ID:       entity.ID,
		Username: entity.Username,
		Email:    entity.Email,
	}, nil
}
