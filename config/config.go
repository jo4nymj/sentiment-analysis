package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Connection struct {
	Conn *sql.DB
}

var Instance *Connection

func GetMySQLDB() (connection *Connection, err error) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "project"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		return nil, err
	}

	return &Connection{Conn: db}, nil
}

func init() {
	GetMySQLDB()
}
