package service

import "alvintanoto.id/blog-htmx-templ/internal/repository"

type Service struct {
	UserService *UserService
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		UserService: NewUserService(repository),
	}
}
