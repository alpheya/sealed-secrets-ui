package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/alpheya/sealed-secrets-ui/model"
	"github.com/alpheya/sealed-secrets-ui/web/ui"
	"github.com/rs/zerolog/log"
)

var escapedBacktick = strings.Join([]string{`\`, "`"}, "")

type sealer interface {
	CreateSealedSecret(context.Context, model.CreateOpts) (string, error)
}

type SealedSecretHandler struct {
	svc sealer
}

func NewSealedSecretHandler(svc sealer) SealedSecretHandler {
	return SealedSecretHandler{svc: svc}
}

func respondError(w http.ResponseWriter, message string) {
	w.Header().Set("HX-Retarget", ".message")

	err := ui.Error(message).Render(context.Background(), w)
	if err != nil {
		log.Err(err).Msg("error rendering error message")
		http.Error(w, "Error rendering error message", http.StatusInternalServerError)
	}
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
		respondError(w, "All fields are required")
		return
	}

	log.Info().Str("scope", scope).Str("namespace", namespace).Str("secretName", secretName).Msg("creating sealed secret")
	keyValues, err := parseKeyValuePairs(valuesToEncrypt)
	if err != nil {
		respondError(w, fmt.Sprintf("Wrongly formatted value(s): %v", err.Error()))
		return
	}

	if keyValues == nil {
		respondError(w, "No key-value pairs found")
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
		respondError(w, "Error creating sealed secret")
		return
	}

	err = ui.CodeArea(yamlManifest).Render(r.Context(), w)
	if err != nil {
		log.Err(err).Msg("error rendering code area")
		http.Error(w, "Error rendering code area", http.StatusInternalServerError)
		return
	}
}

func parseKeyValuePairs(data string) (map[string]string, error) {
	result := make(map[string]string)
	lines := strings.Split(data, "\n")

	if len(lines) == 1 && len(lines[0]) == 0 {
		return nil, errors.New("empty")
	}

	var multilineKey string
	var multilineValue strings.Builder

	for i, line := range lines {

		// inside backticked block
		if len(multilineKey) > 0 {
			var isEndOfBlock bool
			if !strings.HasSuffix(line, escapedBacktick) {
				line, isEndOfBlock = strings.CutSuffix(line, "`")
			}
			line = strings.ReplaceAll(line, escapedBacktick, "`")
			multilineValue.WriteByte('\n')
			multilineValue.WriteString(line)
			if isEndOfBlock {
				result[multilineKey] = multilineValue.String()
				multilineValue.Reset()
				multilineKey = ""
			}
			continue
		}

		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			return nil, fmt.Errorf("Missing '=' at line: %v", i)
		}

		// backticked block starts
		part, ok := strings.CutPrefix(parts[1], "`")
		if ok {
			multilineKey = parts[0]

			var isEndOfBlock bool
			if !strings.HasSuffix(part, escapedBacktick) {
				part, isEndOfBlock = strings.CutSuffix(part, "`")
			}
			part = strings.ReplaceAll(part, escapedBacktick, "`")
			multilineValue.WriteString(part)
			if isEndOfBlock {
				result[multilineKey] = multilineValue.String()
				multilineValue.Reset()
				multilineKey = ""
			}
			continue
		}

		// oneline value
		result[parts[0]] = parts[1]
	}

	if len(multilineKey) != 0 {
		return nil, fmt.Errorf("Backticked block is not closed")
	}

	return result, nil
}
