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

func (s *UserService) GetUserConfig(userID string) (data dto.UserConfigDTO, err error) {
	entities, err := s.repository.UserRepository.GetUserConfig(userID)
	if err != nil {
		log.Println("failed to register user: ", err.Error())
		return nil, err
	}

	data = make(dto.UserConfigDTO)
	for _, entity := range entities {
		data[entity.Key] = entity.Value
	}

	return data, nil
}

func (s *UserService) InsertUserConfig(userID, key, value string) (err error) {
	return s.repository.UserRepository.InsertUserConfig(userID, key, value)
}

func (s *UserService) UpdateUserConfig(userID, key, value string) (err error) {
	return s.repository.UserRepository.UpdateUserConfig(userID, key, value)
}
