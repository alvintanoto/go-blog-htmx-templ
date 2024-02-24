package service

import (
	"time"

	"alvintanoto.id/blog-htmx-templ/internal/dto"
	"alvintanoto.id/blog-htmx-templ/internal/entity"
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

func (s *PostService) CreatePost(userID string, postID string, payloadContent string, action string) (err error) {
	var isDrafting bool = action == "draft"

	// new entry
	if postID == "" {
		err = s.repository.PostRepository.CreateNewPost(userID, payloadContent, isDrafting)
		if err != nil {
			return err
		}

		return nil
	}

	// get existing post
	post, err := s.repository.PostRepository.GetUserPost(userID, postID)
	if err != nil {
		return err
	}

	currentTime := time.Now()
	post.Content = payloadContent
	post.UpdatedAt = currentTime
	post.UpdatedBy = &userID

	if !isDrafting {
		post.IsDraft = false
		post.PostedAt = &currentTime
	}

	err = s.repository.PostRepository.UpdatePost(post)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostService) GetUserPosts(user *dto.UserDTO, lastPosition int) (posts []dto.PostDTO, err error) {
	var entities []*entity.Post

	if lastPosition == 0 {
		entities, err = s.repository.PostRepository.GetUserPosts(user.ID)
		if err != nil {
			return nil, err
		}
	} else {
		entities, err = s.repository.PostRepository.GetMoreUserPosts(user.ID, lastPosition)
		if err != nil {
			return nil, err
		}
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

func (s *PostService) GetPostDetail(postID string) (post *dto.PostDTO, err error) {
	entity, err := s.repository.PostRepository.GetPost(postID)
	if err != nil {
		return nil, err
	}
	post = new(dto.PostDTO)
	post.ID = entity.ID
	post.Content = entity.Content

	return post, nil
}

func (s *PostService) GetUserPost(user *dto.UserDTO, postID string) (post *dto.PostDTO, err error) {
	entity, err := s.repository.PostRepository.GetUserPost(user.ID, postID)
	if err != nil {
		return nil, err
	}

	post = new(dto.PostDTO)
	post.ID = entity.ID
	post.Content = entity.Content

	return post, nil
}

func (s *PostService) GetUserDraft(user *dto.UserDTO, page int) (posts []dto.PostDTO, err error) {
	// TODO: get following user post
	entities, err := s.repository.PostRepository.GetDrafts(user.ID, page)
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
		post.CreatedAt = entity.CreatedAt.Format("02 Jan 2006 15:04:05")
		post.Poster = *user

		posts = append(posts, *post)
	}

	return posts, err
}

func (s *PostService) GetPublicTimeline(lastPosition int) (posts []dto.PostDTO, err error) {
	var entities []*entity.Post

	if lastPosition == 0 {
		entities, err = s.repository.PostRepository.GetPublicTimeline()
		if err != nil {
			return nil, err
		}
	} else {
		entities, err = s.repository.PostRepository.GetMorePublicTimeline(lastPosition)
		if err != nil {
			return nil, err
		}
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
		post.Poster = dto.UserDTO{
			ID:       *entity.CreatedBy,
			Username: entity.Username,
		}

		posts = append(posts, *post)
	}

	return posts, err
}

func (s *PostService) Populate(content []string) error {
	err := s.repository.PostRepository.CreateBatch(content)
	if err != nil {
		return err
	}
	return nil
}
