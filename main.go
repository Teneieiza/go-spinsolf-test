package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Teneieiza/go-spinsolf-test/app"
	"github.com/Teneieiza/go-spinsolf-test/config"
)

func main() {
	cfg := config.LoadConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := config.InitDatabase(ctx, cfg); err != nil {
		log.Fatal(err)
	}
	defer config.DB.Close(context.Background())

	app := app.NewApplication(cfg)

	if err := app.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")
	app.Shutdown()
}
