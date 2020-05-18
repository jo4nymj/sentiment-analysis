package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	if Instance == nil {
		Instance = GetMySQLDB()
	}
}

type Connection struct {
	Conn *sql.DB
}

var Instance *Connection

func GetMySQLDB() (connection *Connection) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "project"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic("Failed to initialize the database")
	}

	return &Connection{Conn: db}
}
