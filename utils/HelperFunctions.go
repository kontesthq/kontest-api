package utils

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func HandleShutdown(server *http.Server) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		slog.Info("Shutting down server...")
		// Create a context with a timeout to allow for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			slog.Error("Server shutdown failed", slog.String("error", err.Error()))
			os.Exit(1)
		}

		slog.Info("Server exited properly")
	}()
}
