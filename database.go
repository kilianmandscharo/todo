package main

import (
	"database/sql"
	"encoding/json"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func newDatabase() (*DB, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }
	os.Mkdir(path.Join(homeDir, ".todo_data"), os.ModePerm)
	db, err := sql.Open("sqlite3", path.Join(homeDir, ".todo_data/data.db"))
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
	_, err := db.db.Exec("CREATE TABLE IF NOT EXISTS ui (id INTEGER, list_order TEXT, UNIQUE(id))")
	if err != nil {
		return err
	}
	_, err = db.db.Exec("INSERT OR IGNORE INTO ui (id, list_order) VALUES(1, '')")
	if err != nil {
		return err
	}
	_, err = db.db.Exec("CREATE TABLE IF NOT EXISTS list (id INTEGER PRIMARY KEY ASC, name TEXT, item_order TEXT)")
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
	listOrder := make(map[int]int)
	for i, l := range lists {
		listOrder[i] = l.ID
		items := l.items
		order := make(map[int]int)
		for j, item := range items {
			order[j] = item.id
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
	listOrderString, err := json.Marshal(listOrder)
	if err != nil {
		return err
	}
	_, err = db.db.Exec("UPDATE ui SET list_order = ? WHERE id = ?", string(listOrderString), 1)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) loadListOrder() (map[int]int, error) {
	row := db.db.QueryRow("SELECT list_order FROM ui WHERE id = 1")
	var orderString string
	if err := row.Scan(&orderString); err != nil {
		return nil, err
	}
	if len(orderString) == 0 {
		return nil, nil
	}
	var order map[int]int
	if err := json.Unmarshal([]byte(orderString), &order); err != nil {
		return nil, err
	}
	return order, nil
}

func (db *DB) loadItemOrder(listId int) (map[int]int, error) {
	row := db.db.QueryRow("SELECT item_order FROM list WHERE id = ?", listId)
	var orderString string
	if err := row.Scan(&orderString); err != nil {
		return nil, err
	}
	if len(orderString) == 0 {
		return nil, nil
	}
	var order map[int]int
	if err := json.Unmarshal([]byte(orderString), &order); err != nil {
		return nil, err
	}
	return order, nil
}
