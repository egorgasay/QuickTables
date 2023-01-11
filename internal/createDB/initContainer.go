package createdb

import (
	"errors"
	"fmt"
	"os"
	"quicktables/internal/userDB"
	"text/template"
)

var composePostgres = `version: "3.9"
services:
  postgres:
    image: postgres:13.3
    environment:
      POSTGRES_DB: "{{.DB.Name}}"
      POSTGRES_USER: "{{.DB.User}}"
      POSTGRES_PASSWORD: "{{.DB.Password}}"
      PGDATA: "/var/lib/postgresql/data/pgdata"
    volumes:
      - ./users/{{.Username}}/2. Init Database:/docker-entrypoint-initdb.d
      - ./{{.DB.Name}}:/var/lib/postgresql/data
    ports:
      - "{{.Port}}:5432"
`

var composeMySQL = `version: "3.9"
services:
  db:
    image: mysql:5.7
    restart: always
    environment:
      MYSQL_DATABASE: '{{.DB.Name}}'
      MYSQL_USER: '{{.DB.User}}'
      MYSQL_PASSWORD: '{{.DB.Password}}'
      # Password for root access
      MYSQL_ROOT_PASSWORD: '{{.DB.Password}}'
    ports:
        - '{{.Port}}:3306'
    expose:
    # Opens port 3306 on the container
        - '3306'
    # Where our data will be persisted
    volumes:
        - .data:/var/lib/mysql
`

// TODO: IMPLEMENT DOWNLOAD ALL IMAGES

func InitContainer(conf *userDB.CustomDB) error {
	if conf == nil {
		return errors.New("conf must be not nil")
	}

	err := os.MkdirAll("users/"+conf.Username+"/"+conf.DB.Name, 0644)
	if err != nil {
		return err
	}

	var compose string

	if conf.Vendor == "postgres" {
		compose = composePostgres
	} else if conf.Vendor == "mysql" {
		compose = composeMySQL
	}

	tmpl := template.Must(template.New("compose").Parse(compose))
	path := fmt.Sprintf("users/%v/%v/%s", conf.Username, conf.DB.Name, "compose.yaml")

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, conf)
	if err != nil {
		return err
	}

	file.Close()

	return RunContainer(path)
}
