package handlers

import (
	"context"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-password/password"
	"log"
	"net/http"
	"quicktables/internal/dockerdb"
	"quicktables/internal/globals"
	"quicktables/internal/userDB"
	"quicktables/pkg"
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
	vendorName := c.PostForm("bdVendorName")

	if connStr == "" {
		c.Redirect(http.StatusFound, "/addDB")
		return
	}

	err := h.logic.AddUserDB(username, dbName, connStr, vendorName)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusOK, "addDB.html", gin.H{"error": err,
			"vendors": globals.AvailableVendors})
		return
	}

	c.Redirect(http.StatusFound, "/switch")
}

func (h Handler) SwitchGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)

	username, _ := user.(string)

	dbs := h.logic.GetAllDBs(username)

	c.HTML(http.StatusOK, "newNav.html", gin.H{"DBs": dbs, "page": "switch"})
}

func (h Handler) SwitchPostHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	username, _ := user.(string)

	dbs := h.logic.GetAllDBs(username)
	dbName := c.PostForm("dbName")

	_, remove := c.GetPostForm("delete")
	if remove {
		err := h.logic.DeleteUserDB(username, dbName)
		if err != nil {
			c.HTML(http.StatusOK, "switch.html", gin.H{"DBs": dbs, "error": err.Error()})
			return
		}

		c.Redirect(http.StatusFound, "/switch")
		return
	}

	err := h.logic.HandleUserDB(username, dbName)
	if err != nil {
		c.HTML(http.StatusOK, "switch.html", gin.H{"DBs": dbs, "error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func (h Handler) ListHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	username, _ := user.(string)

	ctx := context.Background()

	list, err := h.logic.GetListOfUserTables(ctx, username)
	if err != nil {
		c.HTML(http.StatusOK, "newNav.html", gin.H{"Tables": list, "page": "list",
			"error": err})
		return
	}
	dbName := c.Param("name")

	if dbName == "" {
		c.HTML(http.StatusOK, "newNav.html", gin.H{"Tables": list, "page": "list"})
		return
	}

	userTable, err := h.logic.GetUserTable(ctx, username, dbName)
	if err != nil {
		c.HTML(http.StatusOK, "newNav.html", gin.H{"Tables": list, "page": "list",
			"error": err})
		return
	}

	if userTable.HTMLTable != "" {
		c.Writer.Write([]byte(userTable.HTMLTable))
		return
	}

	c.HTML(http.StatusOK, "newNav.html", gin.H{"Tables": list, "page": "list",
		"rows": userTable.Rows, "cols": userTable.Cols})
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

	pswd, err := password.Generate(17, 5, 0, false, false)
	if err != nil {
		log.Fatal(err)
	}

	var port string

	port, err = pkg.GetFreePort()
	if err != nil {
		log.Println("can't get free port:", err)
	}

	if h.logic.BindPort(port) != nil {
		log.Println(err)
	}

	if bdVendorName == "sqlite3" {
		err = h.logic.CreateSqlite(username, dbName)
		if err != nil {
			c.HTML(http.StatusOK, "createDB.html", gin.H{"error": err.Error(),
				"vendors": globals.CreatebleVendors})
		}

		c.Redirect(http.StatusFound, "/")
		return
	}

	conf := &userDB.CustomDB{
		DB: userDB.DB{
			Name:     dbName,
			User:     "admin",
			Password: pswd,
		},
		Username: username,
		Port:     port,
		Vendor:   bdVendorName,
	}

	ddb, err := dockerdb.New(conf)
	if err != nil {
		log.Println(err, "init docker")
		c.HTML(http.StatusOK, "createDB.html", gin.H{"error": err.Error(),
			"vendors": globals.CreatebleVendors})
		return
	}

	ctx := context.TODO()

	err = ddb.Init(ctx)
	if err != nil {
		log.Println(err, "init docker")
		c.HTML(http.StatusOK, "createDB.html", gin.H{"error": err.Error(),
			"vendors": globals.CreatebleVendors})
		return
	}

	ctx = context.TODO()

	err = ddb.Run(ctx)
	if err != nil {
		log.Println(err, "run docker")
		c.HTML(http.StatusOK, "createDB.html", gin.H{"error": err.Error(),
			"vendors": globals.CreatebleVendors})
		return
	}

	err = h.logic.HandleDocker(username, ddb, conf)
	if err != nil {
		log.Println(err, "check docker")
		c.HTML(http.StatusOK, "createDB.html", gin.H{"error": err.Error(),
			"vendors": globals.CreatebleVendors})
		return
	}

	c.Redirect(http.StatusFound, "/switch")
	return
}
