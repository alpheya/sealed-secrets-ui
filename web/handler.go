package web

import (
	"embed"
	"net/http"
	"time"

	"github.com/a-h/templ"
)

//go:embed *.gif
var spinnerFiles embed.FS

func Post(w http.ResponseWriter, r *http.Request) {
	// sleep 2 seconds
	time.Sleep(200 * time.Second)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "success"}`))

}

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/spinner.gif", http.FileServer(http.FS(spinnerFiles)))
	mux.HandleFunc("/sealed-secret", Post)
	mux.Handle("/", templ.Handler(Home()))

	return mux
}
