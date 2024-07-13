package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"

	"time"

	"github.com/rs/zerolog/log"
)

func recoverer(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func(ctx context.Context) {
			if rvr := recover(); rvr != nil {
				log.Error().Ctx(ctx).Err(fmt.Errorf("recovered from panic: %v\nstack: %s", rvr, string(debug.Stack())))

				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(writer, "{\"error\":\"%s\"}", http.StatusText(http.StatusInternalServerError))
			}
		}(request.Context())
		next.ServeHTTP(writer, request)
	})
}

func shutdown(server *http.Server) {
	gracePeriod := 25 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Panic().Err(err).Msg("Server shutdown failed")
	}
	log.Warn().Msg("Server shutdown")
}

func start(server *http.Server) {
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Panic().Err(err).Msg("Server failed to start")
	}
	log.Info().Msg("Server stopped")
}

func Start(port string, handler http.Handler) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	wrappedHandler := recoverer(handler)

	addr := fmt.Sprintf(":%s", port)

	srv := &http.Server{
		Addr:              addr,
		Handler:           http.TimeoutHandler(wrappedHandler, 10*time.Second, "request timed out"),
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	log.Info().Msgf("Server listening on %s", addr)

	go start(srv)

	<-ctx.Done()

	stop()
	shutdown(srv)
}
