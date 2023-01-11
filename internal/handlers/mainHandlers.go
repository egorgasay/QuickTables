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
}

func NewHandler(db *repository.Storage) *Handler {
	return &Handler{service: service.New(db)}
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
		if _, ok := c.GetPostForm("delete"); ok {
			err := h.service.DB.DeleteDB(username, dbName)
			if err != nil {
				log.Println(err)
			}
			c.Redirect(http.StatusFound, "/")
			return
		}

		dbs := h.service.DB.GetAllDBs(username)
		connStr, driver, docker := h.service.DB.GetDBInfobyName(username, dbName)

		if docker == "true" && !userDB.IsDBCached(dbName, username) {
			err := runDBFromDocker(username, dbName)
			if err != nil {
				log.Println(err)
				return
			}

			err = userDB.RecordConnection(dbName, connStr, username, driver)
			if err != nil {
				log.Println(err)
				return
			}
			c.Redirect(http.StatusTemporaryRedirect, "/switch")
			return
		}

		if err := userDB.RecordConnection(dbName, connStr, username, driver); err != nil {
			c.HTML(http.StatusOK, "switch.html", gin.H{"DBs": dbs, "error": err.Error()})
			return
		}

		c.Redirect(http.StatusFound, "/")
		return
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
			rowForTable = append(rowForTable, el.String)
		}

		rowsForTable = append(rowsForTable, rowForTable)
	}

	t.AppendRows(rowsForTable)

	table := t.RenderHTML()

	c.Writer.Write([]byte(table))
}

func (h Handler) NotFoundHandler(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", gin.H{"page": "404"})
}

func (h Handler) HistoryHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)

	username, _ := user.(string)
	name := userDB.GetDbName(username)

	r, err := h.service.DB.GetQueries(username, name)
	if err != nil {
		log.Println(err)
		return
	}

	c.HTML(http.StatusOK, "newNav.html", gin.H{"Queries": r, "page": "history"})
}

func (h Handler) ProfileGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)

	username, _ := user.(string)

	us, err := h.service.DB.GetUserStats(username)
	if err != nil {
		c.HTML(http.StatusOK, "newNav.html", gin.H{"error": err.Error(), "page": "profile"})
		return
	}

	c.HTML(http.StatusOK, "newNav.html", gin.H{"username": username, "us": us, "page": "profile"})

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
