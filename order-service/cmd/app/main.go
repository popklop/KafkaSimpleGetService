package main

import (
	"context"
	"log"
	"wbtech/internal/app"
	"wbtech/internal/config"
)

func main() {
	cfg := config.Load()
	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	if err := application.Run(context.Background()); err != nil {
		log.Fatalf("app run error: %v", err)
	}

}
