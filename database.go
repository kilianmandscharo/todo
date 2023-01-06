package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func newDatabase(name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initDatabase(db *sql.DB) error {
	listStatement := "CREATE TABLE list (id INTEGER PRIMARY KEY ASC)"
	_, err := db.Exec(listStatement)
	if err != nil {
		return err
	}
	itemStatement := "CREATE TABLE item (id INTEGER PRIMARY KEY ASC, content TEXT, list_id INTEGER, FOREIGN KEY (list_id) REFERENCES list (id))"
	_, err = db.Exec(itemStatement)
	if err != nil {
		return err
	}
	return nil
}
