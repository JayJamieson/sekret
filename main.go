package main

import (
	"context"
	"embed"
	"errors"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/JayJamieson/sekret/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

//go:embed ui/dist
var embded embed.FS

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func getFileSystem(useOS bool) fs.FS {
	if useOS {
		log.Print("using live mode")
		return os.DirFS("ui")
	}

	log.Print("using embed mode")
	fsys, err := fs.Sub(embded, "ui/dist")
	if err != nil {
		panic(err)
	}

	return fsys
}

func main() {
	port := os.Getenv("PORT")
	dbURI := os.Getenv("DB_URI")
	useOS := len(os.Args) > 1 && os.Args[1] == "live"

	if port == "" {
		port = "8080"
		log.Print("PORT environment variable must be set, defaulted to 8080")
	}

	server := echo.New()

	uiFS := getFileSystem(useOS)
	assetHandler := http.FileServer(http.FS(uiFS))

	server.Renderer = &Template{
		templates: template.Must(template.ParseFS(uiFS, "*.html")),
	}

	handlers, err := handlers.New(dbURI)

	if err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}

	server.Logger.SetLevel(log.INFO)
	server.Use(middleware.Logger())
	server.Use(middleware.Secure())

	server.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))

	api := server.Group("/api")
	api.POST("/secret", handlers.CreateSecret)
	api.GET("/secret/:key", handlers.GetSecret)

	server.GET("/", handlers.ViewIndex)
	server.POST("/secret", handlers.CreateSecret)
	server.GET("/secret/:key", handlers.ViewSecret)
	server.POST("/secret/:key", handlers.ViewSecret)
	server.GET("/private/:key", handlers.ViewPrivate)

	server.GET("/health", handlers.Health)
	server.GET("/version", func(c echo.Context) error {
		return c.String(http.StatusOK, os.Getenv("ENV_VERSION"))
	})

	go func() {
		if err := server.Start(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			server.Logger.Error(err)
			server.Logger.Fatal("Shutting down")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
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
