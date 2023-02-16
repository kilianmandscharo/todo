package main

import "github.com/gdamore/tcell"

var dark tcell.Color = tcell.ColorBlack
var light tcell.Color = tcell.ColorWhite
var primary tcell.Color = tcell.ColorSeaGreen
var secondary tcell.Color = tcell.ColorLightSkyBlue

// bgcolorFgcolor
var darkPrimary = createStyle(dark, primary)
var darkSecondary = createStyle(dark, secondary)
var darkLight = createStyle(dark, light)
var lightDark = createStyle(light, dark)
var primaryDark = createStyle(primary, dark)
var secondaryDark = createStyle(secondary, dark)

func createStyle(bg, fg tcell.Color) tcell.Style {
	return tcell.StyleDefault.Background(bg).Foreground(fg)
}
