package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"quicktables/internal/globals"
	"quicktables/internal/usecase"
	"strings"
	"time"
)

type Handler struct {
	logic usecase.UseCase
}

func NewHandler(logic usecase.UseCase) *Handler {
	return &Handler{logic: logic}
}

func (h Handler) MainGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	username, _ := user.(string)

	if !h.logic.Service.DB.CheckDB(username) {
		c.Redirect(http.StatusFound, "/addDB")
		return
	}

	dbs := h.logic.Service.DB.GetAllDBs(username)

	vendor, name, err := h.logic.GetVendorAndName(username)
	if err != nil {
		c.HTML(http.StatusOK, "switch.html", gin.H{"DBs": dbs})
		return
	}

	c.HTML(http.StatusOK, "newNav.html", gin.H{"page": "main",
		"current": name, "vendor": vendor})
}

func (h Handler) MainPostHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	username, _ := user.(string)
	dbs := h.logic.Service.DB.GetAllDBs(username)

	if dbName := c.PostForm("dbName"); dbName != "" {
		_, remove := c.GetPostForm("delete")
		if remove {
			err := h.logic.DeleteUserDB(username, dbName)
			if err != nil {
				c.HTML(http.StatusOK, "switch.html", gin.H{"DBs": dbs, "error": err})
				return
			}

			c.Redirect(http.StatusFound, "/")
			return
		}

		var err error
		err = h.logic.HandleUserDB(username, dbName)
		if err != nil {
			c.HTML(http.StatusOK, "switch.html", gin.H{
				"DBs": dbs, "error": err})
			return
		}
	}

	query := c.PostForm("query")

	vendorDB, currentDB, err := h.logic.GetVendorAndName(username)
	if err != nil {
		c.HTML(http.StatusOK, "newNav.html", gin.H{"query": query,
			"page": "main", "current": currentDB, "vendor": vendorDB, "error": err})
		return
	}

	if strings.Trim(query, " ") == "" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	start := time.Now()
	qh, err := h.logic.HandleQuery(query, username)
	if err != nil {
		qh.Status = 2
		saveErr := h.logic.Service.DB.SaveQuery(qh.Status, query, username, currentDB, "0")
		if err != nil {
			log.Printf("%s \n", saveErr)
		}
		c.HTML(http.StatusOK, "newNav.html", gin.H{"query": query,
			"page": "main", "current": currentDB, "vendor": vendorDB, "error": err})
		return
	}

	qh.Status = 1
	err = h.logic.Service.DB.SaveQuery(qh.Status, query, username, currentDB, time.Now().Sub(start).String())
	if err != nil {
		log.Printf("%s \n", err.Error())
	}

	if !qh.IsSelect {
		c.HTML(http.StatusOK, "newNav.html", gin.H{"query": query,
			"msg": "Completed Successfully", "page": "main", "current": currentDB,
			"vendor": vendorDB})
		return
	}

	c.HTML(http.StatusOK, "newNav.html", gin.H{"query": query, "rows": qh.Table.Rows,
		"cols": qh.Table.Cols, "page": "main", "current": currentDB, "vendor": vendorDB})
}

func (h Handler) NotFoundHandler(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", gin.H{"page": "404"})
}

func (h Handler) HistoryHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	username, ok := user.(string)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	r, err := h.logic.GetHistory(username)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusOK, "newNav.html", gin.H{"error": err, "page": "history"})
		return
	}

	var notify string
	if len(r) == 0 {
		notify = "You don't have any query with this db. Let's create one!"
	}

	c.HTML(http.StatusOK, "newNav.html", gin.H{"Queries": r, "page": "history",
		"notify": notify})
}

func (h Handler) ProfileGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)

	username, _ := user.(string)

	us, err := h.logic.GetProfile(username)
	if err != nil {
		c.HTML(http.StatusOK, "newNav.html", gin.H{"uname": username,
			"error": err.Error(), "page": "profile"})
		return
	}

	c.HTML(http.StatusOK, "newNav.html", gin.H{"uname": username,
		"us": us, "page": "profile"})

}

// ProfilePostHandler Handler of user profile
func (h Handler) ProfilePostHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)

	username, _ := user.(string)
	nick, ok := c.GetPostForm("new-nick")
	if ok && nick != "" {
		err := h.logic.Service.DB.ChangeNick(username, nick)
		if err != nil {
			c.HTML(http.StatusOK, "newNav.html", gin.H{"error": err, "page": "profile"})
			return
		}
	}

	oldPassword, okOldPassword := c.GetPostForm("old-password")
	newPassword, okNewPassword := c.GetPostForm("new-password")
	if okOldPassword && okNewPassword && newPassword != "" {
		err := h.logic.Service.DB.ChangePassword(username, oldPassword, newPassword)
		if err != nil {
			c.HTML(http.StatusOK, "newNav.html", gin.H{"error": err.Error(), "page": "profile"})
			return
		}
	}

	c.Redirect(http.StatusFound, "/profile")
}

//func (h Handler) ApiHandler(c *gin.Context) {
//	session := sessions.Default(c)
//	user := session.Get(globals.Userkey)
//	username, _ := user.(string)
//
//	data := userDB.GetUserDataFromDB(username)
//	count := len(data)
//	start, _ := c.Get("start")
//	end, _ := c.Get("length")
//	draw, _ := c.Get("draw")
//
//	fmt.Println(count, start, end, draw)
//}
