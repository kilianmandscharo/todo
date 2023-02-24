package main

import "github.com/gdamore/tcell"

var dark tcell.Color = tcell.NewRGBColor(26, 27, 44)
var light tcell.Color = tcell.NewRGBColor(192, 202, 245)
var primary tcell.Color = tcell.NewRGBColor(40, 52, 74)
var secondary tcell.Color = tcell.NewRGBColor(115, 218, 202)
var tertiary tcell.Color = tcell.NewRGBColor(247, 118, 142)

// bgcolorFgcolor
var darkPrimary = createStyle(dark, primary)
var darkLight = createStyle(dark, light)
var lightDark = createStyle(light, dark)
var primaryLight = createStyle(primary, light)
var secondaryLight = createStyle(secondary, light)
var secondaryDark = createStyle(secondary, dark)
var tertiaryLight = createStyle(tertiary, light)
var tertiaryDark = createStyle(tertiary, dark)

func createStyle(bg, fg tcell.Color) tcell.Style {
	return tcell.StyleDefault.Background(bg).Foreground(fg)
}
