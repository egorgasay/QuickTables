package handlers

import (
	"context"
	"database/sql"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jedib0t/go-pretty/table"
	"log"
	"net/http"
	"quicktables/internal/globals"
	"quicktables/internal/repository"
	"quicktables/internal/service"
	"quicktables/internal/userDB"
	"strings"
	"time"
)

type Handler struct {
	service *service.Service
	//udb     map[string]*userDB.UserDB
}

func NewHandler(db *repository.Storage) *Handler {
	return &Handler{service: service.New(db)}
}

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

	err := h.service.DB.CreateUser(username, password)
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
	status := h.service.DB.CheckPassword(username, password)

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

func (h Handler) AddDBPostHandler(c *gin.Context) {
	session := sessions.Default(c)

	user := session.Get(globals.Userkey)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	username := user.(string)

	dbName := c.PostForm("dbName")
	connStr := c.PostForm("con_str")
	bdVendorName := c.PostForm("bdVendorName")

	if connStr == "" {
		c.HTML(http.StatusOK, "addDB.html", gin.H{})
		return
	}

	err := userDB.RecordConnection(dbName, connStr, username, bdVendorName)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusOK, "addDB.html", gin.H{"err": "Error!",
			"vendors": globals.AvailableVendors})
		return
	}

	err = h.service.DB.AddDB(dbName, connStr, username, bdVendorName)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusOK, "addDB.html", gin.H{"msg": "Server error!",
			"vendors": globals.AvailableVendors})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func (h Handler) AddDBGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	log.Println(user)

	c.HTML(http.StatusOK, "addDB.html", gin.H{"vendors": globals.AvailableVendors})
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

	currentDB, vendorDB := userDB.GetDbNameAndVendor(username)

	c.HTML(http.StatusOK, "newNav.html", gin.H{"page": "main",
		"current": currentDB, "vendor": vendorDB})
}

