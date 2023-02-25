package main

import "github.com/gdamore/tcell"

var dark tcell.Color = tcell.NewRGBColor(26, 27, 44)
var light tcell.Color = tcell.NewRGBColor(192, 202, 245)
var primary tcell.Color = tcell.NewRGBColor(40, 52, 74)
var secondary tcell.Color = tcell.NewRGBColor(115, 218, 202)
var tertiary tcell.Color = tcell.NewRGBColor(247, 118, 142)

// bgcolorFgcolor
var darkLight = createStyle(dark, light)
var lightDark = darkLight.Reverse(true)

var darkPrimary = createStyle(dark, primary)
var darkSecondary = createStyle(dark, secondary)
var darkTertiary = createStyle(dark, tertiary)
var primaryDark = darkPrimary.Reverse(true)
var secondaryDark = darkSecondary.Reverse(true)
var tertiaryDark = darkTertiary.Reverse(true)

var lightPrimary = createStyle(light, primary)
var lightSecondary = createStyle(light, secondary)
var lightTertiary = createStyle(light, tertiary)
var primaryLight = lightPrimary.Reverse(true)
var secondaryLight = lightSecondary.Reverse(true)
var tertiaryLight = lightTertiary.Reverse(true)

func createStyle(bg, fg tcell.Color) tcell.Style {
	return tcell.StyleDefault.Background(bg).Foreground(fg)
}
