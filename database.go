package main

import (
	"database/sql"
	"encoding/json"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func newDatabase() (*DB, error) {
	os.Mkdir("todo_data", os.ModePerm)
	db, err := sql.Open("sqlite3", "todo_data/data.db")
	if err != nil {
		return nil, err
	}
	q := `
  PRAGMA foreign_keys = ON;
  `
	_, err = db.Exec(q)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (db *DB) init() error {
	_, err := db.db.Exec("CREATE TABLE IF NOT EXISTS list (id INTEGER PRIMARY KEY ASC, name TEXT, item_order TEXT)")
	if err != nil {
		return err
	}
	_, err = db.db.Exec("CREATE TABLE IF NOT EXISTS item (id INTEGER PRIMARY KEY ASC, content TEXT, done INTEGER, list_id INTEGER, FOREIGN KEY(list_id) REFERENCES list(id) ON DELETE CASCADE)")
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) close() {
	db.db.Close()
}

func (db *DB) createList() (int, error) {
	var id int
	row := db.db.QueryRow("INSERT INTO list (id, name, item_order) VALUES (null, 'List name', '') RETURNING id")
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (db *DB) deleteList(id int) error {
	_, err := db.db.Exec("DELETE FROM list WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) updateListName(name string, id int) error {
	_, err := db.db.Exec("UPDATE list SET name = ? WHERE id = ?", name, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) getLists() ([]List, error) {
	rows, err := db.db.Query("SELECT id, name from list")
	if err != nil {
		return nil, err
	}
	var lists []List
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		items, _ := db.getItems(id)
		list := List{
			ID:    id,
			name:  name,
			items: items,
		}
		lists = append(lists, list)
	}
	return lists, nil
}

func (db *DB) createItem(listID int) (int, error) {
	var id int
	row := db.db.QueryRow("INSERT INTO item (content, done, list_id) VALUES ('New Entry', ?, ?) RETURNING id", 0, listID)
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (db *DB) deleteItem(id int) error {
	_, err := db.db.Exec("DELETE FROM item WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) updateItemContent(id int, content string) error {
	_, err := db.db.Exec("UPDATE item SET content = ? WHERE id = ?", content, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) updateItemDone(id int, done bool) error {
	var newDone int
	if done {
		newDone = 1
	} else {
		newDone = 0
	}
	_, err := db.db.Exec("UPDATE item SET done = ? WHERE id = ?", newDone, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) getItems(listID int) ([]Item, error) {
	rows, err := db.db.Query("SELECT id, content, done FROM item WHERE list_id = ?", listID)
	if err != nil {
		return nil, err
	}
	var items []Item
	for rows.Next() {
		var id int
		var content string
		var done int
		if err := rows.Scan(&id, &content, &done); err != nil {
			return nil, err
		}
		items = append(items, Item{id: id, content: content, done: done == 1})
	}
	return items, nil
}

func (db *DB) saveOrder(lists []List) error {
	for _, l := range lists {
		items := l.items
		order := make(map[int]int)
		for i, item := range items {
			order[i] = item.id
		}
		orderString, err := json.Marshal(order)
		if err != nil {
			return err
		}
		_, err = db.db.Exec("UPDATE list SET item_order = ? WHERE id = ?", string(orderString), l.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) loadOrder() ([]map[int]int, error) {
	var orders []map[int]int
	rows, err := db.db.Query("SELECT item_order FROM list")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var orderString string
		if err := rows.Scan(&orderString); err != nil {
			return nil, err
		}
		var order map[int]int
		if err := json.Unmarshal([]byte(orderString), &order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
