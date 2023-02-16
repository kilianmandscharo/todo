package main

import (
	"fmt"

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

func (l *List) render(ui *UI) {
	renderHeader(ui, l)
	renderBody(ui, l)
}

func renderHeader(ui *UI, l *List) {
	for col, r := range []rune(l.name + " ") {
		var style tcell.Style
		if ui.mode == editListNameMode && col == l.col {
			style = lightSkyBlueBlack
		} else {
			style = blackWhite
		}
		ui.screen.SetContent(leftOffset+col, 3, r, nil, style)
	}

	total := len(l.items)
	var done int
	var topLine string
	if total == 0 {
		topLine = ""
	} else {
		for i := range l.items {
			if l.items[i].done {
				done++
			}
		}
		topLine = fmt.Sprintf("%d / %d done", done, total)
	}
	renderTopSeparator(ui, separator(ui, topLine), 5)
}

func renderBody(ui *UI, l *List) {
	if len(l.items) == 0 {
		ui.renderLine("Press e + n to create an entry", headerHeight)
	}
	for row, item := range l.items[ui.windowTop:ui.windowBottom] {
		rowWithOffset := row + topOffset + headerHeight
		rowWithW := row + ui.windowTop
		var style tcell.Style
		if ui.mode != editListNameMode && rowWithW == l.row {
			if ui.mode == editMode {
				style = blackWhite
			} else {
				style = whiteBlack
			}
		} else {
			style = blackWhite
		}
		var marker rune
		if item.done {
			marker = 'X'
		} else {
			marker = ' '
		}
		ui.screen.SetContent(0+leftOffset, rowWithOffset, '[', nil, blackWhite)
		ui.screen.SetContent(1+leftOffset, rowWithOffset, marker, nil, blackWhite)
		ui.screen.SetContent(2+leftOffset, rowWithOffset, ']', nil, blackWhite)
		var content string
		if rowWithW == l.row && ui.mode == editMode {
			content = item.content + " "
		} else {
			content = item.content
		}
		for col, r := range []rune(content) {
			colWithOffset := col + leftOffset + 4
			if col == l.col && rowWithW == l.row && ui.mode == editMode {
				ui.screen.SetContent(colWithOffset, rowWithOffset, r, nil, lightSkyBlueBlack)
			} else {
				ui.screen.SetContent(colWithOffset, rowWithOffset, r, nil, style)
			}
		}
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

func (l *List) updateName(db *DB) error {
	if err := db.updateListName(l.name, l.ID); err != nil {
		return err
	}
	return nil
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
	space := ui.listSpaceAvailable() 
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
		content: " ",
	}
	nitems := len(l.items)
	space := ui.listSpaceAvailable()
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
	if r == ' ' && l.col == 0 {
		return
	}
	content := l.currentItem().content
	if len(content) == 1 && content[0] == ' ' {
		l.items[l.row].content = string(r)
	} else {
		l.items[l.row].content = content[:l.col] + string(r) + content[l.col:]
	}
	l.col++
}

func (l *List) deleteRune() {
	item := l.currentItem()
	content := item.content
	if l.col > 0 {
		item.content = content[:l.col-1] + content[l.col:]
		l.col--
		if l.col == 0 && len(item.content) == 0 {
			item.content = " "
		}
	}
}

func (l *List) addRuneToName(r rune) {
	if r == ' ' && l.col == 0 {
		return
	}
	if len(l.name) == 1 && l.name[0] == ' ' {
		l.name = string(r)
	} else {
		l.name = l.name[:l.col] + string(r) + l.name[l.col:]
	}
	l.col++
}

func (l *List) deleteRuneFromName() {
	if l.col > 0 {
		l.name = l.name[:l.col-1] + l.name[l.col:]
		l.col--
		if l.col == 0 && len(l.name) == 0 {
			l.name = " "
		}
	}
}

func (l *List) markItem(db *DB) {
	if len(l.items) == 0 {
		return
	}
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

func (l *List) cursorLeftEntry() {
	if l.col > 0 {
		l.col--
	}
}

func (l *List) cursorRightEntry() {
	if l.col < len(l.currentItem().content) {
		l.col++
	}
}

func (l *List) cursorLeftListName() {
	if l.col > 0 {
		l.col--
	}
}

func (l *List) cursorRightListName() {
	if l.col < len(l.name) {
		l.col++
	}
}
