package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"log"
	"quicktables/config"
	"quicktables/internal/globals"
	"quicktables/internal/handlers"
	"quicktables/internal/middleware"
	"quicktables/internal/repository"
	"quicktables/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	cfg := config.New()
	storage, err := repository.New(cfg.DBConfig)

	if err != nil {
		log.Fatalf("Failed to initialize: %s", err.Error())
	}

	h := handlers.NewHandler(storage)

	r.LoadHTMLGlob("templates/html/*")
	r.Static("/static", "static")
	r.NoRoute(h.NotFoundHandler)

	r.Use(sessions.Sessions("session", cookie.NewStore(globals.Secret)))

	public := r.Group("/")
	routes.PublicRoutes(public, *h)

	private := r.Group("/")
	private.Use(middleware.AuthRequired)
	routes.PrivateRoutes(private, *h)

	r.Run("localhost:8080")
	//log.Fatal(http.ListenAndServe(":8080", r))
}
