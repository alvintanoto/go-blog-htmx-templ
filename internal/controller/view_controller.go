package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"alvintanoto.id/blog-htmx-templ/internal/dto"
	"alvintanoto.id/blog-htmx-templ/view"
)

type ViewController struct {
}

func NewViewController() *ViewController {
	return &ViewController{}
}

func (vc *ViewController) HomepageViewHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		homeDTO := &dto.HomepageDTO{
			Posts: []dto.PostDTO{},
		}

		for i := 0; i < 500; i++ {
			homeDTO.Posts = append(homeDTO.Posts, dto.PostDTO{
				ID:          strconv.Itoa(i + 1),
				Message:     fmt.Sprintf("this is my %d post", i+1),
				Replies:     []dto.PostDTO{},
				ReplyCounts: i,
				Likes:       i,
				Impressions: i,
				Poster:      dto.UserDTO{},
				PostedAt:    time.Now(),
			})
		}

		view.Homepage(homeDTO).Render(r.Context(), w)
	}
}

func (vc *ViewController) NotFoundViewHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("path:", r.URL.Path, "not found")
		view.NotFoundPage().Render(r.Context(), w)
	}
}
