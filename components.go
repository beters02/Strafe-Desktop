package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func createTextLine(text string) *canvas.Text {
	info := canvas.NewText(text, color.White)
	info.TextStyle = fyne.TextStyle{Italic: true}
	info.Alignment = fyne.TextAlignCenter
	return info
}

func createHeader(text string) *canvas.Text {
	strafe := canvas.NewText(text, color.White)
	strafe.TextSize = 50
	strafe.Alignment = fyne.TextAlignCenter
	return strafe
}
