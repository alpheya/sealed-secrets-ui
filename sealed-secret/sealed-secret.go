package sealedsecret

import (
	"crypto/rsa"
	"fmt"

	"gopkg.in/yaml.v2"
)

type EncryptValuesParams struct {
	PubKey    *rsa.PublicKey
	Scope     string
	Namespace string
	Name      string
	Values    map[string]string
}

type Metadata struct {
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type SealedSecret struct {
	APIVersion string           `yaml:"apiVersion"`
	Kind       string           `yaml:"kind"`
	Metadata   Metadata         `yaml:"metadata"`
	Spec       SealedSecretSpec `yaml:"spec"`
}

type SealedSecretSpec struct {
	EncryptedData map[string]string `yaml:"encryptedData"`
	Template      Template          `yaml:"template"`
}

type Template struct {
	Metadata Metadata `yaml:"metadata,omitempty"`
}

type GetLabelParams struct {
	Scope     string
	Namespace string
	Name      string
}

func getLabel(params GetLabelParams) string {
	switch params.Scope {
	case "cluster":
		return ""
	case "namespace":
		return params.Namespace
	default:
		return fmt.Sprintf("%s/%s", params.Namespace, params.Name)
	}
}

func encryptValues(params EncryptValuesParams) (map[string]string, error) {
	encryptedData := make(map[string]string)
	for key, value := range params.Values {
		enc, err := HybridEncrypt(EncryptParams{
			PublicKey: params.PubKey,
			Value:     value,
			Label:     getLabel(GetLabelParams{Scope: params.Scope, Namespace: params.Namespace, Name: params.Name}),
		})

		if err != nil {
			return nil, err
		}

		encryptedData[key] = enc
	}

	return encryptedData, nil
}

func GetSealedSecret(params EncryptValuesParams) (string, error) {
	encryptedData, err := encryptValues(params)
	if err != nil {
		return "", err
	}

	annotations := make(map[string]string)
	if params.Scope == "cluster" {
		annotations["sealedsecrets.bitnami.com/cluster-wide"] = "true"
	} else if params.Scope == "namespace" {
		annotations["sealedsecrets.bitnami.com/namespace-wide"] = "true"
	}

	sealedSecret := SealedSecret{
		APIVersion: "bitnami.com/v1alpha1",
		Kind:       "SealedSecret",
		Metadata: Metadata{
			Name:        params.Name,
			Namespace:   params.Namespace,
			Annotations: annotations,
		},
		Spec: SealedSecretSpec{
			EncryptedData: encryptedData,
			Template: Template{
				Metadata: Metadata{
					Name:        params.Name,
					Namespace:   params.Namespace,
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
