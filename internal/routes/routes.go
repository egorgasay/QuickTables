package routes

import (
	"github.com/gin-gonic/gin"
	"quicktables/internal/handlers"
)

func PublicRoutes(r *gin.RouterGroup, h handlers.Handler) {
	r.Any("/reg", h.RegisterHandler)
	r.Any("/login", h.LoginHandler)
}

func PrivateRoutes(r *gin.RouterGroup, h handlers.Handler) {
	r.GET("/addDB", h.AddDBGetHandler)
	r.POST("/addDB", h.AddDBPostHandler)
	r.GET("/", h.MainGetHandler)
	r.POST("/", h.MainPostHandler)
	r.GET("/logout", h.LogoutHandler)
	r.GET("/history", h.HistoryHandler)
	r.GET("/switch", h.SwitchGetHandler)
	r.POST("/switch", h.SwitchPostHandler)
	r.GET("/list", h.ListHandler)
	r.GET("/list/:name", h.ListHandler)
	r.GET("/profile", h.ProfileGetHandler)
	r.POST("/profile", h.ProfilePostHandler)
	r.GET("/createDB", h.CreateDBGetHandler)
	r.POST("/createDB", h.CreateDBPostHandler)
	//r.POST("/api", h.ApiHandler)
}
