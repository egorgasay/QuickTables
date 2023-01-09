package createdb

import (
	"errors"
	"fmt"
	"os"
	"quicktables/internal/userDB"
	"text/template"
)

var compose = `version: "3.9"
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

func InitContainer(conf *userDB.CustomDB) error {
	if conf == nil {
		return errors.New("conf must be not nil")
	}

	os.MkdirAll("users/"+conf.Username+"/"+conf.DB.Name, 0644)
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

	return RunContainer(path)
}
