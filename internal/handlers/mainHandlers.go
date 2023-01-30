package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"quicktables/internal/globals"
	"quicktables/internal/repository"
	"quicktables/internal/service"
	"quicktables/internal/usecase"
	"quicktables/internal/userDB"
	"strings"
	"time"
)

type Handler struct {
	service *service.Service
	userDBs *userDB.UserDBs
}

func NewHandler(db *repository.Storage, userDBs *userDB.UserDBs) *Handler {
	return &Handler{service: service.New(db),
		userDBs: userDBs}
}

func (h Handler) MainGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	username, _ := user.(string)

	if !h.service.DB.CheckDB(username) {
		c.Redirect(http.StatusFound, "/addDB")
		return
	}

	dbs := h.service.DB.GetAllDBs(username)

	if (*h.userDBs)[username] == nil {
		c.HTML(http.StatusOK, "switch.html", gin.H{"DBs": dbs})
		return
	}

	if (*h.userDBs)[username].DBs == nil {
		c.HTML(http.StatusOK, "switch.html", gin.H{"DBs": dbs})
		return
	}

	activeDB := (*h.userDBs)[username].Active
	if activeDB == nil {
		c.Redirect(http.StatusFound, "/switch")
		return
	}

	vendor := activeDB.Driver

	c.HTML(http.StatusOK, "newNav.html", gin.H{"page": "main",
		"current": activeDB.Name, "vendor": vendor})
}

func (h Handler) MainPostHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	username, _ := user.(string)
	dbs := h.service.DB.GetAllDBs(username)

	udbs := (*h.userDBs)[username]
	if udbs == nil {
		udbs = &userDB.ConnStorage{}
		(*h.userDBs)[username] = udbs
	}

	if dbName := c.PostForm("dbName"); dbName != "" {
		_, remove := c.GetPostForm("delete")
		if remove {
			err := usecase.DeleteUserDB(h.service.DB, username, dbName)
			if err != nil {
				c.HTML(http.StatusOK, "newNav.html", gin.H{"page": "main",
					"DBs": dbs, "error": err})
				return
			}

			c.Redirect(http.StatusFound, "/")
			return
		}

		err := usecase.HandleUserDB(h.service.DB, username, dbName, udbs)
		if err != nil {
			c.HTML(http.StatusOK, "newNav.html", gin.H{"page": "main",
				"DBs": dbs, "error": err})
			return
		}
	}

	currentDB, vendorDB := udbs.Active.Name, udbs.Active.Driver

	query := c.PostForm("query")

	if strings.Trim(query, " ") == "" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	start := time.Now()
	qh, err := usecase.HandleQuery(udbs, query)
	if err != nil {
		qh.Status = 2
		h.service.DB.SaveQuery(qh.Status, query, username, currentDB, "0")
		c.HTML(http.StatusOK, "newNav.html", gin.H{"query": query,
			"page": "main", "current": currentDB, "vendor": vendorDB, "error": err})
		return
	}

	qh.Status = 1
	h.service.DB.SaveQuery(qh.Status, query, username, currentDB, time.Now().Sub(start).String())

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

	udbs := (*h.userDBs)[username]

	name := udbs.Active.Name

	r, err := h.service.DB.GetQueries(username, name)
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

	us, err := h.service.DB.GetUserStats(username)
	if err != nil {
		c.HTML(http.StatusOK, "newNav.html", gin.H{"uname": username,
			"error": err.Error(), "page": "profile"})
		return
	}

	c.HTML(http.StatusOK, "newNav.html", gin.H{"uname": username,
		"us": us, "page": "profile"})

}

func (h Handler) ProfilePostHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)

	username, _ := user.(string)
	nick, ok := c.GetPostForm("new-nick")
	if ok && nick != "" {
		err := h.service.DB.ChangeNick(username, nick)
		if err != nil {
			c.HTML(http.StatusOK, "newNav.html", gin.H{"error": err.Error(), "page": "profile"})
			return
		}
	}

	oldPassword, okOldPassword := c.GetPostForm("old-password")
	newPassword, okNewPassword := c.GetPostForm("new-password")
	if okOldPassword && okNewPassword && newPassword != "" {
		err := h.service.DB.ChangePassword(username, oldPassword, newPassword)
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
