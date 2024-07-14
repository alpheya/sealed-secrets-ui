package model

type CreateOpts struct {
	Scope      string
	Namespace  string
	SecretName string
	Values     map[string]string
}
