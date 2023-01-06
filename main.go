package main

import (
	"log"
	"os"

	"github.com/gdamore/tcell"
)

func main() {
    s, err := tcell.NewScreen()
    if err != nil {
        log.Fatal(err)
    }
    if err := s.Init(); err != nil {
        log.Fatal(err)
    }

    current := 0

    defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
    hiStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
    s.SetStyle(defStyle)

    s.Clear()

    todos := []string{
        "Clean bedroom",
        "Practice piano",
        "Prepare dinner",
    }

    renderTodos(current, todos, s, defStyle, hiStyle)   

    for {
        s.Show()
        ev := s.PollEvent()

        switch ev := ev.(type) {
        case *tcell.EventResize:
            s.Sync()
        case *tcell.EventKey: 
            if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
                s.Fini()
                os.Exit(0)
            }
            if ev.Rune() == 'j' {
                if current + 1 <= len(todos) - 1 {
                    current++
                    renderTodos(current, todos, s, defStyle, hiStyle)
                }             
            } 
            if ev.Rune() == 'k' {
                if current - 1 >= 0 {
                    current--
                    renderTodos(current, todos, s, defStyle, hiStyle)
                }             
            } 
        
        }
    }
}

func deleteTodo(todos []string, index int) {

}

func renderTodos(current int, todos []string, s tcell.Screen, dstyle tcell.Style, hstyle tcell.Style) {
    for row, todo := range todos {
        if row == current {
            renderLine(s, hstyle, row, todo)
        } else {
            renderLine(s, dstyle, row, todo)
        }
    }

}

func renderLine(s tcell.Screen, style tcell.Style, row int, line string) {
    for col, r := range []rune(line) {
        s.SetContent(col, row, r, nil, style)
    }
}
