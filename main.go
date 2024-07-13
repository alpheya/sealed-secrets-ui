package main

import (
	"fmt"
	"os"
	"time"

	"github.com/quantum-wealth/sealed-secrets-ui/k8s"
	sealedsecret "github.com/quantum-wealth/sealed-secrets-ui/sealed-secret"
	"github.com/quantum-wealth/sealed-secrets-ui/web"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setupLogging() {
	logLevel := os.Getenv("LOG_LEVEL")
	level, _ := zerolog.ParseLevel(logLevel) //nolint: errcheck
	if level == zerolog.NoLevel {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = time.RFC3339
	l := log.Level(level)
	zerolog.DefaultContextLogger = &l
}

func main() {
	setupLogging()
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

	handler := web.NewHandler()

	web.Start("8080", handler)
}
