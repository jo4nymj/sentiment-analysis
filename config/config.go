package config

import (
	"database/sql"
	"fmt"
	"os"

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
	/*
		var (
			dbUser                 = MustGetenv("DB_USER")
			dbPwd                  = MustGetenv("DB_PWD")
			instanceConnectionName = MustGetenv("INSTANCE_CONNECTION_NAME")
			dbName                 = MustGetenv("DB_NAME")
		)

		var dbURI string
		dbURI = fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s", dbUser, dbPwd, instanceConnectionName, dbName)

		dbPool, err := sql.Open("mysql", dbURI)
		if err != nil {
			fmt.Printf("Cannot access to database")
			//fmt.Errorf("sql.Open: %v", err)
			return nil
		}

		return &Connection{Conn: dbPool}
	*/
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

func MustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		fmt.Printf("Warning: %s environment variable not set.\n", k)
	}
	return v
}
