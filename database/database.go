package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"code.sentiments/config"
)

var Instance *Connection

func init() {
	if Instance == nil {
		Instance = GetMySQLDB()
	}
}

type Connection struct {
	Conn *sql.DB
}

func GetMySQLDB() (connection *Connection) {
	if config.IsProduction {
		var dbURI string
		dbURI = fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s",
			config.DBConn.DBUser, config.DBConn.DBPwd, config.DBConn.InstanceConnectionName, config.DBConn.DBName)

		dbPool, err := sql.Open("mysql", dbURI)
		if err != nil {
			panic("Failed to initialize the database")
		}

		return &Connection{Conn: dbPool}
	}

	db, err := sql.Open(config.DBConn.DBDriver,
		config.DBConn.DBUser+":"+config.DBConn.DBPwd+"@/"+config.DBConn.DBName)
	if err != nil {
		panic("Failed to initialize the database")
	}

	return &Connection{Conn: db}
}
