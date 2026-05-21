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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app := server.NewApp(ctx, server.LoadConfig(rootDir))
	err = app.Start()

	if err != nil {
		panic(err)
	}
}
