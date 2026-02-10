package model

type CreateOpts struct {
	Namespace  string
	SecretName string
	Values     map[string]string
}
