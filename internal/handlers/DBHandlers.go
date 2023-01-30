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
	"os"
	"quicktables/internal/dockerdb"
	"quicktables/internal/globals"
	"quicktables/internal/usecase"
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
	udbs := (*h.userDBs)[username]

	dbName := c.PostForm("dbName")
	connStr := c.PostForm("con_str")
	bdVendorName := c.PostForm("bdVendorName")

	if connStr == "" {
		c.Redirect(http.StatusFound, "/addDB")
		return
	}

	err := udbs.RecordConnection(dbName, connStr, bdVendorName)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusOK, "addDB.html", gin.H{"error": err,
			"vendors": globals.AvailableVendors})
		return
	}

	err = h.service.DB.AddDB(dbName, connStr, username, bdVendorName, "")
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
	udbs := (*h.userDBs)[username]

	dbs := h.service.DB.GetAllDBs(username)
	dbName := c.PostForm("dbName")

	_, remove := c.GetPostForm("delete")
	if remove {
		err := usecase.DeleteUserDB(h.service.DB, username, dbName)
		if err != nil {
			c.HTML(http.StatusOK, "switch.html", gin.H{"DBs": dbs, "error": err.Error()})
			return
		}

		c.Redirect(http.StatusFound, "/switch")
		return
	}

	err := usecase.HandleUserDB(h.service.DB, username, dbName, udbs)
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

	udbs := (*h.userDBs)[username]
	ctx := context.Background()

	list, err := usecase.GetListOfUserTables(ctx, udbs, username)
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

	userTable, err := usecase.GetUserTable(ctx, udbs, username, dbName)
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
	//c.HTML(http.StatusOK, "loadDB.html", gin.H{})
	udbs := (*h.userDBs)[username]

	pswd, err := password.Generate(17, 5, 0, false, false)
	if err != nil {
		log.Fatal(err)
	}

	var port string

	port, err = GetFreePort()
	if err != nil {
		log.Println("can't get free port:", err)
	}

	if h.service.DB.BindPort(port) != nil {
		log.Println(err)
	}

	if bdVendorName == "sqlite3" {
		path := fmt.Sprintf("users/%s/", username)
		err = os.MkdirAll(path, 777)
		if err != nil {
			return
		}

		err = h.service.DB.AddDB(dbName, path+dbName, username, bdVendorName, "")
		if err != nil {
			log.Println(err)
		}

		err = udbs.RecordConnection(dbName, path+dbName, "sqlite3")
		if err != nil {
			log.Println(err)
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

	connStr := udbs.StrConnBuilder(conf)

	err = udbs.CheckConnDocker(connStr, conf.Vendor)
	if err != nil {
		log.Println(err, "check docker")
		c.HTML(http.StatusOK, "createDB.html", gin.H{"error": err.Error(),
			"vendors": globals.CreatebleVendors})
		return
	}

	err = h.service.DB.AddDB(dbName, connStr, username, bdVendorName, ddb.ID)
	if err != nil {
		log.Println(err, "AddDB")
		c.HTML(http.StatusOK, "createDB.html", gin.H{"error": err.Error(),
			"vendors": globals.CreatebleVendors})
		return
	}

	err = udbs.RecordConnection(dbName, connStr, bdVendorName)
	if err != nil {
		log.Println(err, "RecordConnection")
		c.HTML(http.StatusOK, "createDB.html", gin.H{"error": err.Error(),
			"vendors": globals.CreatebleVendors})
		return
	}

	//err = userDB.AddDockerCli(cli, conf)
	//if err != nil {
	//	log.Println(err, "AddDockerCli")
	//	c.HTML(http.StatusOK, "createDB.html", gin.H{"error": err.Error(),
	//		"vendors": globals.CreatebleVendors})
	//	return
	//}

	c.Redirect(http.StatusFound, "/")
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
