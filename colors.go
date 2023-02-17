package main

import "github.com/gdamore/tcell"

var dark tcell.Color = tcell.NewRGBColor(26, 27, 38)
var light tcell.Color = tcell.NewRGBColor(240, 240, 240)
var primary tcell.Color = tcell.NewRGBColor(23, 111, 160)
var secondary tcell.Color = tcell.NewRGBColor(63, 148, 100)
var tertiary tcell.Color = tcell.ColorOrangeRed

// bgcolorFgcolor
var darkPrimary = createStyle(dark, primary)
var darkLight = createStyle(dark, light)
var lightDark = createStyle(light, dark)
var primaryLight = createStyle(primary, light)
var secondaryLight = createStyle(secondary, light)
var tertiaryLight = createStyle(tertiary, light)

func createStyle(bg, fg tcell.Color) tcell.Style {
	return tcell.StyleDefault.Background(bg).Foreground(fg)
}
