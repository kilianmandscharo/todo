package main

import (
	"github.com/gdamore/tcell"
	"log"
	"os"
)

type UI struct {
	screen  tcell.Screen
	db      *DB
	lists   []List
	current int
	dstyle  tcell.Style
	hstyle  tcell.Style
	estyle  tcell.Style
}

func newUI() UI {
	dstyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	hstyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
	estyle := tcell.StyleDefault.Background(tcell.ColorPurple).Foreground(tcell.ColorBlack)
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Init(); err != nil {
		log.Fatal(err)
	}
	s.SetStyle(dstyle)
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
  // Create new list if there are none in the database
	if len(lists) == 0 {
		ui.db.createList()
		lists = append(lists, List{ID: 1})
	}
  // Create a new item for each empty list and set the ui field for each list
	for i := range lists {
		if len(lists[i].items) == 0 {
			ui.db.createItem("New Entry", lists[i].ID)
			lists[i].items = append(lists[i].items, Item{id: 1, content: "New Entry"})
		}
	}
	ui.lists = lists
}

func (ui *UI) clear() {
	ui.screen.Clear()
}

func (ui *UI) show() {
	ui.screen.Show()
}

func (ui *UI) addList() {
	todos := []Item{
		{done: false, content: "New Entry"},
	}
	list := List{
		row:   0,
		col:   0,
		edit:  false,
		items: todos,
	}
	ui.lists = append(ui.lists, list)
}

func (ui *UI) currentList() *List {
	return &ui.lists[ui.current]
}

func (ui *UI) handleEvent(ev tcell.Event) {
	list := ui.currentList()
	switch ev := ev.(type) {
	case *tcell.EventResize:
		ui.screen.Sync()
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyCtrlC {
			ui.screen.Fini()
			os.Exit(0)
		}
		if list.edit {
			list.addRune(ev.Rune())
			if ev.Rune() == 127 {
				list.deleteRune()
			}
			if ev.Key() == tcell.KeyEscape {
        list.updateItem(ui.db)
				list.exitEdit()
			}
		} else {
			if ev.Rune() == 'j' {
				list.down()
			}
			if ev.Rune() == 'k' {
				list.up()
			}
			if ev.Rune() == 'J' {
				list.switchDown()
			}
			if ev.Rune() == 'K' {
				list.switchUp()
			}
			if ev.Rune() == 'd' {
				list.delete(ui.db)
			}
			if ev.Rune() == 'n' {
				list.add(ui.db)
			}
			if ev.Rune() == 'e' {
				list.enterEdit()
			}
			if ev.Rune() == 13 {
				list.markItem(ui.db)
			}
		}
	}
}
