package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/dreambrother/oas-ls/internal/handler"
	"github.com/owenrumney/go-lsp/server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	srv := server.NewServer(handler.NewHandler(logger), server.WithLogger(logger))
	if err := srv.Run(context.Background(), server.RunStdio()); err != nil {
		log.Fatal(err)
	}
}
