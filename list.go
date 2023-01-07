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
	ui    *UI
}

type Item struct {
	id       int
	position int
	content  string
	done     bool
}

func (l *List) enterEdit() {
	l.edit = true
	l.items[l.row].content = " "
}

func (l *List) exitEdit() {
	l.edit = false
	l.col = 0
}

func (l *List) render() {
	for row, item := range l.items {
		var style tcell.Style
		if row == l.row {
			if l.edit {
				style = l.ui.estyle
			} else {
				style = l.ui.hstyle
			}
		} else {
			style = l.ui.dstyle
		}
		var marker rune
		if item.done {
			marker = 'X'
		} else {
			marker = ' '
		}
		l.ui.screen.SetContent(0, row, '[', nil, l.ui.dstyle)
		l.ui.screen.SetContent(1, row, marker, nil, l.ui.dstyle)
		l.ui.screen.SetContent(2, row, ']', nil, l.ui.dstyle)
		for col, r := range []rune(item.content) {
			l.ui.screen.SetContent(col+4, row, r, nil, style)
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

func (l *List) delete() {
	l.ui.db.deleteItem(l.currentItem().id)
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

func (l *List) add() {
	i := l.row
	content := "New Entry"
	var position int
	if len(l.items) == 0 {
		position = 0
	} else {
		position = i + 1
    l.row++
	}
	id, _ := l.ui.db.createItem(content, position, l.ID)
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

func (l *List) markItem() {
	item := l.currentItem()
	item.done = !item.done
}

func (l *List) currentItem() *Item {
	return &l.items[l.row]
}
