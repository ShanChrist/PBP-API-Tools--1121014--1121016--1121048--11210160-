package controllers

import (
	"database/sql"
	"log"
)

func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/db_explorasi_tools")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
