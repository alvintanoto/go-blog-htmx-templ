package service

import "alvintanoto.id/blog-htmx-templ/internal/repository"

type Service struct {
	AuthenticationService *AuthenticationService
	UserService           *UserService
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		AuthenticationService: NewAuthenticationService(repository),
		UserService:           NewUserService(repository),
	}
}
