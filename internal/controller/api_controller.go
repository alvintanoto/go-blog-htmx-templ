package controller

import (
	"log"
	"net/http"

	"alvintanoto.id/blog-htmx-templ/internal/dto"
	"alvintanoto.id/blog-htmx-templ/internal/service"
	"alvintanoto.id/blog-htmx-templ/view"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
)

type ApiController struct {
	Store   *sessions.CookieStore
	Service *service.Service
}

func NewApiController(service *service.Service, store *sessions.CookieStore) *ApiController {
	return &ApiController{
		Service: service,
		Store:   store,
	}
}

func (ac *ApiController) PreviewPost() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := schema.NewDecoder()

		var payload dto.PreviewPostDTO
		err := r.ParseForm()
		if err != nil {
			log.Println("error parsing form: ", err.Error())
			http.Redirect(w, r, "/post/new-post", http.StatusPermanentRedirect)
			return
		}

		err = decoder.Decode(&payload, r.PostForm)
		if err != nil {
			log.Println("error decoding payload: ", err.Error())
			http.Redirect(w, r, "/post/new-post", http.StatusPermanentRedirect)
			return
		}

		if payload.Preview {
			maybeUnsafeHTML := markdown.ToHTML([]byte(payload.Content), nil, nil)
			html := bluemonday.UGCPolicy().SanitizeBytes(maybeUnsafeHTML)

			view.PreviewPostContainer(string(html), payload.Content).Render(r.Context(), w)
			return
		}

		view.EditorContainer(payload.Content).Render(r.Context(), w)
	}
}
