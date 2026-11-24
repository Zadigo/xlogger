package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Zadigo/goxlogger/internal/backend"
	"github.com/Zadigo/goxlogger/internal/mainapp"
	"github.com/Zadigo/goxlogger/internal/utils"
	"github.com/joho/godotenv"
)

func preStartup(ctx context.Context, options *utils.ServerOptions) {
	options.Redis = backend.NewRedisBackend(ctx)
}

func main() {
	godotenv.Load(".env")
	log.Print("⚡️ Starting Go-xLogger...")

	rootDir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ctx = context.WithValue(ctx, "rootDir", rootDir)
	ctx = context.WithValue(ctx, "debugMode", os.Getenv("DEBUG_MODE") == "true")

	config := &utils.ServerOptions{}
	config.Load(ctx)

	preStartup(ctx, config)
	app := mainapp.NewMainServerApp(ctx, config)
	err = app.Start(ctx)
	if err != nil {
		log.Print(err)
	}
}
