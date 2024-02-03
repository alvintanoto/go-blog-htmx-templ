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
		store, _ := vc.Session.Get(r, "default")
		user := store.Values["user"]
		if user != nil {
			view.NotFoundPage(user.(*dto.UserDTO)).Render(r.Context(), w)
			return
		}

		view.NotFoundPage(nil).Render(r.Context(), w)
	}
}

func (vc *ViewController) SignInHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")
		dto := &dto.SignInPageDTO{}

		for _, flash := range store.Flashes("error") {
			dto.Error = flash.(string)
		}

		sessions.Save(r, w)
		view.SignInPage(dto).Render(r.Context(), w)
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
				Error: "",
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

			for _, flash := range store.Flashes("error") {
				payload.Error = flash.(string)
			}

			err := sessions.Save(r, w)
			if err != nil {
				log.Println("error saving session :", err.Error())
			}
			view.RegisterPage(payload).Render(r.Context(), w)
			return
		}

		view.RegisterPage(nil).Render(r.Context(), w)
	}
}

func (vc *ViewController) HomepageViewHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		homeDTO := &dto.HomepageDTO{
			Posts: []dto.PostDTO{},
		}

		store, _ := vc.Session.Get(r, "default")
		userStore := store.Values["user"]
		if userStore != nil {
			homeDTO.User = userStore.(*dto.UserDTO)
		}

		view.Homepage(homeDTO).Render(r.Context(), w)
	}
}

func (vc *ViewController) CreateNewPostViewHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")
		user := store.Values["user"].(*dto.UserDTO)

		createNewPostDTO := &dto.CreateNewPostDTO{
			User: user,
		}

		view.CreateNewPostPage(createNewPostDTO).Render(r.Context(), w)
	}
}

func (vc *ViewController) CreatePostHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (vc *ViewController) ProfileHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")
		user := store.Values["user"].(*dto.UserDTO)

		profileDTO := &dto.ProfilePageDTO{
			User:  user,
			Posts: []dto.PostDTO{},
		}

		for i := 0; i < 25; i++ {
			profileDTO.Posts = append(profileDTO.Posts, dto.PostDTO{
				ID:          strconv.Itoa(i + 1),
				Content:     fmt.Sprintf("this is my %d post", i+1),
				Replies:     []dto.PostDTO{},
				ReplyCounts: i,
				Likes:       i,
				SavedCounts: i,
				Impressions: i,
				Poster: dto.UserDTO{
					Username: user.Username,
				},
				PostedAt: time.Now(),
			})
		}

		view.ProfilePage(profileDTO).Render(r.Context(), w)
	}
}

func (vc *ViewController) SettingsHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")
		user := store.Values["user"].(*dto.UserDTO)

		view.SettingsPage(&dto.SettingsPageDto{
			User: user,
		}).Render(r.Context(), w)
	}
}
