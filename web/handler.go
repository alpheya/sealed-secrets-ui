package web

import (
	"net/http"

	"github.com/a-h/templ"
)

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", templ.Handler(Home()))

	return mux
}
