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
  db, err := newDatabase("data.db")
  if err != nil {
    log.Fatal(err)
  }
  initDatabase(db)
  defer db.Close()

	ui := newUI()
	ui.addList()

	for {
		ui.clear()
		ui.currentList().render()
		ui.show()
		ev := ui.screen.PollEvent()
		ui.handleEvent(ev)
	}
}
