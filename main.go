package main

import (
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

func main() {
	var debug bool
	if len(os.Args) > 1 {
		debug = os.Args[1] == "debug"
	}
	ui := newUI(debug)
	ui.load()
	ui.calculateWindow()

	defer ui.db.close()

	if debug {
		os.Exit(0)
	}

	for {
		ui.clear()
    ui.debugPrint()
		ui.render()
		ui.show()
		ev := ui.screen.PollEvent()
		ui.handleEvent(ev)
	}
}
