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

func homePage(window fyne.Window, uploadPage *fyne.Container) *fyne.Container {
	strafeHeader := createHeader("Strafe")
	info := createTextLine("GUI Wrapper for CS Labs' file manager.")
	text := container.NewGridWithRows(2, strafeHeader, info)
	txtc := container.NewCenter(text)

	uploadbutton := widget.NewButton("Upload", func() {
		window.SetContent(uploadPage)
		window.CenterOnScreen()
	})
	downloadbutton := widget.NewButton("Download", func() {})
	buttons := container.NewHBox(uploadbutton, downloadbutton)
	bc := container.NewCenter(buttons)

	all := container.NewVBox(txtc, bc)
	return container.NewCenter(all)
}

// upload/download page
func filePage(net strafe.Net) *fyne.Container {
	label := createHeader("Upload Assets")
	files := getFilesAt(net, "/")
	fileGrid := container.NewAdaptiveGrid(len(files))

	for _, file := range files {
		fileGrid.Add(widget.NewButton(file.Name(), func() {}))
	}

	fc := container.NewVBox(fileGrid, widget.NewLabel(""))
	return container.NewBorder(nil, fc, nil, nil, label) //fc //container.NewWithoutLayout(lc, fc)
}

func main() {
	a := app.New()
	w := a.NewWindow("Strafe")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(400, 300)) //w.Resize(fyne.NewSize(1024, 768))

	net := strafe.NetConnect()

	// temp create file page first
	uc := filePage(net)
	hc := homePage(w, uc)

	net.Disconnect()

	w.SetContent(hc)
	w.ShowAndRun()
}
