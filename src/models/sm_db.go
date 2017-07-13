package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func getsmDBconnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:roialpha@cloudsql(project-alpha-170622:sessions)/management")
	return db, err
}

func getptDBconnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "activedev:longsitter32@tcp(127.0.0.1:3306)/prompts")
	return db, err
}