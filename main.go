package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	strafe "github.com/beters02/Strafe"
)

func homePage(pressed func(page string)) *fyne.Container {
	strafeHeader := createHeader("Strafe")
	info := createTextLine("GUI Wrapper for CS Labs' file manager.")
	text := container.NewGridWithRows(2, strafeHeader, info)
	txtc := container.NewCenter(text)

	uploadbutton := widget.NewButton("Upload", func() {
		pressed("upload")
	})
	downloadbutton := widget.NewButton("Download", func() {
		pressed("download")
	})
	buttons := container.NewHBox(uploadbutton, downloadbutton)
	bc := container.NewCenter(buttons)

	all := container.NewVBox(txtc, bc)
	return container.NewCenter(all)
}

func main() {

	// initialize app
	a := app.New()
	w := a.NewWindow("Strafe")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(400, 300)) //w.Resize(fyne.NewSize(1024, 768))

	// connect to the asset library
	net := strafe.NetConnect()

	// init page headers
	uhead := createHeader("Upload Assets")
	dhead := createHeader("Download Assets")

	// init page switching functionality
	currentPage := "none"

	// create container with buttons for all weapons, since thats all we are storing right now
	files := getServerFilesAt(net, "/")
	fileGrid := container.NewAdaptiveGrid(len(files))

	// finish pages

	// nice! i forgot the upload page needs to display a dialog box or something.
	// the file dialog for fyne is actually pretty bad... gotta make a custom file explorer for this.
	// ugh.
	uploadPage := container.NewBorder(nil, container.NewVBox(fileGrid, widget.NewLabel("")), nil, nil, uhead)

	// this guy is fine though
	downloadPage := container.NewBorder(nil, container.NewVBox(fileGrid, widget.NewLabel("")), nil, nil, dhead)

	// create home page
	hc := homePage(func(page string) {
		if page == "upload" {
			w.Resize(fyne.NewSize(1024, 768))
			currentPage = "upload"
			w.SetContent(uploadPage)
			w.CenterOnScreen()
		} else {
			w.Resize(fyne.NewSize(1024, 768))
			currentPage = "download"
			w.SetContent(downloadPage)
			w.CenterOnScreen()
		}
	})

	// add all weapon buttons
	for _, file := range files {
		fileGrid.Add(widget.NewButton(file.Name(), func() {
			if currentPage == "none" {
				return
			} else if currentPage == "upload" {
				println("Print all weapons from local")
				println("JK This shouldnt be happening at all")
			} else {
				println("Print all weapons from ssh")
			}
		}))
	}

	// cant forget this bad boy
	fileGrid.Add(widget.NewButton("Back", func() {
		w.SetContent(hc)
		w.Resize(fyne.NewSize(400, 300))
		w.CenterOnScreen()
	}))

	// cant forget this bad boy
	net.Disconnect()

	// lets go baby
	w.SetContent(hc)
	w.ShowAndRun()
}
