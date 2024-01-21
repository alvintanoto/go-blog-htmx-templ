package controller

import (
	"log"
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
		log.Println("path:", r.URL.Path, "not found")
		view.NotFoundPage().Render(r.Context(), w)
	}
}
