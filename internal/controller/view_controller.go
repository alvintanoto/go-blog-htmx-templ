package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"alvintanoto.id/blog-htmx-templ/internal/dto"
	"alvintanoto.id/blog-htmx-templ/view"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

type ViewController struct {
	Session *sessions.CookieStore
}

func NewViewController(session *sessions.CookieStore) *ViewController {
	return &ViewController{
		Session: session,
	}
}

func (vc *ViewController) NotFoundViewHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("path:", r.URL.Path, "not found")
		view.NotFoundPage(nil).Render(r.Context(), w)
	}
}

func (vc *ViewController) SignInHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		view.SignInPage().Render(r.Context(), w)
	}
}

func (vc *ViewController) RegisterHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")

		var payload dto.RegisterUserRequestDTO
		decoder := schema.NewDecoder()

		r.ParseForm()
		err := decoder.Decode(&payload, r.PostForm)
		if err == nil {
			payload := &dto.RegisterPageDTO{
				RegisterFieldDTO: &dto.RegisterFieldDTO{
					Username: struct {
						Value  string
						Errors []string
					}{
						Value: payload.Username,
					},
					Email: struct {
						Value  string
						Errors []string
					}{Value: payload.Email},
					PasswordErrors:        []string{},
					ConfirmPasswordErrors: []string{},
				},
			}

			// redirect result check for flashes

			for _, flash := range store.Flashes("username") {
				payload.RegisterFieldDTO.Username.Errors = append(payload.RegisterFieldDTO.Username.Errors, flash.(string))
			}

			for _, flash := range store.Flashes("email") {
				payload.RegisterFieldDTO.Email.Errors = append(payload.RegisterFieldDTO.Email.Errors, flash.(string))
			}

			for _, flash := range store.Flashes("password") {
				payload.RegisterFieldDTO.PasswordErrors = append(payload.RegisterFieldDTO.PasswordErrors, flash.(string))
			}

			for _, flash := range store.Flashes("confirm_password") {
				payload.RegisterFieldDTO.ConfirmPasswordErrors = append(payload.RegisterFieldDTO.ConfirmPasswordErrors, flash.(string))
			}

			err := sessions.Save(r, w)
			if err != nil {
				log.Println("error saving session :", err.Error())
			}
			view.RegisterPage(payload).Render(r.Context(), w)
			return
		}

		fmt.Println("fresh page")
		view.RegisterPage(nil).Render(r.Context(), w)
	}
}

func (vc *ViewController) HomepageViewHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		homeDTO := &dto.HomepageDTO{
			Posts: []dto.PostDTO{},
		}

		for i := 0; i < 25; i++ {
			homeDTO.Posts = append(homeDTO.Posts, dto.PostDTO{
				ID:          strconv.Itoa(i + 1),
				Message:     fmt.Sprintf("this is my %d post", i+1),
				Replies:     []dto.PostDTO{},
				ReplyCounts: i,
				Likes:       i,
				SavedCounts: i,
				Impressions: i,
				Poster:      dto.UserDTO{},
				PostedAt:    time.Now(),
			})
		}

		view.Homepage(homeDTO).Render(r.Context(), w)
	}
}

func (vc *ViewController) CreateNewPostHandler() func(http.ResponseWriter, *http.Request) {
	// TODO: check session
	// if not logged in redirect to sign in page
	// else show create new post page

	return func(w http.ResponseWriter, r *http.Request) {
		createNewPostDTO := &dto.CreateNewPostDTO{
			User: nil,
		}

		view.CreateNewPostPage(createNewPostDTO).Render(r.Context(), w)
	}
}
