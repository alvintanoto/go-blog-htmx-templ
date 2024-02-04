package service

import "alvintanoto.id/blog-htmx-templ/internal/repository"

type PostService struct {
	repository *repository.Repository
}

func NewPostService(repository *repository.Repository) *PostService {
	return &PostService{
		repository: repository,
	}
}

func (s *PostService) CreateNewPost(userID, content string) (err error) {
	err = s.repository.PostRepository.CreateNewPost(userID, content)
	if err != nil {
		return err
	}

	return nil
}
