package main

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"unicode"

	"github.com/gdamore/tcell"
)

type Mode int

const (
	normalMode Mode = iota
	editMode
	editListNameMode
	deleteListMode
)

const headerHeight = 6
const footerHeight = 2
const leftOffset = 1
const topOffset = 1
const bottomOffset = 1
const navPosition = 1

var modeTitleMap = map[Mode]string{
	normalMode:       "Normal",
	editListNameMode: "Insert",
	deleteListMode:   "Delete",
	editMode:         "Insert",
}

var modeStyleMap = map[Mode]tcell.Style{
	normalMode:       lightDark,
	editListNameMode: secondaryLight,
	deleteListMode:   tertiaryLight,
	editMode:         secondaryLight,
}

type UI struct {
	screen       tcell.Screen
	db           *DB
	lists        []List
	current      int
	windowTop    int
	windowBottom int
	mode         Mode
}

func newUI(debug bool) *UI {
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
	ui.mode = normalMode
	ui.screen.SetStyle(darkLight)
	return ui
}

func (ui *UI) load() {
	lists, err := ui.db.getLists()
	if err != nil {
		log.Fatal(err)
	}
	ui.loadOrder(lists)
	ui.calculateWindow()
}

func (ui *UI) clear() {
	ui.screen.Clear()
}

func (ui *UI) show() {
	ui.screen.Show()
}

func (ui *UI) currentList() *List {
	if len(ui.lists) == 0 {
		return nil
	}
	return &ui.lists[ui.current]
}

func (ui *UI) switchList(r rune) {
	if !unicode.IsDigit(r) || r == '0' {
		return
	}
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
    ui.current = len(ui.lists) - 1
    ui.calculateWindow()
    ui.mode = editListNameMode
    ui.currentList().name = " "
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

func (ui *UI) loadOrder(lists []List) error {
	order, err := ui.db.loadListOrder()
	if err != nil {
		return err
	}
	if order == nil || len(order) == 0 || len(order) != len(lists) {
		ui.lists = lists
	} else {
		var orderedLists []List
		for i := 0; i < len(order); i++ {
			for j := range lists {
				if lists[j].ID == order[i] {
					orderedLists = append(orderedLists, lists[j])
					continue
				}
			}
		}
		ui.lists = orderedLists
	}

	for i := range ui.lists {
		order, err := ui.db.loadItemOrder(ui.lists[i].ID)
		if err != nil {
			return err
		}
		if order == nil || len(order) == 0 || len(order) != len(ui.lists[i].items) {
			continue
		}
		var orderedItems []Item
		for j := 0; j < len(order); j++ {
			for _, item := range ui.lists[i].items {
				if item.id == order[j] {
					orderedItems = append(orderedItems, item)
				}
			}
		}
		ui.lists[i].items = orderedItems
	}
	return nil
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
		case editListNameMode:
			handleEditListNameModeEv(ui, ev.Key(), ev.Rune())
		case deleteListMode:
			handleDeleteListModeEv(ui, ev.Key(), ev.Rune())
		case editMode:
			handleEditModeEv(ui, ev.Key(), ev.Rune())
		}
	}
}

