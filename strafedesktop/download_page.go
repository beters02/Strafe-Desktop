package main

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
)

func downloadFolderAsZip(mainapp *Application, folderPath string) {
	pathArr := []string{}
	folderName := getFileName(folderPath)

	for _, wepFile := range getServerFilesAt(mainapp.net, "/"+folderName) {
		path := mainapp.net.Download(wepFile.Name(), folderName)
		pathArr = append(pathArr, path)
	}

	os.Mkdir("temp", 0700)
	for _, path := range pathArr {
		fn := getFileName(path)
		localFileCopy(path, "temp/"+fn, true)
	}

	ZipWriter("temp/", "misc/downloads/"+folderName+".zip")
	err := os.RemoveAll("temp")
	fmt.Println(err)
}

func downloadFiles(mainapp *Application, s *FileSelector, tree *fyne.Container) {
	f := s.GetFiles().(map[string]string)
	if len(f) == 0 {
		return
	}

	for dir := range f {
		deselectButton(dir, s, tree)

		ft := getFileType(dir)
		if ft == "folder" {
			downloadFolderAsZip(mainapp, dir)
			continue
		}

		fn := getFileName(dir)
		folderSlash := secToLastIndex(dir, "/")
		nameSlash := strings.LastIndex(dir, "/")
		folder := substr(dir, folderSlash, nameSlash-folderSlash)
		folder = getFileName(folder)

		fmt.Printf("Downloading %v from %v\n", fn, folder)
		mainapp.net.Download(fn, folder)
	}

	d := dialog.NewCustom("Success", "Dismiss", widget.NewLabel("Files downloaded!"), mainapp.window)
	d.Show()
}

func createDownloadPage(mainapp *Application) *Page {

	// create container with buttons for all weapons, since thats all we are storing right now
	weaponFiles := getServerFilesAt(mainapp.net, "/")

	// create a file cache
	if isLocalFileExists("misc/cache") {
		os.Remove("misc/cache")
	}
	os.Mkdir("misc/cache", 0700)

	// populate cache
	for _, wepFile := range weaponFiles {
		wn := wepFile.Name()
		os.Mkdir("misc/cache/"+wn, 0700)
		for _, file := range getServerFilesAt(mainapp.net, "/"+wn) {
			os.Create("misc/cache/" + wn + "/" + file.Name())
		}
	}

	// now we create the container
	s := NewFileSelector(true)
	selectFromTree := xwidget.NewFileTree(storage.NewFileURI("misc/cache"))
	selectedTree := container.NewGridWithRows(10)

	selectFromTree.OnSelected = func(uid string) {
		exists := s.GetFile(uid)
		if exists {
			deselectButton(uid, s, selectedTree)
		} else {
			selectButton(uid, mainapp, s, selectedTree, false)
		}
	}

	frombox := container.NewBorder(canvas.NewText("Select files to upload.", color.White), nil, nil, nil, selectFromTree)
	tobox := container.NewBorder(canvas.NewText("Selected files. 10 maximum.", color.White), nil, nil, nil, selectedTree)
	backbut := widget.NewButton("Back", func() { SetPage("Home", mainapp.window) })
	upbut := widget.NewButton("Download", func() {})
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
		downloadFiles(mainapp, s, selectedTree)
	}

	return NewPageClass("Download", mainapp, mc)
}
