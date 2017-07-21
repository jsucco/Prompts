package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"log"
	_"fmt"
	"fmt"
)

func getProxyConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:roialpha@tcp(127.0.0.1:3306)/management")
	return db, err
}

func getSessionsConnection() (*sql.DB, error) {

	if mustGetenv("DEBUG-MODE") == "true" {
		return getProxyConnection()
	}

	connectionName := mustGetenv("CLOUDSQL_CONNECTION_NAME")
	user := mustGetenv("CLOUDSQL_USER")
	password := os.Getenv("CLOUDSQL_PASSWORD")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@cloudsql(%s)/", user, password, connectionName))

	return db, err
}

func getptDBconnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "activedev:longsitter32@tcp(127.0.0.1:3306)/prompts")
	return db, err
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return v
}