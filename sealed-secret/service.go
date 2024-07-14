package sealedsecret

import (
	"context"
	"crypto/rsa"
	"fmt"

	"github.com/quantum-wealth/sealed-secrets-ui/model"
	"gopkg.in/yaml.v2"
)

type SealedSecretService struct {
	sealedSecretControllerName      string
	sealedSecretControllerNamespace string
}

type encryptRequest struct {
	pubKey     *rsa.PublicKey
	secretName string
	namespace  string
	scope      string
	values     map[string]string
}

func NewSealedSecretService(controllerNamespace, controllerName string) SealedSecretService {
	return SealedSecretService{
		sealedSecretControllerNamespace: controllerNamespace,
		sealedSecretControllerName:      controllerName,
	}
}

func (s SealedSecretService) CreateSealedSecret(ctx context.Context, opts model.CreateOpts) (string, error) {
	existingData, err := getSecretData(ctx, opts.Namespace, opts.SecretName)
	if err != nil {
		return "", fmt.Errorf("failed to get existing secret data: %w", err)
	}

	// we need to encrypt all the existing data as well as the new data
	valuesToEncrypt := make(map[string]string)
	for key, value := range existingData {
		valuesToEncrypt[key] = value
	}

	for key, value := range opts.Values {
		valuesToEncrypt[key] = value
	}

	// we need to get the public key every time we create a sealed secret because the
	// sealed-secrets controller rotates the public key every X hours
	pubKey, err := s.getPublicKey(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get public key: %w", err)
	}

	req := encryptRequest{
		pubKey:     pubKey,
		secretName: opts.SecretName,
		namespace:  opts.Namespace,
		values:     valuesToEncrypt,
		scope:      opts.Scope,
	}

	encryptedData, err := s.encryptValues(req)
	if err != nil {
		return "", err
	}

	annotations := make(map[string]string)
	if req.scope == "cluster" {
		annotations["sealedsecrets.bitnami.com/cluster-wide"] = "true"
	} else if req.scope == "namespace" {
		annotations["sealedsecrets.bitnami.com/namespace-wide"] = "true"
	}
	// if scope == strict we not need any annotations

	sealedSecret := model.SealedSecret{
		APIVersion: "bitnami.com/v1alpha1",
		Kind:       "SealedSecret",
		Metadata: model.Metadata{
			Name:        req.secretName,
			Namespace:   req.namespace,
			Annotations: annotations,
		},
		Spec: model.SealedSecretSpec{
			EncryptedData: encryptedData,
			Template: model.Template{
				Metadata: model.Metadata{
					Name:        req.secretName,
					Namespace:   req.namespace,
					Annotations: annotations,
				},
			},
		},
	}

	yamlData, err := yaml.Marshal(sealedSecret)
	if err != nil {
		return "", fmt.Errorf("failed to marshal sealed secret to YAML: %w", err)
	}

	return string(yamlData), nil
}

func (s SealedSecretService) getLabel(req encryptRequest) string {
	switch req.scope {
	case "cluster":
		return ""
	case "namespace":
		return req.namespace
	default:
		return fmt.Sprintf("%s/%s", req.namespace, req.secretName)
	}
}

func (s SealedSecretService) encryptValues(req encryptRequest) (map[string]string, error) {
	encryptedData := make(map[string]string)
	for key, value := range req.values {
		enc, err := hybridEncrypt(req.pubKey, value, s.getLabel(req))
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt value: %w", err)
		}

		encryptedData[key] = enc
	}

	return encryptedData, nil
}
