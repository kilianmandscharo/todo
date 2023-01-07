package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func newDatabase(name string) (*DB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (db *DB) init() error {
	listStatement := "CREATE TABLE IF NOT EXISTS list (id INTEGER PRIMARY KEY ASC)"
	_, err := db.db.Exec(listStatement)
	if err != nil {
		return err
	}
	itemStatement := "CREATE TABLE IF NOT EXISTS item (id INTEGER PRIMARY KEY ASC, content TEXT, list_id INTEGER, FOREIGN KEY (list_id) REFERENCES list (id))"
	_, err = db.db.Exec(itemStatement)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) close() {
	db.db.Close()
}

func (db *DB) createList() error {
	_, err := db.db.Exec("INSERT INTO list (id) VALUES (null)")
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) deleteList(id int) error {
	stmt := fmt.Sprintf("DELETE FROM list WHERE id = %d", id)
	_, err := db.db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) getLists() ([]List, error) {
	stmt := "SELECT id from list"
	rows, err := db.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	var lists []List
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items, _ := db.getItems(id)
		list := List{
			ID:    id,
			items: items,
		}
		lists = append(lists, list)
	}
	return lists, nil
}

func (db *DB) createItem(item string, listID int) (int, error) {
	stmt := fmt.Sprintf("INSERT INTO item (content, list_id) VALUES (\"%s\", %d) RETURNING id", item, listID)
  var id int
	row := db.db.QueryRow(stmt)
  row.Scan(&id)
	return id, nil
}

func (db *DB) deleteItem(id int) error {
	stmt := fmt.Sprintf("DELETE FROM item WHERE id = %d", id)
	_, err := db.db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) updateItem(id int, content string) error {
	stmt := fmt.Sprintf("UPDATE item SET content = \"%s\" WHERE id = %d", content, id)
	_, err := db.db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) getItems(listID int) ([]Item, error) {
	stmt := fmt.Sprintf("SELECT id, content FROM item WHERE list_id = \"%d\"", listID)
	rows, err := db.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	var items []Item
	for rows.Next() {
		var id int
		var content string
		if err := rows.Scan(&id, &content); err != nil {
			return nil, err
		}
		items = append(items, Item{id: id, content: content})
	}
	rows.Close()
	return items, nil
}
