package web

import (
	"embed"
	"net/http"
	"time"

	"github.com/a-h/templ"
)

//go:embed *.gif
var spinnerFiles embed.FS

func CreateSealedSecret(w http.ResponseWriter, r *http.Request) {
	// sleep 2 seconds
	time.Sleep(4 * time.Second)

	CodeArea("asd").Render(r.Context(), w)
}

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/spinner.gif", http.FileServer(http.FS(spinnerFiles)))
	mux.HandleFunc("/sealed-secret", CreateSealedSecret)
	mux.Handle("/", templ.Handler(Home()))

	return mux
}
