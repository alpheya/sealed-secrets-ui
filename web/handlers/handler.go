package handlers

import (
	"net/http"
	"strings"

	"github.com/quantum-wealth/sealed-secrets-ui/k8s"
	sealedsecret "github.com/quantum-wealth/sealed-secrets-ui/sealed-secret"
	"github.com/quantum-wealth/sealed-secrets-ui/web/ui"
	"github.com/rs/zerolog/log"
)

func parseKeyValuePairs(data string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(data, "\n")

	if len(lines) == 0 {
		return nil
	}

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

	if scope == "" || namespace == "" || secretName == "" || valuesToEncrypt == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	log.Info().Str("scope", scope).Str("namespace", namespace).Str("secretName", secretName).Msg("creating sealed secret")
	keyValues := parseKeyValuePairs(valuesToEncrypt)

	if keyValues == nil {
		http.Error(w, "Error parsing key value pairs", http.StatusBadRequest)
		return
	}

	pubKey, err := k8s.GetPublicKey()
	if err != nil {
		log.Ctx(r.Context()).Err(err).Msg("error getting public key")
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
	log.Info().Str("yaml", yamlManifest).Msg("sealed-secret created")

	if err != nil {
		log.Ctx(r.Context()).Err(err).Msg("error creating sealed secret")
		http.Error(w, "Error creating sealed secret", http.StatusInternalServerError)
		return
	}

	err = ui.CodeArea(yamlManifest).Render(r.Context(), w)
	if err != nil {
		log.Err(err).Msg("error rendering code area")
		http.Error(w, "Error rendering code area", http.StatusInternalServerError)
		return
	}
}
