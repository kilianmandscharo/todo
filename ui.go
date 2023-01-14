package main

import (
	"log"
	"os"
	"strconv"
	"unicode"

	"github.com/gdamore/tcell"
)

type UI struct {
	screen  tcell.Screen
	db      *DB
	lists   []List
	current int
	dstyle  tcell.Style
	hstyle  tcell.Style
	estyle  tcell.Style
	edit    bool
}

func newUI(debug bool) UI {
	dstyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	hstyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
	estyle := tcell.StyleDefault.Background(tcell.ColorPurple).Foreground(tcell.ColorBlack)
	var s tcell.Screen
	if !debug {
		screen, err := tcell.NewScreen()
		if err != nil {
			log.Fatal(err)
		}
		if err := screen.Init(); err != nil {
			log.Fatal(err)
		}
		screen.SetStyle(dstyle)
		s = screen
	}
	db, err := newDatabase("data.db")
	if err != nil {
		log.Fatal(err)
	}
	db.init()
	return UI{
		screen: s,
		db:     db,
		dstyle: dstyle,
		hstyle: hstyle,
		estyle: estyle,
	}
}

func (ui *UI) load() {
	lists, err := ui.db.getLists()
	if err != nil {
		log.Fatal(err)
	}
	ui.lists = lists
	ui.loadOrder()
}

func (ui *UI) clear() {
	ui.screen.Clear()
}

func (ui *UI) show() {
	ui.screen.Show()
}

func (ui *UI) addList() {
	if len(ui.lists) == 9 {
		return
	}
	id, err := ui.db.createList()
	if err != nil {
		return
	}
	ui.lists = append(ui.lists, List{ID: id})
}

func (ui *UI) deleteList() {
	if len(ui.lists) == 0 {
		return
	}
	err := ui.db.deleteList(ui.currentList().ID)
	if err != nil {
		return
	}
	if len(ui.lists) == 1 {
		ui.lists = nil
		return
	}
	i := ui.current
	if i == len(ui.lists)-1 {
		ui.current--
	}
	newLists := ui.lists[:i]
	newLists = append(newLists, ui.lists[i+1:]...)
	ui.lists = newLists
}

func (ui *UI) currentList() *List {
	return &ui.lists[ui.current]
}

func (ui *UI) handleEvent(ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventResize:
		ui.screen.Sync()
	case *tcell.EventKey:
		ui.clear()
		ui.show()
		if ev.Key() == tcell.KeyCtrlC {
			ui.saveOrder()
			ui.screen.Fini()
			os.Exit(0)
		}
		if ui.edit {
			if ev.Key() == tcell.KeyEscape {
				ui.currentList().updateItem(ui.db)
				ui.exitEdit()
				return
			}
		} else {
			if ev.Rune() == 'c' {
				ui.addList()
				return
			}
			if ev.Rune() == 'r' {
				ui.deleteList()
				return
			}
			if unicode.IsDigit(ev.Rune()) {
				ui.switchList(ev.Rune())
				return
			}
			if ev.Rune() == 'e' {
				ui.enterEdit()
				return
			}
		}
		if len(ui.lists) != 0 {
			list := ui.currentList()
			if ui.edit {
				if ev.Rune() == 127 {
					list.deleteRune()
				} else {
					list.addRune(ev.Rune())
				}
				return
			} else {
				if ev.Rune() == 'j' {
					list.down()
					return
				}
				if ev.Rune() == 'k' {
					list.up()
					return
				}
				if ev.Rune() == 'J' {
					list.switchDown()
					return
				}
				if ev.Rune() == 'K' {
					list.switchUp()
					return
				}
				if ev.Rune() == 'd' {
					list.delete(ui.db)
					return
				}
				if ev.Rune() == 'n' {
					list.add(ui.db)
					return
				}
				if ev.Rune() == 13 {
					list.markItem(ui.db)
					return
				}
			}
		}
	}
}

func (ui *UI) saveOrder() {
	err := ui.db.saveOrder(ui.lists)
	if err != nil {
		log.Fatal(err)
	}
}

func (ui *UI) loadOrder() {
	orders, err := ui.db.loadOrder()
	if err != nil {
		log.Fatal(err)
	}
	for i, order := range orders {
		if len(order) == 0 {
			continue
		}
		var newItems []Item
		for _, id := range order {
			newItems = append(newItems, *ui.lists[i].itemById(id))
		}
		ui.lists[i].items = newItems
	}
}

func (ui *UI) switchList(r rune) {
	val := int(r - '0')
	ui.screen.SetContent(0, 1, r, nil, ui.dstyle)
	if val > len(ui.lists) {
		return
	}
	ui.current = val - 1
}

func (ui *UI) render() {
	renderHeader(ui)
	renderListNav(ui)
	renderCurrentList(ui)
}

func renderCurrentList(ui *UI) {
	if len(ui.lists) == 0 {
		ui.renderLine("Press c to create a new list", 3)
		return
	}
	ui.lists[ui.current].render(ui, 1, 4)
}

func renderListNav(ui *UI) {
	var style tcell.Style
	for i := range ui.lists {
		if i == ui.current {
			style = ui.hstyle
		} else {
			style = ui.dstyle
		}
		r := strconv.Itoa(i + 1)
		ui.screen.SetContent(i*2+1, 2, []rune(r)[0], nil, style)
	}
}

func renderHeader(ui *UI) {
	header := "========== todo =========="
	ui.renderLine(header, 0)
}

func (ui *UI) renderLine(line string, row int) {
	for col, r := range []rune(line) {
		ui.screen.SetContent(col+1, row+1, r, nil, ui.dstyle)
	}
}

func (ui *UI) enterEdit() {
	l := ui.currentList()
	l.items[l.row].content = " "
	ui.edit = true
}

func (ui *UI) exitEdit() {
	ui.edit = false
	ui.currentList().col = 0
}
