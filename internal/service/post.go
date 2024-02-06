package service

import (
	"alvintanoto.id/blog-htmx-templ/internal/dto"
	"alvintanoto.id/blog-htmx-templ/internal/repository"
)

type PostService struct {
	repository *repository.Repository
}

func NewPostService(repository *repository.Repository) *PostService {
	return &PostService{
		repository: repository,
	}
}

func (s *PostService) CreateNewPost(userID, content string, isDraft bool) (err error) {
	err = s.repository.PostRepository.CreateNewPost(userID, content, isDraft)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostService) GetUserPostByUserID(user *dto.UserDTO, page int) (posts []dto.PostDTO, err error) {
	entities, err := s.repository.PostRepository.GetUserPostByUserID(user.ID, page)
	if err != nil {
		return nil, err
	}

	for _, entity := range entities {
		var post = new(dto.PostDTO)
		post.ID = entity.ID
		post.Content = entity.Content
		post.ReplyCounts = entity.ReplyCount
		post.Likes = entity.LikeCount
		post.Dislikes = entity.DislikeCount
		post.Impressions = entity.ImpressionCount
		post.SavedCounts = entity.SaveCount
		post.PostedAt = entity.PostedAt.Format("02 Jan 2006 15:04:05")
		post.Poster = *user

		posts = append(posts, *post)
	}

	return posts, err
}
