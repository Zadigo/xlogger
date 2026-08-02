package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/Zadigo/goxlogger/internal/httpapp"
	"github.com/Zadigo/goxlogger/internal/utils"
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

	ctx = context.WithValue(ctx, "rootDir", rootDir)
	ctx = context.WithValue(ctx, "debugMode", os.Getenv("DEBUG_MODE") == "true")

	config := &utils.ServerConfig{}
	config.Load(ctx)

	app := httpapp.NewApp(ctx)
	err = app.Start()

	if err != nil {
		panic(err)
	}
}
