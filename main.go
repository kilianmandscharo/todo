package main

import (
	"os"
)

func main() {
	var debug bool
	if len(os.Args) > 1 {
		debug = os.Args[1] == "debug"
	}

	ui := newUI(debug)
	ui.load()
	defer ui.closeDB()

	if debug {
		os.Exit(0)
	}

	for {
		ui.clear()
		ui.render()
		ui.show()
		ev := ui.event()
		ui.handleEvent(ev)
	}
}
