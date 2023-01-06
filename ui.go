package main

import (
	"log"
	"os"
	"github.com/gdamore/tcell"
)

type UI struct {
	screen  tcell.Screen
	lists   []List
	current int
}

func newUI() UI {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Init(); err != nil {
		log.Fatal(err)
	}
	return UI{screen: s}
}

func (ui *UI) clear() {
	ui.screen.Clear()
}

func (ui *UI) show() {
	ui.screen.Show()
}

func (ui *UI) addList() {
	dstyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	hstyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
	estyle := tcell.StyleDefault.Background(tcell.ColorPurple).Foreground(tcell.ColorBlack)
	ui.screen.SetStyle(dstyle)
	todos := []string{
		"New Entry",
	}
	list := List{row: 0, col: 0, edit: false, items: todos, done: []bool{false}, dstyle: dstyle, hstyle: hstyle, estyle: estyle, screen: ui.screen}
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
				list.delete()
			}
			if ev.Rune() == 'n' {
				list.add()
			}
			if ev.Rune() == 'e' {
				list.enterEdit()
			}
			if ev.Rune() == 13 {
				list.markItem()
			}
		}
	}
}

