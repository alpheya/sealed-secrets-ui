package web

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/quantum-wealth/sealed-secrets-ui/web/assets"
	"github.com/quantum-wealth/sealed-secrets-ui/web/handlers"
	"github.com/quantum-wealth/sealed-secrets-ui/web/ui"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/spinner.gif", http.FileServer(http.FS(assets.SpinnerFiles)))
	mux.HandleFunc("/sealed-secret", handlers.CreateSealedSecret)
	mux.Handle("/", templ.Handler(ui.Home()))

	return mux
}
