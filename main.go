package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var store map[string]Secret = make(map[string]Secret)

type (
	Secret struct {
		Data      string `json:"data"`
		CreatedAt string `json:"createdAt"`
		Owner     string `json:"owner"`
	}

	CreatedResponse struct {
		Name string `json:"name"`
	}
)

func createSecret(c echo.Context) error {

	createdAt := strconv.FormatInt(time.Now().UnixNano(), 10)

	secret := new(Secret)
	secret.CreatedAt = createdAt

	if err := c.Bind(secret); err != nil {
		return err
	}

	var key string = namesgenerator.GetRandomName(0)

	if _, ok := store[key]; ok {
		key = namesgenerator.GetRandomName(1)
	}

	// TODO Add to database/in memory store?
	store[key] = *secret

	// TODO create random easy name ie grizzly-bear, wild-fish...
	response := &CreatedResponse{key}

	return c.JSON(http.StatusAccepted, response)
}

func fetchSecret(c echo.Context) error {
	key := c.Param("key")

	value, ok := store[key]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	// TODO check if value expired, delete and return no content not found

	delete(store, key)

	return c.JSON(http.StatusOK, value)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	server := echo.New()

	server.Use(middleware.Logger())

	server.POST("/secret", createSecret)
	server.GET("/secret/:key", fetchSecret)

	server.GET("/version", func(c echo.Context) error {
		return c.String(http.StatusOK, os.Getenv("ENV_VERSION"))
	})

	server.Logger.Fatal(server.Start(":" + port))
}
