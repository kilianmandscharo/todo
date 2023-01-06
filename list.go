package main

import (
	"unicode"
	"github.com/gdamore/tcell"
)

type List struct {
	row    int
	col    int
	edit   bool
	items  []string
	done   []bool
	dstyle tcell.Style
	hstyle tcell.Style
	estyle tcell.Style
	screen tcell.Screen
}

func (l *List) enterEdit() {
	l.edit = true
	l.items[l.row] = " "
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
				style = l.estyle
			} else {
				style = l.hstyle
			}
		} else {
			style = l.dstyle
		}
		var marker rune
		if l.done[row] {
			marker = 'X'
		} else {
			marker = ' '
		}
		l.screen.SetContent(0, row, '[', nil, l.dstyle)
		l.screen.SetContent(1, row, marker, nil, l.dstyle)
		l.screen.SetContent(2, row, ']', nil, l.dstyle)
		for col, r := range []rune(item) {
			l.screen.SetContent(col+4, row, r, nil, style)
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
		l.done[i], l.done[i-1] = l.done[i-1], l.done[i]
		l.row--
	}
}

func (l *List) switchDown() {
	i := l.row
	if i+1 <= len(l.items)-1 {
		l.items[i], l.items[i+1] = l.items[i+1], l.items[i]
		l.done[i], l.done[i+1] = l.done[i+1], l.done[i]
		l.row++
	}
}

func (l *List) delete() {
	if len(l.items) == 1 {
		l.items = nil
		l.done = nil
		return
	}
	i := l.row
	if l.row == len(l.items)-1 {
		l.row--
	}
	newItems := l.items[:i]
	newItems = append(newItems, l.items[i+1:]...)
	l.items = newItems
	newDone := l.done[:i]
	newDone = append(newDone, l.done[i+1:]...)
	l.done = newDone
}

func (l *List) add() {
	newItem := "New Entry"
	i := l.row
	nitems := len(l.items)
	if nitems == 0 || nitems-1 == i {
		l.done = append(l.done, false)
		l.items = append(l.items, newItem)
		return
	}
	l.items = append(l.items[:i+1], l.items[i:]...)
	l.items[i+1] = newItem
	l.done = append(l.done[:i+1], l.done[i:]...)
	l.done[i+1] = false
}

func (l *List) addRune(r rune) {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '_' || r == '.' || r == '-' {
		if len(l.items[l.row]) == 1 && l.items[l.row][0] == ' ' {
			l.items[l.row] = string(r)
		} else {

			l.items[l.row] += string(r)
		}
		l.col++
	}
}

func (l *List) deleteRune() {
	last := len(l.items[l.row]) - 1
	if last > 0 {
		l.items[l.row] = l.items[l.row][:last]
		l.col--
	}
	if last == 0 {
		l.items[l.row] = " "
	}
}

func (l *List) markItem() {
	l.done[l.row] = !l.done[l.row]
}

