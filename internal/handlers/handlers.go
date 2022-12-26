package handlers

import (
	"context"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"quicktables/internal/globals"
	"quicktables/internal/repository"
	"quicktables/internal/service"
	"quicktables/internal/userDB"
	"time"
)

type Handler struct {
	service *service.Service
	udb     map[string]*userDB.UserDB
}

func NewHandler(db *repository.Storage) *Handler {
	return &Handler{service: service.New(db), udb: userDB.New()}
}

type RegStruct struct {
	msg     string
	cssName string
}

var rs = RegStruct{"Sign Up", "reg_style.css"}

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
		c.HTML(http.StatusOK, "reg.html", gin.H{"msg": rs.msg,
			"css": rs.cssName})
		return
	} else if password != password2 {
		rs.msg = "Passwords don't match"
		rs.cssName = "bad_reg_style.css"
		c.HTML(http.StatusOK, "reg.html", gin.H{"msg": rs.msg,
			"css": rs.cssName})
		return
	}

	err := h.service.DB.CreateUser(username, password)
	if err != nil {
		rs.msg = "Username is already taken"
		rs.cssName = "bad_reg_style.css"
		c.HTML(http.StatusOK, "reg.html", gin.H{"msg": rs.msg,
			"css": rs.cssName})
		return
	}

	rs.msg = "Sign Up"
	rs.cssName = "reg_style.css"

	c.Redirect(http.StatusPermanentRedirect, "/login")
}

type LoginStruct struct {
	msg     string
	cssName string
}

var ls = RegStruct{"Sign In", "reg_style.css"}

func (h Handler) LoginHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	if user != nil {
		c.Redirect(http.StatusFound, "/logout")
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")
	status := h.service.DB.CheckPassword(username, password)

	if !status && password != "" {
		ls.msg = "Wrong password or username"
		ls.cssName = "bad_reg_style.css"

		c.HTML(http.StatusOK, "login.html", gin.H{"msg": ls.msg,
			"css": ls.cssName})

		return
	} else if status && password != "" {
		session.Set(globals.Userkey, username)
		if err := session.Save(); err != nil {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{"content": "Failed to save session"})
			return
		}

		c.Redirect(http.StatusFound, "/")

		return
	}

	ls.msg = "Sign In"
	ls.cssName = "reg_style.css"

	c.HTML(http.StatusOK, "login.html", gin.H{"msg": ls.msg,
		"css": ls.cssName})
}

type AddDBStruct struct {
	msg     string
	cssName string
}

var ads = AddDBStruct{msg: "Enter connection string", cssName: "reg_style.css"}

func (h Handler) AddDBPostHandler(c *gin.Context) {
	session := sessions.Default(c)

	user := session.Get(globals.Userkey)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	username := user.(string)

	connStr := c.PostForm("con_str")
	bdVendorName := c.PostForm("bdVendorName")

	if connStr == "" {
		ads.msg = "Enter connection string"
		ads.cssName = "reg_style.css"
		c.HTML(http.StatusOK, "addDB.html", gin.H{"msg": ads.msg,
			"css": ads.cssName})
		return
	}

	err := userDB.RecordConnection(connStr, username, bdVendorName)
	if err != nil {
		log.Println(err)
		ads.msg = "Error!"
		ads.cssName = "bad_reg_style.css"
		c.HTML(http.StatusOK, "addDB.html", gin.H{"msg": ads.msg,
			"css": ads.cssName})
		ads.msg = "Enter connection string"
		ads.cssName = "reg_style.css"
		return
	}

	err = h.service.DB.AddDB(connStr, username, bdVendorName)
	if err != nil {
		log.Println(err)
		ads.msg = "Server error!"
		ads.cssName = "bad_reg_style.css"
		c.HTML(http.StatusOK, "addDB.html", gin.H{"msg": ads.msg,
			"css": ads.cssName})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func (h Handler) AddDBGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	log.Println(user)

	c.HTML(http.StatusOK, "addDB.html", gin.H{"msg": ads.msg,
		"css": ads.cssName})
}

func (h Handler) MainGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	username := user.(string)
	if !checkUserDB(c, username, h.service.DB) {
		return
	}

	c.HTML(http.StatusOK, "newNav.html", gin.H{"page": "main"})
}

func checkUserDB(c *gin.Context, username string, db service.IService) bool {
	if !db.CheckDB(username) {
		c.Redirect(http.StatusFound, "/addDB")
		return false
	}

	if !userDB.CheckConn(username) {
		connStr, driver := db.GetDB(username)
		err := userDB.RecordConnection(connStr, username, driver)

		if err != nil {
			log.Println(err)
			ads.msg = "Error!"
			ads.cssName = "bad_reg_style.css"
			c.HTML(http.StatusOK, "addDB.html", gin.H{"error": ads.msg,
				"css": ads.cssName})
			return false
		}
	}

	return true
}

func (h Handler) MainPostHandler(c *gin.Context) {
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

	if !checkUserDB(c, username, h.service.DB) {
		return
	}

	//dbName := userDB.GetDbName(username)
	query := c.PostForm("query")
	ctx := context.Background()

	start := time.Now()

	rows, err := userDB.Query(ctx, username, query)
	if err != nil {
		c.HTML(http.StatusOK, "newNav.html", gin.H{"query": query, "error": err.Error(),
			"page": "main"})
		h.service.DB.SaveQuery(2, query, username, "", "0")
		return
	}

	cols, err := rows.Columns()
	if err != nil {
		c.HTML(http.StatusOK, "newNav.html", gin.H{"query": query, "error": err.Error(),
			"page": "main"})
		h.service.DB.SaveQuery(2, query, username, "", "0")
		return
	}

	err = h.service.DB.SaveQuery(1, query, username, "", time.Now().Sub(start).String())
	if err != nil {
		log.Println(err)
	}

	readCols := make([]interface{}, len(cols))
	writeCols := make([]string, len(cols))

	for i, _ := range writeCols {
		readCols[i] = &writeCols[i]
	}

	rowsArr := make([][]string, 0, 4)
	for rows.Next() {
		err := rows.Scan(readCols...)
		if err != nil {
			panic(err)
		}
		rowsArr = append(rowsArr, writeCols)
	}
	c.HTML(http.StatusOK, "newNav.html", gin.H{"query": query, "rows": rowsArr,
		"cols": cols, "page": "main"})
}

func (h Handler) LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)

	session.Delete(globals.Userkey)
	if err := session.Save(); err != nil {
		log.Println("Failed to save session:", err)
		return
	}

	c.Redirect(http.StatusFound, "/login")
}

func (h Handler) NotFoundHandler(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", gin.H{"page": "404"})
}

func (h Handler) HistoryHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)

	username, _ := user.(string)

	r, err := h.service.DB.GetQueries(username, "")
	if err != nil {
		log.Println(err)
		return
	}

	c.HTML(http.StatusOK, "newNav.html", gin.H{"Queries": r, "page": "history"})
}
