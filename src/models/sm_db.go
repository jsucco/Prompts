package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"log"
	"fmt"
)

func getsmDBconnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:roialpha@cloudsql(project-alpha-170622:sessions)/management")
	return db, err
}

func getSessionsConnection() (*sql.DB, error) {

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