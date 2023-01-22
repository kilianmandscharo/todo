package main

import (
	"fmt"
	"unicode"

	"github.com/gdamore/tcell"
)

type List struct {
	ID    int
	name  string
	row   int
	col   int
	items []Item
}

type Item struct {
	id      int
	content string
	done    bool
}

func (l *List) renderListName(ui *UI) {
	for col, r := range []rune(l.name) {
		var style tcell.Style
		if ui.editListName {
				style = ui.estyle
		} else {
			style = ui.dstyle
		}
		ui.screen.SetContent(xoffset+col, 3, r, nil, style)
	}
}

func (l *List) render(ui *UI) {
	l.renderListName(ui)
	if len(l.items) == 0 {
		ui.renderLine("Press 'n' to create an entry", headerHeight)
		return
	}
	done := 0
	for row, item := range l.items[ui.windowTop:ui.windowBottom] {
		if item.done {
			done++
		}
		rowWithOffset := row + topOffset + headerHeight
		rowWithW := row + ui.windowTop
		var style tcell.Style
		if !ui.editListName && rowWithW == l.row {
			if ui.edit {
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
		ui.screen.SetContent(0+xoffset, rowWithOffset, '[', nil, ui.dstyle)
		ui.screen.SetContent(1+xoffset, rowWithOffset, marker, nil, ui.dstyle)
		ui.screen.SetContent(2+xoffset, rowWithOffset, ']', nil, ui.dstyle)
		for col, r := range []rune(item.content) {
			colWithOffset := col + xoffset + 4
			ui.screen.SetContent(colWithOffset, rowWithOffset, r, nil, style)
		}
	}
	total := len(l.items)
	topLine := fmt.Sprintf("%d / %d done", done, total)
	ui.screen.SetContent(len(l.name)+2, 3, '-', nil, ui.dstyle)
	ui.screen.SetContent(len(l.name)+3, 3, '-', nil, ui.dstyle)
	ui.screen.SetContent(len(l.name)+4, 3, '-', nil, ui.dstyle)
	for col, r := range []rune(topLine) {
		ui.screen.SetContent(len(l.name)+6+col, 3, r, nil, ui.dstyle)
	}
}

func (l *List) down(ui *UI) {
	if l.row+1 <= len(l.items)-1 {
		if l.row+1 == ui.windowBottom {
			ui.windowBottom++
			ui.windowTop++
		}
		l.row++
	}
}

func (l *List) up(ui *UI) {
	if l.row-1 >= 0 {
		if l.row == ui.windowTop {
			ui.windowTop--
			ui.windowBottom--
		}
		l.row--
	}
}

func (l *List) switchUp(ui *UI) {
	i := l.row
	if i-1 >= 0 {
		if l.row == ui.windowTop {
			ui.windowTop--
			ui.windowBottom--
		}
		l.items[i], l.items[i-1] = l.items[i-1], l.items[i]
		l.row--
	}
}

func (l *List) switchDown(ui *UI) {
	i := l.row
	if i+1 <= len(l.items)-1 {
		if l.row+1 == ui.windowBottom {
			ui.windowBottom++
			ui.windowTop++
		}
		l.items[i], l.items[i+1] = l.items[i+1], l.items[i]
		l.row++
	}
}

func (l *List) delete(db *DB, ui *UI) {
	if len(l.items) == 0 {
		return
	}
	err := db.deleteItem(l.currentItem().id)
	if err != nil {
		return
	}
	nitems := len(l.items)
	space := ui.height() - headerHeight - topOffset - bottomOffset
	if nitems-1 < space {
		ui.windowBottom--
	} else if ui.windowTop > 0 {
		ui.windowBottom--
		ui.windowTop--
	}
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

func (l *List) add(db *DB, ui *UI) {
	i := l.row
	id, err := db.createItem(l.ID)
	if err != nil {
		return
	}
	newItem := Item{
		id:      id,
		content: "New Entry",
	}
	nitems := len(l.items)
	space := ui.height() - headerHeight - topOffset - bottomOffset
	if nitems+1 <= space && ui.windowTop < space {
		ui.windowBottom++
	}
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

func (l *List) itemById(id int) *Item {
	for i := range l.items {
		if l.items[i].id == id {
			return &l.items[i]
		}
	}
	return nil
}

func (l *List) addRuneToName(r rune) {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '_' || r == '.' || r == '-' {
		if len(l.name) == 1 && l.name[0] == ' ' {
			l.name = string(r)
		} else {
			l.name += string(r)
		}
		l.col++
	}
}

func (l *List) deleteRuneFromName() {
	last := len(l.name) - 1
	if last > 0 {
		l.name = l.name[:last]
		l.col--
	}
	if last == 0 {
		l.name = " "
	}
}

func (l *List) updateName(db *DB) error {
	if err := db.updateListName(l.name, l.ID); err != nil {
		return err
	}
	return nil
}
