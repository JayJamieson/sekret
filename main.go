package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/JayJamieson/sekret/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	server := echo.New()

	sekretServer := handlers.NewSekret()

	server.Logger.SetLevel(log.INFO)

	server.Use(middleware.Logger())

	server.POST("/secret", sekretServer.CreateSecret)
	server.GET("/secret/:key", sekretServer.FetchSecret)

	server.GET("/version", func(c echo.Context) error {
		return c.String(http.StatusOK, os.Getenv("ENV_VERSION"))
	})

	go func() {
		if err := server.Start(":" + port); err != nil && err != http.ErrServerClosed {
			server.Logger.Fatal("Shutting down")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		server.Logger.Fatal(err)
	}
}
