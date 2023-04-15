package main

import (
	"context"
	"github.com/egorgasay/dockerdb/v2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"log"
	"quicktables/config"
	"quicktables/internal/globals"
	"quicktables/internal/handlers"
	"quicktables/internal/middleware"
	"quicktables/internal/repository"
	"quicktables/internal/routes"
	"quicktables/internal/service"
	"quicktables/internal/usecase"
	"quicktables/internal/userDB"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	cfg := config.New()

	storage, err := repository.New(cfg.DBConfig)
	if err != nil {
		log.Fatalf("Failed to initialize: %s", err.Error())
	}

	logic := usecase.New(service.New(storage), userDB.New())

	ctx := context.TODO()
	for _, vendor := range globals.DownloadableVendors {
		err = dockerdb.Pull(ctx, vendor)
		if err != nil {
			log.Fatalf("Failed to download images: %s", err.Error())
		}
	}

	h := handlers.NewHandler(logic)

	r.LoadHTMLGlob("templates/html/*")
	r.Static("/static", "static")
	r.NoRoute(h.NotFoundHandler)

	r.Use(sessions.Sessions("session", cookie.NewStore(globals.Secret)))

	public := r.Group("/")
	routes.PublicRoutes(public, *h)

	private := r.Group("/")
	private.Use(middleware.AuthRequired)
	routes.PrivateRoutes(private, *h)

	r.Run(":8000")
	//log.Fatal(http.ListenAndServe(":8080", r))
}
