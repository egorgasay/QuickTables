package handlers

import (
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-password/password"
	"log"
	"net"
	"net/http"
	createdb "quicktables/internal/createDB"
	"quicktables/internal/globals"
	"quicktables/internal/userDB"
	"strconv"
)

func (h Handler) AddDBGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	log.Println(user)

	c.HTML(http.StatusOK, "addDB.html", gin.H{"vendors": globals.AvailableVendors})
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
		c.Redirect(http.StatusFound, "/addDB")
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

	if _, ok := c.GetPostForm("delete"); ok {
		err := h.service.DB.DeleteDB(username, dbName)
		if err != nil {
			log.Println(err)
		}
		dbs := h.service.DB.GetAllDBs(username)
		c.HTML(http.StatusOK, "newNav.html", gin.H{"DBs": dbs, "page": "switch"})
		return
	}

	connStr, driver := h.service.DB.GetDBbyName(username, dbName)
	err := userDB.SetMainDbByName(dbName, username, connStr, driver)

	if err != nil {
		dbs := h.service.DB.GetAllDBs(username)
		c.HTML(http.StatusOK, "newNav.html", gin.H{"DBs": dbs, "error": err.Error(), "page": "switch"})
		return
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
		query := fmt.Sprintf(`SELECT * FROM "%s"`, name)

		rows, err := userDB.Query(ctx, username, query)
		if err != nil {
			query = fmt.Sprintf(`SELECT * FROM %s`, name)
			rows, err = userDB.Query(ctx, username, query)

			if err != nil {
				log.Println(err)
				c.HTML(http.StatusOK, "newNav.html", gin.H{"error": err.Error(),
					"page": "list"})
			}
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

func (h Handler) CreateDBGetHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "createDB.html", gin.H{"vendors": globals.CreatebleVendors})
}

func (h Handler) CreateDBPostHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	username, _ := user.(string)
	dbName := c.PostForm("dbName")
	bdVendorName := c.PostForm("bdVendorName")
	c.HTML(http.StatusOK, "loadDB.html", gin.H{})

	go func(c *gin.Context) {
		pswd, err := password.Generate(17, 5, 0, false, false)
		if err != nil {
			log.Fatal(err)
		}

		var port string

		for {
			port, err = GetFreePort()
			if err != nil {
				log.Fatal(err)
			}

			if h.service.DB.BindPort(port) == nil {
				break
			}
		}

		if bdVendorName == "sqlite3" {
			c.Redirect(http.StatusTemporaryRedirect, "/addDB")
			return
		}

		v := &userDB.CustomDB{
			DB: userDB.DB{
				Name:     dbName,
				User:     "admin",
				Password: pswd,
			},
			Username: username,
			Port:     port,
		}

		err = createdb.InitContainer(v)
		if err != nil {
			log.Println(err)
			return
		}

		switch bdVendorName {
		case "postgres":
			connStr, err := userDB.RecordConnPostgres(v)
			if err != nil {
				log.Println(err)
			}

			err = h.service.DB.AddDB(dbName, connStr, username, bdVendorName)
			if err != nil {
				log.Println(err)
			}
		}
	}(c)
	return
}

func GetFreePort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "0", err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "0", err
	}

	defer l.Close()
	port := strconv.Itoa(l.Addr().(*net.TCPAddr).Port)

	return port, nil
}
