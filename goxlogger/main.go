package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/Zadigo/goxlogger/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	log.Print("⚡️ Starting Go-xLogger...")

	rootDir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	app := server.NewApp(server.LoadConfig(rootDir))
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err = app.Start(ctx)
	if err != nil {
		panic(err)
	}
}
