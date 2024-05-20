package main

import (
	"image/color"
	"io/fs"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	strafe "github.com/beters02/Strafe"
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

func getFilesAt(net strafe.Net, path string) []fs.FileInfo {
	prePath := "/shr/strafe"
	fi, err := net.Client.ReadDir(prePath + path)
	if err != nil {
		panic(err)
	}
	return fi
}

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
	files := getFilesAt(net, "/")
	fileGrid := container.NewAdaptiveGrid(len(files))

	// finish pages

	// nice! i forgot the upload page needs to display a dialog box or something.
	// the file dialog for fyne is actually pretty bad... gotta make a custom file explorer for this.
	// ugh.
	ufc := container.NewBorder(nil, container.NewVBox(fileGrid, widget.NewLabel("")), nil, nil, uhead)

	// this guy is fine though
	dfc := container.NewBorder(nil, container.NewVBox(fileGrid, widget.NewLabel("")), nil, nil, dhead)

	setPage := func(page string) {
		if page == "upload" {
			w.Resize(fyne.NewSize(1024, 768))
			currentPage = "upload"
			w.SetContent(ufc)
			w.CenterOnScreen()
		} else {
			w.Resize(fyne.NewSize(1024, 768))
			currentPage = "download"
			w.SetContent(dfc)
			w.CenterOnScreen()
		}
	}

	// create home page
	hc := homePage(setPage)

	// add all weapon buttons
	for _, file := range files {
		fileGrid.Add(widget.NewButton(file.Name(), func() {
			if currentPage == "none" {
				return
			} else if currentPage == "upload" {
				println("Print all weapons from local")
			} else {
				println("Print all weapons from ssh")
			}
		}))
	}

	fileGrid.Add(widget.NewButton("Back", func() {
		w.SetContent(hc)
		w.Resize(fyne.NewSize(400, 300))
		w.CenterOnScreen()
	}))

	net.Disconnect()

	// lets go baby
	w.SetContent(hc)
	w.ShowAndRun()
}
