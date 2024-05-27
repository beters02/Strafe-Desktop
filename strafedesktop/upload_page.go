package main

import (
	"fmt"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
)

func uploadFiles(mainapp *Application, s *FileSelector, tree *fyne.Container) {

	// lets find out where these files are going
	f := s.GetFiles().(map[string]string)

	ignored := []string{}
	println("Attempting file upload... ")
	for dir, loc := range f {
		if loc == "" {
			loc = "Misc"
		}

		fn := getFileName(dir)
		nd := "misc/uploads/" + fn

		dd := fixUid(dir)
		fmt.Println(dd)
		success := localFileCopy(dd, nd, false)
		if !success {
			ignored = append(ignored, dir)
			deselectButton(dir, s, tree)
		} else {
			mainapp.net.Upload(fn, loc)
			deselectButton(dir, s, tree)
			os.Remove(nd)
		}
	}

	if len(ignored) > 0 {
		str := "Could not upload some files: \n"
		for _, dir := range ignored {
			str = str + dir + "\n"
		}
		d := dialog.NewCustom("Failed", "Dismiss", widget.NewLabel(str), mainapp.window)
		d.Show()
	} else {
		d := dialog.NewCustom("Success", "Dismiss", widget.NewLabel("Files uploaded to strafe!"), mainapp.window)
		d.Show()
	}
}

func createUploadPage(mainapp *Application) *Page {
	s := NewFileSelector(false)
	selectFromTree := xwidget.NewFileTree(storage.NewFileURI(getLocalHomeDir()))
	selectedTree := container.NewGridWithRows(10)

	selectFromTree.OnSelected = func(uid string) {
		exists := s.GetFile(uid)
		if exists {
			deselectButton(uid, s, selectedTree)
		} else {
			selectButton(uid, mainapp, s, selectedTree, true)
		}
	}

	frombox := container.NewBorder(canvas.NewText("Select files to upload.", color.White), nil, nil, nil, selectFromTree)
	tobox := container.NewBorder(canvas.NewText("Selected files. 10 maximum.", color.White), nil, nil, nil, selectedTree)
	backbut := widget.NewButton("Back", func() { SetPage("Home", mainapp.window) })
	upbut := widget.NewButton("Upload", func() {})
	mc := container.NewWithoutLayout(frombox, tobox, backbut, upbut)

	// config stuff
	screenSize := V2{x: 1024, y: 768}
	var padding float32 = 20

	butSizeV := V2{x: 100, y: 40}
	backbutPosV := V2{x: screenSize.x - (butSizeV.x + padding), y: screenSize.y - (butSizeV.y + padding)}
	upbutPosV := V2{x: backbutPosV.x, y: backbutPosV.y - butSizeV.y - (padding / 2)}

	fromSize := fyne.NewSize(screenSize.x/2, screenSize.y)
	fromPos := fyne.NewPos(0, 0)

	toSize := fyne.NewSize(screenSize.x/2, screenSize.y-(butSizeV.y*2)-(padding*2))
	toPos := fyne.NewPos(512, 0)

	frombox.Resize(fromSize)
	frombox.Move(fromPos)
	tobox.Resize(toSize)
	tobox.Move(toPos)
	backbut.Resize(fyne.NewSize(butSizeV.x, butSizeV.y))
	backbut.Move(fyne.NewPos(backbutPosV.x, backbutPosV.y))

	upbut.Resize(fyne.NewSize(butSizeV.x, butSizeV.y))
	upbut.Move(fyne.NewPos(upbutPosV.x, upbutPosV.y))
	upbut.OnTapped = func() {
		uploadFiles(mainapp, s, selectedTree)
	}

	page := NewPageClass("Upload", mainapp, mc)
	return page
}

//
