package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"

	"github.com/gdamore/tcell"
)

const headerHeight = 4
const leftOffset = 1
const topOffset = 1
const bottomOffset = 1

const navPosition = 1

type UI struct {
	screen       tcell.Screen
	db           *DB
	lists        []List
	current      int
	dstyle       tcell.Style
	hstyle       tcell.Style
	estyle       tcell.Style
	pstyle       tcell.Style
	successStyle tcell.Style
	errorStyle   tcell.Style
	edit         bool
	windowTop    int
	windowBottom int
	editListName bool
	deletingList bool
}

func newUI(debug bool) UI {
	var s tcell.Screen
	if !debug {
		screen, err := tcell.NewScreen()
		if err != nil {
			log.Fatal(err)
		}
		if err := screen.Init(); err != nil {
			log.Fatal(err)
		}
		s = screen
	}
	db, err := newDatabase()
	if err != nil {
		log.Fatal(err)
	}
	db.init()
	ui := &UI{screen: s, db: db}
	ui = setStyles(ui)
	return *ui
}

func (ui *UI) load() {
	lists, err := ui.db.getLists()
	if err != nil {
		log.Fatal(err)
	}
	ui.lists = lists
	ui.loadOrder()
	ui.calculateWindow()
}

func (ui *UI) clear() {
	ui.screen.Clear()
}

func (ui *UI) show() {
	ui.screen.Show()
}

func (ui *UI) currentList() *List {
	return &ui.lists[ui.current]
}

func (ui *UI) switchList(r rune) {
	val := int(r - '0')
	if val > len(ui.lists) {
		return
	}
	ui.current = val - 1
	ui.calculateWindow()
}

func (ui *UI) addList() {
	if len(ui.lists) == 9 {
		return
	}
	id, err := ui.db.createList()
	if err != nil {
		return
	}
	ui.lists = append(ui.lists, List{ID: id, name: "List name"})
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
	ui.calculateWindow()
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
		logToFile(fmt.Sprintf("Order: %v", order))
		var newItems []Item
		for index := 0; index < len(order); index++ {
			id := order[index]
			item := *ui.lists[i].itemById(id)
			newItems = append(newItems, item)
		}
		ui.lists[i].items = newItems
	}
}

func (ui *UI) handleEvent(ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventResize:
		ui.calculateWindow()
		ui.screen.Sync()
		return
	case *tcell.EventKey:
		ui.clear()
		ui.show()
		if ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'x' {
			ui.saveOrder()
			ui.screen.Fini()
			os.Exit(0)
		}
		if ui.editListName {
			if ev.Key() == tcell.KeyEscape {
				ui.exitNameEdit()
				return
			}
			if ev.Rune() == 127 {
				ui.currentList().deleteRuneFromName()
			} else {
				ui.currentList().addRuneToName(ev.Rune())
			}
			return
		}
		if ui.edit {
			if ev.Key() == tcell.KeyEscape {
				ui.exitEdit()
				return
			}
		} else {
			if ui.deletingList {
				if ev.Rune() == 'y' {
					ui.deleteList()
					ui.deletingList = false
					return
				}
				if ev.Rune() == 'n' {
					ui.deletingList = false
					return
				}
				return
			}
			if ev.Rune() == 'c' {
				ui.addList()
				return
			}
			if ev.Rune() == 'r' {
				if len(ui.lists) != 0 {
					ui.deletingList = true
				}
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
				if ev.Rune() == 'b' {
					ui.enterNameEdit()
				}
				if ev.Rune() == 'j' {
					list.down(ui)
					return
				}
				if ev.Rune() == 'k' {
					list.up(ui)
					return
				}
				if ev.Rune() == 'J' {
					list.switchDown(ui)
					return
				}
				if ev.Rune() == 'K' {
					list.switchUp(ui)
					return
				}
				if ev.Rune() == 'd' {
					list.delete(ui.db, ui)
					return
				}
				if ev.Rune() == 'n' {
					list.add(ui.db, ui)
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

func (ui *UI) render() {
	renderListNav(ui)
	renderCurrentList(ui)
	renderFooter(ui)
}

func renderCurrentList(ui *UI) {
	if len(ui.lists) == 0 {
		ui.renderLine("Press c to create a new list", headerHeight-2)
		return
	}
	ui.lists[ui.current].render(ui)
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
		ui.screen.SetContent(i*2+leftOffset, 1, []rune(r)[0], nil, style)
	}
	if len(ui.lists) != 0 {
		ui.screen.SetContent(ui.current*2+leftOffset, 2, '^', nil, ui.pstyle)
	}
	ui.renderLine(separator(ui), 3)
}

func separator(ui *UI) string {
	var line bytes.Buffer
	for i := 0; i < ui.width()-2; i++ {
		line.WriteRune('=')
	}
	return line.String()
}

func (ui *UI) renderLine(line string, row int) {
	for col, r := range []rune(line) {
		ui.screen.SetContent(col+leftOffset, row+topOffset, r, nil, ui.dstyle)
	}
}

func renderFooter(ui *UI) {
	if ui.deletingList {
		ui.renderLine("Delete current list? y / n", ui.height()-2)
	} else {
		ui.renderLine(separator(ui), ui.height()-2)
	}
}

func (ui *UI) enterEdit() {
	l := ui.currentList()
	l.items[l.row].content = " "
	ui.edit = true
}

func (ui *UI) exitEdit() {
	ui.edit = false
	ui.currentList().updateItem(ui.db)
	ui.currentList().col = 0
}

func (ui *UI) enterNameEdit() {
	ui.currentList().name = " "
	ui.editListName = true
}

func (ui *UI) exitNameEdit() {
	ui.editListName = false
	ui.currentList().updateName(ui.db)
	ui.currentList().col = 0
}

func (ui *UI) calculateWindow() {
	if len(ui.lists) == 0 {
		return
	}
	_, h := ui.screen.Size()
	topOffset := topOffset + headerHeight
	bottomOffset := 1
	listLength := len(ui.currentList().items)
	spaceNeeded := listLength + topOffset + bottomOffset
	var newWindowBottom int
	if spaceNeeded > h {
		newWindowBottom = max(listLength-(spaceNeeded-h), 0)
	} else {
		newWindowBottom = max(listLength, 0)
	}
	ui.windowBottom = newWindowBottom
}

func (ui *UI) closeDB() {
	ui.db.close()
}

func (ui *UI) height() int {
	_, h := ui.screen.Size()
	return h
}

func (ui *UI) width() int {
	w, _ := ui.screen.Size()
	return w
}

func (ui *UI) event() tcell.Event {
	return ui.screen.PollEvent()
}

func setStyles(ui *UI) *UI {
	dstyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	hstyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
	estyle := tcell.StyleDefault.Background(tcell.ColorPurple).Foreground(tcell.ColorBlack)
	pstyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlue)
	successStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)
	errorStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
	ui.dstyle = dstyle
	ui.hstyle = hstyle
	ui.estyle = estyle
	ui.pstyle = pstyle
	ui.successStyle = successStyle
	ui.errorStyle = errorStyle
	ui.screen.SetStyle(dstyle)
	return ui
}
