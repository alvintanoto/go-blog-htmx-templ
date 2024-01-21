package controller

import (
	"net/http"

	"alvintanoto.id/blog-htmx-templ/view"
)

type ViewController struct {
}

func NewViewController() *ViewController {
	return &ViewController{}
}

func (vc ViewController) NotFoundViewHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		view.NotFoundPage().Render(r.Context(), w)
	}
}