func handleNormalModeEv(ui *UI, key tcell.Key, r rune) {
	if r == 'x' {
		ui.exit()
	} else if r == 'h' {
		ui.left()
	} else if r == 'l' {
		ui.right()
	} else if r == 'H' {
		ui.switchListLeft()
	} else if r == 'L' {
		ui.switchListRight()
	} else if r == 'D' {
		ui.enterDeleteListMode()
	} else if r == 'I' {
		ui.enterNameEdit()
	} else if r == 'N' {
		ui.addList()
	} else if r == 'j' {
		ui.listDown()
	} else if r == 'k' {
		ui.listUp()
	} else if r == 'J' {
		ui.listSwitchDown()
	} else if r == 'K' {
		ui.listSwitchUp()
	} else if r == 'n' {
		ui.listAddEntry()
	} else if r == 'd' {
		ui.listDeleteEntry()
	} else if r == 'i' {
		ui.enterEdit()
	} else if r == 13 {
		ui.listMarkEntry()
	} else {
		ui.switchList(r)
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
	list := ui.currentList()
	if key == tcell.KeyEscape || key == tcell.KeyEnter {
		ui.exitNameEdit()
	} else if r == 127 {
		ui.currentList().deleteRuneFromName()
	} else if key == tcell.KeyLeft {
		list.cursorLeftListName()
	} else if key == tcell.KeyRight {
		list.cursorRightListName()
	} else {
		ui.currentList().addRuneToName(r)
	}
}

func handleEditModeEv(ui *UI, key tcell.Key, r rune) {
	list := ui.currentList()
	if key == tcell.KeyEscape || key == tcell.KeyEnter {
		ui.exitEdit()
	} else if r == 127 {
		list.deleteRune()
	} else if key == tcell.KeyLeft {
		list.cursorLeftEntry()
	} else if key == tcell.KeyRight {
		list.cursorRightEntry()
	} else {
		list.addRune(r)
	}
}

func (ui *UI) listDeleteEntry() {
	if list := ui.currentList(); list != nil {
		list.delete(ui.db, ui)
	}
}

func (ui *UI) listMarkEntry() {
	if list := ui.currentList(); list != nil {
		list.markItem(ui.db)
	}
}

func (ui *UI) listAddEntry() {
	if list := ui.currentList(); list != nil {
		list.add(ui.db, ui)
		list.down(ui)
		ui.mode = editMode
	}
}

func (ui *UI) listDown() {
	if list := ui.currentList(); list != nil {
		list.down(ui)
	}
}

func (ui *UI) listUp() {
	if list := ui.currentList(); list != nil {
		list.up(ui)
	}
}

func (ui *UI) listSwitchDown() {
	if list := ui.currentList(); list != nil {
		list.switchDown(ui)
	}
}

func (ui *UI) listSwitchUp() {
	if list := ui.currentList(); list != nil {
		list.switchUp(ui)
	}
}

func (ui *UI) enterDeleteListMode() {
	if len(ui.lists) != 0 {
		ui.mode = deleteListMode
	}
}

func (ui *UI) render() {
	renderListNav(ui)
	renderCurrentList(ui)
	renderFooter(ui)
	// uiValsDebugPrint(ui)
}

func renderCurrentList(ui *UI) {
	if len(ui.lists) == 0 {
		ui.renderLine("Press N to create a new list", headerHeight-4)
		return
	}
	ui.currentList().render(ui)
}

func renderListNav(ui *UI) {
	var style tcell.Style
	for i := range ui.lists {
		if i == ui.current {
			style = primaryLight
		} else {
			style = darkLight
		}
		r := strconv.Itoa(i + 1)
		ui.screen.SetContent(i*2+leftOffset, 1, []rune(r)[0], nil, style)
	}
	if len(ui.lists) != 0 {
		ui.screen.SetContent(ui.current*2+leftOffset, 2, '^', nil, darkPrimary)
	} else {
		renderTopSeparator(ui, separator(ui, "", 0), 5)
	}
}

func separator(ui *UI, s string, offset int) string {
	w := ui.width() - offset
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

func renderTopSeparator(ui *UI, line string, ypos int) {
	for col, r := range []rune(line) {
		ui.screen.SetContent(
			col+leftOffset,
			ypos,
			r,
			nil,
			primaryLight)
	}
}

func (ui *UI) renderLine(line string, row int) {
	for col, r := range []rune(line) {
		ui.screen.SetContent(col+leftOffset, row+topOffset, r, nil, darkLight)
	}
}

func renderFooter(ui *UI) {
	footerYPos := ui.height() - 2
	var line string
	if ui.mode == deleteListMode {
		line = "Delete current list? y / n"
	} else if ui.mode == editListNameMode {
		line = "List name - (esc)ape"
	} else if ui.mode == editMode {
		line = "Entry name - (esc)ape"
	} else if len(ui.lists) != 0 {
		line = "(enter) mark"
		line += " - e(x)it"
	}
	modeString := padChunk(modeTitleMap[ui.mode])
	line = separator(ui, padChunk(line), len(modeString))
	renderChunk(ui, modeString, modeStyleMap[ui.mode], 0, footerYPos)
	renderChunk(ui, line, primaryLight, len(modeString), footerYPos)
}

func renderChunk(ui *UI, s string, style tcell.Style, col, row int) {
	for i, r := range []rune(s) {
		ui.screen.SetContent(col+i+leftOffset, row, r, nil, style)
	}
}

func padChunk(s string) string {
	return " " + s + " "
}

func (ui *UI) enterEdit() {
	if l := ui.currentList(); l != nil && len(l.items) != 0 {
		itemLength := len(l.currentItem().content)
		if itemLength == 1 && l.currentItem().content[0] == ' ' {
			l.col = 0
		} else {
			l.col = itemLength
		}
		ui.mode = editMode
	}
}

func (ui *UI) exitEdit() {
	ui.mode = normalMode
	ui.currentList().updateItem(ui.db)
	ui.currentList().col = 0
}

func (ui *UI) enterNameEdit() {
	if len(ui.lists) == 0 {
		return
	}
	l := ui.currentList()
	if len(l.name) == 1 && l.name[0] == ' ' {
		l.col = 0
	} else {
		l.col = len(l.name)
	}
	ui.mode = editListNameMode
}

func (ui *UI) exitNameEdit() {
	ui.mode = normalMode
	ui.currentList().updateName(ui.db)
	ui.currentList().col = 0
}

func (ui *UI) listSpaceAvailable() int {
	return ui.height() - topOffset - headerHeight - bottomOffset - footerHeight
}

func (ui *UI) calculateWindow() {
	if len(ui.lists) == 0 {
		return
	}
	_, h := ui.screen.Size()
	topOffset := topOffset + headerHeight
	bottomOffset := bottomOffset + footerHeight
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

func (ui *UI) left() {
	if ui.current > 0 {
		ui.current--
		ui.calculateWindow()
	}
}

func (ui *UI) right() {
	if ui.current < len(ui.lists)-1 {
		ui.current++
		ui.calculateWindow()
	}
}

func (ui *UI) switchListLeft() {
	if ui.current == 0 {
		return
	}
	ui.lists[ui.current], ui.lists[ui.current-1] = ui.lists[ui.current-1], ui.lists[ui.current]
	ui.current--
}

func (ui *UI) switchListRight() {
	if ui.current == len(ui.lists)-1 {
		return
	}
	ui.lists[ui.current], ui.lists[ui.current+1] = ui.lists[ui.current+1], ui.lists[ui.current]
	ui.current++
}
