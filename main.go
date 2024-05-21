package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	strafe "github.com/beters02/Strafe"
)

type Application struct {
	app fyne.App
	net strafe.Net

	// this is the main window. i could see this being a problem at some point,
	// if i wanted to use multiple windows or something
	window fyne.Window
}

func main() {

	// initialize app
	a := app.New()
	mainapp := Application{
		app:    a,
		net:    strafe.NetConnect(),
		window: a.NewWindow("Strafe"),
	}

	mainapp.window.CenterOnScreen()
	mainapp.window.Resize(fyne.NewSize(400, 300)) //w.Resize(fyne.NewSize(1024, 768))

	// create pages
	InitPages(&mainapp)

	// lets go baby
	SetPage("Home", mainapp.window)
	mainapp.window.ShowAndRun()

	// cant forget this bad boy
	mainapp.net.Disconnect()
}
