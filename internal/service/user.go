package service

import (
	"log"

	"alvintanoto.id/blog-htmx-templ/internal/dto"
	"alvintanoto.id/blog-htmx-templ/internal/repository"
)

type UserService struct {
	repository *repository.Repository
}

func NewUserService(repository *repository.Repository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) RegisterUser(username, email, password string) (user *dto.UserDTO, err error) {
	entity, err := s.repository.UserRepository.RegisterUser(username, email, password)
	if err != nil {
		log.Println("failed to register user: ", err.Error())
		return nil, err
	}

	return &dto.UserDTO{
		ID:        entity.ID,
		Email:     entity.Email,
		Username:  entity.Username,
		CreatedAt: entity.CreatedAt.Format("02 January 2006"),
	}, nil
}
