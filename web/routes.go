package web

import (
	"net/http"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/a-h/templ"
	sealedsecret "github.com/quantum-wealth/sealed-secrets-ui/sealed-secret"
	"github.com/quantum-wealth/sealed-secrets-ui/web/assets"
	"github.com/quantum-wealth/sealed-secrets-ui/web/handlers"
	"github.com/quantum-wealth/sealed-secrets-ui/web/ui"
)

func NewRouter() http.Handler {
	controllerNamespace := os.Getenv("SEALED_SECRETS_CONTROLLER_NAMESPACE")
	controllerName := os.Getenv("SEALED_SECRETS_CONTROLLER_NAME")

	if controllerNamespace == "" {
		controllerNamespace = "kube-system" // default namespace if sealed-secrets was installed with Helm
	}

	if controllerName == "" {
		controllerName = "sealed-secrets-controller" // default controllerName if sealed-secrets was installed with Helm
	}

	svc, err := sealedsecret.NewSealedSecretService(controllerNamespace, controllerName)
	if err != nil {
		log.Panic().Err(err).Msg("failed to create sealed secret service")
	}

	handler := handlers.NewSealedSecretHandler(svc)

	mux := http.NewServeMux()
	mux.Handle("/spinner.gif", http.FileServer(http.FS(assets.SpinnerFiles)))
	mux.HandleFunc("/sealed-secret", handler.CreateSealedSecretHandler)
	mux.HandleFunc("/healthz", handlers.HealthHandler)
	mux.Handle("/", templ.Handler(ui.Home()))

	return mux
}
