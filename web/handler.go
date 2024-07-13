package web

import (
	"embed"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/quantum-wealth/sealed-secrets-ui/k8s"
	sealedsecret "github.com/quantum-wealth/sealed-secrets-ui/sealed-secret"
	"github.com/rs/zerolog/log"
)

//go:embed *.gif
var spinnerFiles embed.FS

func parseKeyValuePairs(data string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

func CreateSealedSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	scope := r.FormValue("scope")
	namespace := r.FormValue("namespace")
	secretName := r.FormValue("secretName")
	valuesToEncrypt := r.FormValue("values")

	log.Info().Str("scope", scope).Str("namespace", namespace).Str("secretName", secretName).Msg("Creating sealed secret")
	keyValues := parseKeyValuePairs(valuesToEncrypt)

	pubKey, err := k8s.GetPublicKey()
	if err != nil {
		http.Error(w, "Error getting public key", http.StatusInternalServerError)
		return
	}

	params := sealedsecret.EncryptValuesParams{
		PubKey:    pubKey,
		Namespace: namespace,
		Scope:     scope,
		Name:      secretName,
		Values:    keyValues,
	}

	yamlManifest, err := sealedsecret.GetSealedSecret(params)
	if err != nil {
		http.Error(w, "Error creating sealed secret", http.StatusInternalServerError)
		return
	}

	CodeArea(yamlManifest).Render(r.Context(), w)
}

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/spinner.gif", http.FileServer(http.FS(spinnerFiles)))
	mux.HandleFunc("/sealed-secret", CreateSealedSecret)
	mux.Handle("/", templ.Handler(Home()))

	return mux
}
