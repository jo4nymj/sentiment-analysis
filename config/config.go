package config

import (
	"fmt"
	"os"
)

var (
	Key          string
	IsProduction bool
	DBConn       DatabaseConnection
)

type DatabaseConnection struct {
	DBUser                 string
	DBPwd                  string
	DBName                 string
	DBDriver               string
	InstanceConnectionName string
}

func init() {
	if mustGetenv("PRODUCTION") == "true" {
		IsProduction = true

		Key = mustGetenv("KEY")

		DBConn.DBUser = mustGetenv("DB_USER")
		DBConn.DBPwd = mustGetenv("DB_PWD")
		DBConn.InstanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME")
		DBConn.DBName = mustGetenv("DB_NAME")
		return
	}
	DBConn.DBUser = "root"
	DBConn.DBPwd = ""
	DBConn.DBDriver = "mysql"
	DBConn.DBName = "project"
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		fmt.Printf("Warning: %s environment variable not set.\n", k)
	}
	return v
}
