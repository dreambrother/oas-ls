package main

import (
	"context"
	"log"

	"github.com/dreambrother/specpls/internal/handler"
	"github.com/owenrumney/go-lsp/server"
)

func main() {
	srv := server.NewServer(handler.NewHandler())
	if err := srv.Run(context.Background(), server.RunStdio()); err != nil {
		log.Fatal(err)
	}
}
