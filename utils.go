package main

import (
	"fmt"
	"log"
	"os"
)

func logToFile(s string) {
	f, err := os.OpenFile("text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(s + "\n"); err != nil {
		log.Println(err)
	}
}

func uiValsDebugPrint(ui *UI) {
	if len(ui.lists) == 0 {
		return
	}
	w, h := ui.screen.Size()
	line := fmt.Sprintf("Width: %d, Height: %d, Wtop: %d, Wbottom: %d, row: %d", w, h, ui.windowTop, ui.windowBottom, ui.currentList().row)
	for col, r := range []rune(line) {
		ui.screen.SetContent(col+1, 0, r, nil, ui.dstyle)
	}
}

func max(val1, val2 int) int {
	if val1 >= val2 {
		return val1
	}
	return val2
}
