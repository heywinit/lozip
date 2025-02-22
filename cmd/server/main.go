package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/heywinit/lozip/internal/server"
)

type App struct {
	config *server.Config
	fiber  *fiber.App
}

func NewApp() (*App, error) {
	// Load configuration
	config, err := server.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Lozip API",
	})

	return &App{
		config: config,
		fiber:  app,
	}, nil
}

func (a *App) Start() error {
	// Setup routes
	server.SetupRoutes(a.fiber, a.config)

	// Start server
	addr := fmt.Sprintf(":%s", a.config.Port)
	log.Printf("Server starting on %s", addr)
	return a.fiber.Listen(addr)
}

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}


