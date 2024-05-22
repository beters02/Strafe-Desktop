package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func newTextLine(text string) *canvas.Text {
	info := canvas.NewText(text, color.White)
	info.TextStyle = fyne.TextStyle{Italic: true}
	info.Alignment = fyne.TextAlignCenter
	return info
}

func newHeader(text string) *canvas.Text {
	strafe := canvas.NewText(text, color.White)
	strafe.TextSize = 50
	strafe.Alignment = fyne.TextAlignCenter
	return strafe
}

func newDropdownSelection(text string, items []string, callback func(string)) *fyne.Container {
	l := canvas.NewText(text, color.White)
	s := widget.NewSelect(items, callback)
	return container.NewVBox(l, s)
}
