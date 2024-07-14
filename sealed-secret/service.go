package sealedsecret

import (
	"context"
	"crypto/rsa"
	"fmt"

	"github.com/quantum-wealth/sealed-secrets-ui/model"
	"gopkg.in/yaml.v2"
)

type SealedSecretService struct {
	pubKey                          *rsa.PublicKey
	sealedSecretControllerName      string
	sealedSecretControllerNamespace string
}

func NewSealedSecretService(controllerNamespace, controllerName string) SealedSecretService {
	return SealedSecretService{
		sealedSecretControllerNamespace: controllerNamespace,
		sealedSecretControllerName:      controllerName,
	}
}

func (s SealedSecretService) CreateSealedSecret(ctx context.Context, opts model.CreateOpts) (string, error) {
	// we need to get the public key every time we create a sealed secret because the
	// sealed-secrets controller rotates the public key every X hours
	pubKey, err := s.getPublicKey(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get public key: %w", err)
	}

	s.pubKey = pubKey

	encryptedData, err := s.encryptValues(opts)
	if err != nil {
		return "", err
	}

	annotations := make(map[string]string)
	if opts.Scope == "cluster" {
		annotations["sealedsecrets.bitnami.com/cluster-wide"] = "true"
	} else if opts.Scope == "namespace" {
		annotations["sealedsecrets.bitnami.com/namespace-wide"] = "true"
	}
	// if scope == strict we not need any annotations

	sealedSecret := model.SealedSecret{
		APIVersion: "bitnami.com/v1alpha1",
		Kind:       "SealedSecret",
		Metadata: model.Metadata{
			Name:        opts.SecretName,
			Namespace:   opts.Namespace,
			Annotations: annotations,
		},
		Spec: model.SealedSecretSpec{
			EncryptedData: encryptedData,
			Template: model.Template{
				Metadata: model.Metadata{
					Name:        opts.SecretName,
					Namespace:   opts.Namespace,
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

func (s SealedSecretService) getLabel(opts model.CreateOpts) string {
	switch opts.Scope {
	case "cluster":
		return ""
	case "namespace":
		return opts.Namespace
	default:
		return fmt.Sprintf("%s/%s", opts.Namespace, opts.SecretName)
	}
}

func (s SealedSecretService) encryptValues(opts model.CreateOpts) (map[string]string, error) {
	encryptedData := make(map[string]string)
	for key, value := range opts.Values {
		enc, err := s.hybridEncrypt(value, s.getLabel(opts))
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt value: %w", err)
		}

		encryptedData[key] = enc
	}

	return encryptedData, nil
}
