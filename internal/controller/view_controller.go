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

		for i := 0; i < 25; i++ {
			homeDTO.Posts = append(homeDTO.Posts, dto.PostDTO{
				ID: strconv.Itoa(i + 1),
				Message: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus sed justo eget metus interdum dapibus. Quisque nisi sem, molestie a luctus sit amet, pellentesque a justo. In elementum malesuada odio sed finibus. Nunc accumsan metus in velit malesuada, quis egestas nisl venenatis. Sed aliquam nisl eu mi lobortis, quis euismod libero commodo. Mauris pulvinar id velit sed porttitor. Nullam vitae est ullamcorper libero venenatis imperdiet. Maecenas fringilla cursus nibh vitae varius. Donec a posuere ex. Etiam ipsum nulla, molestie vulputate leo sit amet, fermentum aliquam quam.

				Praesent ullamcorper ex in turpis congue, at vestibulum mi interdum. Phasellus sodales orci a massa semper, in finibus metus elementum. Fusce mi enim, imperdiet eu dignissim eget, consequat vel sapien. Phasellus lacus dui, euismod sed tempus ac, interdum in sem. Duis maximus ornare orci, sodales venenatis sapien iaculis et. Nunc sed interdum nibh. Cras interdum risus non tincidunt lobortis. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Aliquam congue lectus non ante porta accumsan. In porta sit amet ipsum vel hendrerit. Nullam massa nunc, tincidunt a auctor at, egestas non tortor.
				
				Proin placerat, risus eget rutrum euismod, elit risus consequat dolor, a tristique turpis purus a elit. Praesent a magna quis justo ullamcorper iaculis. Curabitur a dictum leo. Nullam condimentum interdum nisi vitae aliquet. Etiam lacus dui, aliquam at fringilla et, facilisis eget dui. Ut sapien turpis, porttitor eget eros eget, tincidunt venenatis libero. Maecenas viverra consequat ullamcorper. Fusce in scelerisque tortor. Nulla non luctus sem. Ut ac porttitor ex. Vestibulum finibus rutrum dui eu dictum. Vivamus non accumsan lectus, a tincidunt diam. Curabitur sit amet nulla augue. Vivamus commodo arcu at nunc mattis bibendum. Curabitur sed convallis ex. Maecenas in mauris ut turpis consectetur iaculis.
				
				In nec metus sed urna imperdiet molestie. Curabitur in mattis velit, ac aliquet lorem. Sed posuere ac nisi mollis venenatis. Phasellus sit amet nisi tempor, semper orci a, finibus nisi. Fusce a neque non lectus auctor facilisis a bibendum turpis. Nullam auctor nisi quis pharetra viverra. Suspendisse luctus, libero at posuere tristique, dolor massa tristique lectus, a commodo augue purus pulvinar ante. Fusce magna tortor, vulputate vel commodo in, tincidunt et magna. Etiam imperdiet libero eget nunc aliquet, id rhoncus tortor sodales. Nulla tortor tortor, bibendum eu dapibus a, scelerisque nec libero.
				
				Praesent eget metus tristique, dictum dui quis, tincidunt ex. Curabitur interdum non velit ac facilisis. Fusce mi felis, ultricies vitae porta mollis, gravida pulvinar mi. Nulla facilisi. Vestibulum non scelerisque massa. Integer hendrerit magna in risus pulvinar efficitur vitae eu massa. Cras molestie ante nec gravida luctus. Nam dapibus velit quis rutrum suscipit. Praesent placerat purus a lorem pellentesque, ut lobortis tortor porta. Praesent tempor odio lectus, eget viverra purus vehicula sed.
				
				`,
				Replies:     []dto.PostDTO{},
				ReplyCounts: i,
				Likes:       i,
				SavedCounts: i,
				Impressions: i,
				Poster: dto.UserDTO{
					Username: "[deleted_user]",
				},
				PostedAt: time.Now(),
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
		store, _ := vc.Session.Get(r, "default")
		user := store.Values["user"].(*dto.UserDTO)

		createNewPostDTO := &dto.CreateNewPostDTO{
			User: user,
		}

		view.CreateNewPostPage(createNewPostDTO).Render(r.Context(), w)
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
				Message:     fmt.Sprintf("this is my %d post", i+1),
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
