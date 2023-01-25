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

type Mode int

const (
	normalMode Mode = iota
	listMode
	entryMode
	editListNameMode
	deleteListMode
	editMode
)

const headerHeight = 6
const leftOffset = 1
const topOffset = 1
const bottomOffset = 1

const navPosition = 1

type UI struct {
	screen       tcell.Screen
	db           *DB
	lists        []List
	current      int
	windowTop    int
	windowBottom int
	styles       Styles
	mode         Mode
}

type Styles struct {
	success   tcell.Style
	error     tcell.Style
	primary   tcell.Style
	highlight tcell.Style
	def       tcell.Style
	edit      tcell.Style
    rib       tcell.Style
    nav       tcell.Style 
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
	ui.mode = normalMode
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
    ui.windowBottom = 0
    ui.windowTop = 0
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

		if ev.Key() == tcell.KeyCtrlC {
			ui.exit()
		}

		switch ui.mode {
		case normalMode:
			handleNormalModeEv(ui, ev.Key(), ev.Rune())
		case listMode:
			handleListModeEv(ui, ev.Key(), ev.Rune())
		case entryMode:
			handleEntryModeEv(ui, ev.Key(), ev.Rune())
		case editListNameMode:
			handleEditListNameModeEv(ui, ev.Key(), ev.Rune())
		case deleteListMode:
			handleDeleteListModeEv(ui, ev.Key(), ev.Rune())
		case editMode:
			handleEditMode(ui, ev.Key(), ev.Rune())
		}
	}
}

func handleNormalModeEv(ui *UI, key tcell.Key, r rune) {
	if r == 'x' {
		ui.exit()
	}
	if r == 'l' {
		ui.mode = listMode
	}
	if len(ui.lists) != 0 {
		list := ui.currentList()
		if unicode.IsDigit(r) {
			ui.switchList(r)
		} else if r == 'j' {
			list.down(ui)
		} else if r == 'k' {
			list.up(ui)
		} else if r == 'J' {
			list.switchDown(ui)
		} else if r == 'K' {
			list.switchUp(ui)
		} else if r == 13 {
			list.markItem(ui.db)
		} else if r == 'e' {
			ui.mode = entryMode
		}
	}
}

func handleListModeEv(ui *UI, key tcell.Key, r rune) {
	if r == 'd' && len(ui.lists) != 0 {
		ui.mode = deleteListMode
	} else if r == 'n' {
		ui.addList()
		ui.mode = normalMode
	} else if r == 'e' && len(ui.lists) != 0 {
		ui.enterNameEdit()
	} else if r == 'b' {
		ui.mode = normalMode
	}
}

func handleEntryModeEv(ui *UI, key tcell.Key, r rune) {
	if r == 'd' && len(ui.currentList().items) != 0 {
		ui.currentList().delete(ui.db, ui)
		ui.mode = normalMode
	} else if r == 'n' {
		ui.currentList().add(ui.db, ui)
		ui.mode = normalMode
	} else if r == 'e' && len(ui.currentList().items) != 0 {
		ui.enterEdit()
	} else if r == 'b' {
		ui.mode = normalMode
	}
}

func handleDeleteListModeEv(ui *UI, key tcell.Key, r rune) {
	if r == 'y' {
		ui.deleteList()
		ui.mode = normalMode
	} else if r == 'n' {
		ui.mode = normalMode
	}
}

func handleEditListNameModeEv(ui *UI, key tcell.Key, r rune) {
	if key == tcell.KeyEscape {
		ui.exitNameEdit()
	} else if r == 127 {
		ui.currentList().deleteRuneFromName()
	} else {
		ui.currentList().addRuneToName(r)
	}
}

func handleEditMode(ui *UI, key tcell.Key, r rune) {
	list := ui.currentList()
	if key == tcell.KeyEscape {
		ui.exitEdit()
	} else if r == 127 {
		list.deleteRune()
	} else {
		list.addRune(r)
	}
}

func (ui *UI) render() {
	renderListNav(ui)
	renderCurrentList(ui)
	renderFooter(ui)
}