func checkUserDB(c *gin.Context, username string, db service.IService) bool {
	if !db.CheckDB(username) {
		c.Redirect(http.StatusFound, "/addDB")
		return false
	}

	if !userDB.CheckConn(username) {
		dbs := db.GetAllDBs(username)
		c.HTML(http.StatusOK, "switch.html", gin.H{"DBs": dbs})
		return false
		//dbName, connStr, driver := db.GetDB(username)
		//err := userDB.RecordConnection(dbName, connStr, username, driver)
		//
		//if err != nil {
		//	log.Println(err)
		//	dbs := db.GetAllDBs(username)
		//	c.HTML(http.StatusOK, "switch.html", gin.H{"DBs": dbs, "error": err.Error()})
		//	return false
		//}
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

	username, _ := user.(string)

	if dbName := c.PostForm("dbName"); dbName != "" {
		dbs := h.service.DB.GetAllDBs(username)
		connStr, driver := h.service.DB.GetDBbyName(username, dbName)

		if err := userDB.RecordConnection(dbName, connStr, username, driver); err != nil {
			c.HTML(http.StatusOK, "switch.html", gin.H{"DBs": dbs, "error": err.Error()})
			return
		}

		c.Redirect(http.StatusFound, "/")
	}

	query := c.PostForm("query")

	if strings.Trim(query, " ") == "" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	ctx := context.Background()
	currentDB, vendorDB := userDB.GetDbNameAndVendor(username)
	start := time.Now()

	rows, err := userDB.Query(ctx, username, query)
	if err != nil {
		c.HTML(http.StatusOK, "newNav.html", gin.H{"query": query, "error": err.Error(),
			"page": "main", "current": currentDB, "vendor": vendorDB})
		h.service.DB.SaveQuery(2, query, username, currentDB, "0")
		return
	}

	cols, err := rows.Columns()
	if err != nil {
		c.HTML(http.StatusOK, "newNav.html", gin.H{"query": query, "error": err.Error(),
			"page": "main", "current": currentDB, "vendor": vendorDB})
		h.service.DB.SaveQuery(2, query, username, currentDB, "0")
		return
	}

	err = h.service.DB.SaveQuery(1, query, username, currentDB, time.Now().Sub(start).String())
	if err != nil {
		log.Println(err)
	}

	rowsArr := doTableFromData(cols, rows)
	if len(rowsArr) > 1000 {
		doLargeTable(c, cols, rowsArr)
		return
	}

	c.HTML(http.StatusOK, "newNav.html", gin.H{"query": query, "rows": rowsArr,
		"cols": cols, "page": "main", "current": currentDB, "vendor": vendorDB})
}

func doTableFromData(cols []string, rows *sql.Rows) [][]sql.NullString {
	readCols := make([]interface{}, len(cols))
	writeCols := make([]sql.NullString, len(cols))

	rowsArr := make([][]sql.NullString, 0, 1000)
	for i := 0; rows.Next(); i++ {

		for i, _ := range writeCols {
			readCols[i] = &writeCols[i]
		}

		err := rows.Scan(readCols...)
		if err != nil {
			panic(err)
		}
		rowsArr = append(rowsArr, make([]sql.NullString, len(cols)))
		copy(rowsArr[i], writeCols)
	}

	return rowsArr
}

func doLargeTable(c *gin.Context, cols []string, rowsArr [][]sql.NullString) {
	t := table.NewWriter()

	colsForTable := make(table.Row, 0, 10)
	for _, el := range cols {
		colsForTable = append(colsForTable, el)
	}

	t.AppendHeader(colsForTable)

	rowsForTable := make([]table.Row, 0, 2000)
	for _, el := range rowsArr {
		rowForTable := make(table.Row, 0, 10)

		for _, el := range el {
			rowForTable = append(rowForTable, el)
		}

		rowsForTable = append(rowsForTable, rowForTable)
	}

	t.AppendRows(rowsForTable)

	table := t.RenderHTML()

	c.Writer.Write([]byte(table))
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
	name, _ := userDB.GetDbNameAndVendor(username)

	r, err := h.service.DB.GetQueries(username, name)
	if err != nil {
		log.Println(err)
		return
	}

	c.HTML(http.StatusOK, "newNav.html", gin.H{"Queries": r, "page": "history"})
}

func (h Handler) SwitchGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)

	username, _ := user.(string)

	dbs := h.service.DB.GetAllDBs(username)

	c.HTML(http.StatusOK, "newNav.html", gin.H{"DBs": dbs, "page": "switch"})
}

func (h Handler) SwitchPostHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	username, _ := user.(string)

	dbName := c.PostForm("dbName")
	connStr, driver := h.service.DB.GetDBbyName(username, dbName)
	err := userDB.SetMainDbByName(dbName, username, connStr, driver)

	if err != nil {
		dbs := h.service.DB.GetAllDBs(username)
		c.HTML(http.StatusOK, "newNav.html", gin.H{"DBs": dbs, "error": err.Error(), "page": "switch"})
	}

	c.Redirect(http.StatusFound, "/")
}

func (h Handler) ListHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	username, _ := user.(string)

	ctx := context.Background()
	list, err := userDB.GetAllTables(ctx, username)

	if err != nil {
		c.HTML(http.StatusOK, "newNav.html", gin.H{"error": err.Error(),
			"page": "list"})
		return
	}

	if name := c.Param("name"); name != "" {
		ctx := context.Background()

		rows, err := userDB.Query(ctx, username, `SELECT * FROM "`+name+`"`)
		if err != nil {
			log.Println(err)
		}

		cols, _ := rows.Columns()

		rowsArr := doTableFromData(cols, rows)

		if len(rowsArr) > 1000 {
			doLargeTable(c, cols, rowsArr)
			return
		}

		c.HTML(http.StatusOK, "newNav.html", gin.H{"Tables": list,
			"page": "list", "rows": rowsArr, "cols": cols})
		return
	}

	c.HTML(http.StatusOK, "newNav.html", gin.H{"Tables": list, "page": "list"})
}
