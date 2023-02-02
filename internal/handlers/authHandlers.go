package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"quicktables/internal/globals"
)

func (h Handler) RegisterHandler(c *gin.Context) {
	session := sessions.Default(c)

	user := session.Get(globals.Userkey)
	if user != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")
	password2 := c.PostForm("password2")

	if password == "" && password2 == "" {
		c.HTML(http.StatusOK, "reg.html", gin.H{})
		return
	} else if password != password2 {
		c.HTML(http.StatusOK, "reg.html", gin.H{"err": "Passwords don't match"})
		return
	}

	err := h.logic.Service.DB.CreateUser(username, password)
	if err != nil {
		c.HTML(http.StatusOK, "reg.html", gin.H{"err": "Username is already taken"})
		return
	}

	c.Redirect(http.StatusPermanentRedirect, "/login")
}

func (h Handler) LoginHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	if user != nil {
		c.Redirect(http.StatusFound, "/logout")
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")
	status := h.logic.Service.DB.CheckPassword(username, password)

	if !status && password != "" {
		c.HTML(http.StatusOK, "login.html", gin.H{"err": "Wrong password or username"})
		return
	} else if status && password != "" {
		session.Set(globals.Userkey, username)
		if err := session.Save(); err != nil {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{"err": "Failed to save session"})
			return
		}

		c.Redirect(http.StatusFound, "/")

		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func (h Handler) LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)

	session.Delete(globals.Userkey)
	if err := session.Save(); err != nil {
		log.Println("Failed to save session:", err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/login")
}
