package main

import (
	"unicode"

	"github.com/gdamore/tcell"
)

type List struct {
	ID    int
	row   int
	col   int
	edit  bool
	items []Item
}

type Item struct {
	id      int
	content string
	done    bool
}

func (l *List) enterEdit() {
	l.edit = true
	l.items[l.row].content = " "
}

func (l *List) exitEdit() {
	l.edit = false
	l.col = 0
}

func (l *List) render(ui *UI) {
	for row, item := range l.items {
		var style tcell.Style
		if row == l.row {
			if l.edit {
				style = ui.estyle
			} else {
				style = ui.hstyle
			}
		} else {
			style = ui.dstyle
		}
		var marker rune
		if item.done {
			marker = 'X'
		} else {
			marker = ' '
		}
		ui.screen.SetContent(0, row, '[', nil, ui.dstyle)
		ui.screen.SetContent(1, row, marker, nil, ui.dstyle)
		ui.screen.SetContent(2, row, ']', nil, ui.dstyle)
		for col, r := range []rune(item.content) {
			ui.screen.SetContent(col+4, row, r, nil, style)
		}
	}
}

func (l *List) down() {
	if l.row+1 <= len(l.items)-1 {
		l.row++
	}
}

func (l *List) up() {
	if l.row-1 >= 0 {
		l.row--
	}
}

func (l *List) switchUp() {
	i := l.row
	if i-1 >= 0 {
		l.items[i], l.items[i-1] = l.items[i-1], l.items[i]
		l.row--
	}
}

func (l *List) switchDown() {
	i := l.row
	if i+1 <= len(l.items)-1 {
		l.items[i], l.items[i+1] = l.items[i+1], l.items[i]
		l.row++
	}
}

func (l *List) delete(db *DB) {
	db.deleteItem(l.currentItem().id)
	if len(l.items) == 1 {
		l.items = nil
		return
	}
	i := l.row
	if l.row == len(l.items)-1 {
		l.row--
	}
	newItems := l.items[:i]
	newItems = append(newItems, l.items[i+1:]...)
	l.items = newItems
}

func (l *List) add(db *DB) {
	i := l.row
	content := "New Entry"
	id, _ := db.createItem(content, l.ID)
	newItem := Item{
		id:      id,
		content: content,
	}
	nitems := len(l.items)
	if nitems == 0 || nitems-1 == i {
		l.items = append(l.items, newItem)
		return
	}
	l.items = append(l.items[:i+1], l.items[i:]...)
	l.items[i+1] = newItem
}

func (l *List) addRune(r rune) {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '_' || r == '.' || r == '-' {
		if len(l.items[l.row].content) == 1 && l.items[l.row].content[0] == ' ' {
			l.items[l.row].content = string(r)
		} else {
			l.items[l.row].content += string(r)
		}
		l.col++
	}
}

func (l *List) deleteRune() {
	last := len(l.items[l.row].content) - 1
	if last > 0 {
		l.items[l.row].content = l.items[l.row].content[:last]
		l.col--
	}
	if last == 0 {
		l.items[l.row].content = " "
	}
}

func (l *List) markItem(db *DB) {
	item := l.currentItem()
	item.done = !item.done
	db.updateItemDone(item.id, item.done)
}

func (l *List) updateItem(db *DB) {
	item := l.currentItem()
	db.updateItemContent(item.id, item.content)
}

func (l *List) currentItem() *Item {
	return &l.items[l.row]
}
