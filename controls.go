package main

import "fmt"

type Keymap struct {
	key  string
	text string
}

func newKeyMap(key string, text string) Keymap {
	return Keymap{key: key, text: text}
}

var keymaps = []Keymap{
	newKeyMap("n", "create new entry"),
	newKeyMap("N", "create new list"),
	newKeyMap("", ""),
	newKeyMap("d", "delete entry"),
	newKeyMap("D", "delete list"),
	newKeyMap("", ""),
	newKeyMap("i", "edit entry"),
	newKeyMap("I", "edit list name"),
	newKeyMap("", ""),
	newKeyMap("j", "go one entry down"),
	newKeyMap("k", "go one entry up"),
	newKeyMap("J", "switch entry with the one belowe"),
	newKeyMap("K", "switch entry with the on above"),
	newKeyMap("", ""),
	newKeyMap("h", "go one list to the left"),
	newKeyMap("l", "go one list to the right"),
	newKeyMap("H", "switch list with the one to the left"),
	newKeyMap("L", "switch list with the one to the right"),
	newKeyMap("", ""),
	newKeyMap("0-9", "switch to list (0-9)"),
	newKeyMap("enter", "toggle entry"),
	newKeyMap("x", "exit"),
}

func printKeymaps() {
	for _, k := range keymaps {
		if len(k.key) != 0 {
			fmt.Printf("%-6s %s\n", k.key, k.text)
		} else {
			fmt.Printf("\n")
    }
	}
}