func renderCurrentList(ui *UI) {
	if len(ui.lists) == 0 {
		ui.renderLine("Press l + n to create a new list", headerHeight-4)
		return
	}
	ui.lists[ui.current].render(ui)
}

func renderListNav(ui *UI) {
	var style tcell.Style
	for i := range ui.lists {
		if i == ui.current {
			style = ui.styles.nav 
		} else {
			style = ui.styles.def
		}
		r := strconv.Itoa(i + 1)
		ui.screen.SetContent(i*2+leftOffset, 1, []rune(r)[0], nil, style)
	}
	if len(ui.lists) != 0 {
		ui.screen.SetContent(ui.current*2+leftOffset, 2, '^', nil, ui.styles.primary)
	} else {
        renderSeparator(ui, separator(ui, ""), 5)
    }
}

func separator(ui *UI, s string) string {
	w := ui.width()
	var line bytes.Buffer
	var spaceTaken int
	if w-2 < len(s) {
		spaceTaken = 0
	} else {
		spaceTaken = len(s)
		line.WriteString(s)
	}
	for i := 0; i < w-2-spaceTaken; i++ {
		line.WriteRune(' ')
	}
	return line.String()
}

func renderSeparator(ui *UI, line string, ypos int) {
	for col, r := range []rune(line) {
		ui.screen.SetContent(
			col+leftOffset,
			ypos,
			r,
			nil,
            ui.styles.rib)
	}
}

func (ui *UI) renderLine(line string, row int) {
	for col, r := range []rune(line) {
		ui.screen.SetContent(col+leftOffset, row+topOffset, r, nil, ui.styles.def)
	}
}

func renderFooter(ui *UI) {
	footerYPos := ui.height() - 2
	var line string
	if ui.mode == listMode {
		line = "List: (n)ew - (d)elete - (e)dit - (b)ack"
	} else if ui.mode == entryMode {
		line = "Entry: (n)ew - (d)elete - (e)dit - (b)ack"
	} else if ui.mode == deleteListMode {
		line = "Delete current list? y / n"
	} else if ui.mode == editListNameMode {
    line = "Editing list name... - (esc)ape"
  } else if ui.mode == editMode {
    line = "Editing entry... - (esc)ape"
  } else if len(ui.lists) != 0 {
		line = "(enter) mark - (l)ist mode"
    if len(ui.currentList().items) != 0 {
      line += " - (e)ntry mode"
    }
    line += " - e(x)it"
	}
	renderSeparator(ui, separator(ui, line), footerYPos)
}

func (ui *UI) enterEdit() {
	l := ui.currentList()
	l.items[l.row].content = " "
	ui.mode = editMode
}

func (ui *UI) exitEdit() {
	ui.mode = normalMode
	ui.currentList().updateItem(ui.db)
	ui.currentList().col = 0
}

func (ui *UI) enterNameEdit() {
	ui.currentList().name = " "
	ui.mode = editListNameMode
}

func (ui *UI) exitNameEdit() {
	ui.mode = normalMode
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

func (ui *UI) exit() {
	ui.saveOrder()
	ui.screen.Fini()
	os.Exit(0)
}

func setStyles(ui *UI) *UI {
	dstyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	hstyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
	estyle := tcell.StyleDefault.Background(tcell.ColorLightSkyBlue).Foreground(tcell.ColorBlack)
	pstyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorLightSkyBlue)
	successStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorSeaGreen)
	errorStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorPaleVioletRed)
    ribStyle := tcell.StyleDefault.Background(tcell.ColorSeaGreen).Foreground(tcell.ColorBlack)
    navStyle := tcell.StyleDefault.Background(tcell.ColorLightSkyBlue).Foreground(tcell.ColorBlack)
	styles := Styles{
		edit:      estyle,
		highlight: hstyle,
		primary:   pstyle,
		def:       dstyle,
		success:   successStyle,
		error:     errorStyle,
        rib:       ribStyle,
        nav:       navStyle,
	}
	ui.screen.SetStyle(dstyle)
	ui.styles = styles
	return ui
}
