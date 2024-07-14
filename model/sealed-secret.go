package model

type Metadata struct {
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type SealedSecretSpec struct {
	EncryptedData map[string]string `yaml:"encryptedData"`
	Template      Template          `yaml:"template"`
}

type Template struct {
	Metadata Metadata `yaml:"metadata,omitempty"`
}

type SealedSecret struct {
	APIVersion string           `yaml:"apiVersion"`
	Kind       string           `yaml:"kind"`
	Metadata   Metadata         `yaml:"metadata"`
	Spec       SealedSecretSpec `yaml:"spec"`
}
