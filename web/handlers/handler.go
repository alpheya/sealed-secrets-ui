package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/quantum-wealth/sealed-secrets-ui/model"
	"github.com/quantum-wealth/sealed-secrets-ui/web/ui"
	"github.com/rs/zerolog/log"
)

type sealer interface {
	CreateSealedSecret(context.Context, model.CreateOpts) (string, error)
}

type SealedSecretHandler struct {
	svc sealer
}

func NewSealedSecretHandler(svc sealer) SealedSecretHandler {
	return SealedSecretHandler{svc: svc}
}

func (s SealedSecretHandler) CreateSealedSecretHandler(w http.ResponseWriter, r *http.Request) {
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

	createOpts := model.CreateOpts{
		Scope:      scope,
		Namespace:  namespace,
		SecretName: secretName,
		Values:     keyValues,
	}

	yamlManifest, err := s.svc.CreateSealedSecret(r.Context(), createOpts)

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
