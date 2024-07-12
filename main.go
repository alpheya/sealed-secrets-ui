package main

import (
	"fmt"

	"github.com/quantum-wealth/sealed-secrets-ui/k8s"
	sealedsecret "github.com/quantum-wealth/sealed-secrets-ui/sealed-secret"
)

func main() {
	pubKey, err := k8s.GetPublicKey()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	params := sealedsecret.EncryptValuesParams{
		PubKey:    pubKey,
		Namespace: "default",
		Scope:     "",
		Name:      "my-sealed-secret",
		Values:    map[string]string{"password": "supersecret"},
	}

	yamlManifest, err := sealedsecret.GetSealedSecret(params)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(yamlManifest)
}
