package main

import (
	"os"
	"time"

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
	handler := web.NewHandler()

	web.Start("8080", handler)
}
