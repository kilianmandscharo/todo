package main

import (
	"flag"
	"os"
)

var cFlag = flag.Bool("controls", false, "set to print controls overview")

func main() {
	flag.Parse()

	if *cFlag {
		printKeymaps()
		return
	}

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
