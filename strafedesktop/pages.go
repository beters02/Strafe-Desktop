package main

import (
	"fmt"
	"image/color"
	"io/fs"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
)

type V2 struct {
	x float32
	y float32
}

// Page
type PageContainers map[string]*fyne.Container
type PageButtons map[string]*widget.Button
type PageObjects map[string]any // Unorganized page objects.
type Page struct {
	name       string
	mainapp    *Application
	containers PageContainers
	buttons    PageButtons
	objects    PageObjects
}

func page_open(page *Page) {
	page.mainapp.window.SetContent(page.containers["main"])
	if page.name == "Home" {
		page.mainapp.window.Resize(fyne.NewSize(400, 300))
	} else {
		page.mainapp.window.Resize(fyne.NewSize(1024, 768))
	}
	page.mainapp.window.CenterOnScreen()
}

func page_close(page *Page) {}

func (page *Page) Open() {
	page_open(page)
}

func (page *Page) Close() {
	page_close(page)
}

var Pages map[string]*Page // init main pages map
var CurrentPage string
var WeaponGrids map[string]*fyne.Container
var TreeArr []*widget.Button

// The correct way to create a page.
func NewPageClass(pagename string, mainapp *Application, maincontainer *fyne.Container) *Page {
	Pages[pagename] = &Page{
		name:    pagename,
		mainapp: mainapp,
		containers: PageContainers{
			"main": maincontainer,
		},
		buttons: PageButtons{},
		objects: PageObjects{},
	}
	return Pages[pagename]
}

func InitPages(ma *Application) {
	Pages = map[string]*Page{}
	CurrentPage = "Home"
	WeaponGrids = map[string]*fyne.Container{}
	TreeArr = []*widget.Button{}

	createHomePage(ma)
	createDownloadPage(ma)
	createUploadPage(ma)
}

func GetPage(page string) *Page {
	return Pages[page]
}

func SetPage(pageName string, w fyne.Window) {
	page := GetPage(pageName)
	currpage := GetPage(CurrentPage)

	currpage.Close()
	page.Open()

	CurrentPage = page.name
}

//

//Selector

type FileSelector struct {
	currFiles  map[string]string
	currFilesN int
	maxFiles   int
}

func NewFileSelector() *FileSelector {
	return &FileSelector{
		currFiles:  map[string]string{},
		currFilesN: 0,
		maxFiles:   10,
	}
}

func (s *FileSelector) GetFile(dir string) bool {
	_, ok := s.currFiles[dir]
	return ok
}

func (s *FileSelector) GetFilesAmount() int {
	return s.currFilesN
}

func (s *FileSelector) GetFiles() interface{} {
	return s.currFiles
}

// Change to FileData, for now you can edit the value in the dict
func (s *FileSelector) SetFileString(file string, str string) {
	s.currFiles[file] = str
}

func (s *FileSelector) GetFileString(file string) string {
	return s.currFiles[file]
}

func (s *FileSelector) Select(dir string) (didSelect bool) {
	if s.GetFilesAmount() >= s.maxFiles || s.GetFile(dir) {
		return false
	}
	s.currFiles[dir] = ""
	s.currFilesN++
	return true
}

func (s *FileSelector) Deselect(dir string) (didDeselect bool) {
	if !s.GetFile(dir) {
		return false
	}
	delete(s.currFiles, dir)
	s.currFilesN--
	return true
}

//

func createHomePage(ma *Application) *Page {

	// create main container and init page
	mc := container.NewCenter()
	page := NewPageClass("Home", ma, mc)

	// label
	strafeHeader := newHeader("Strafe")
	info := newTextLine("GUI Wrapper for CS Labs' file manager.")
	txtc := container.NewCenter(container.NewGridWithRows(2, strafeHeader, info))

	// buttons
	uploadbutton := widget.NewButton("Upload", func() { SetPage("Upload", page.mainapp.window) })
	downloadbutton := widget.NewButton("Download", func() { SetPage("Download", page.mainapp.window) })
	buttons := container.NewHBox(uploadbutton, downloadbutton)

	// clean it all up
	mc.Add(container.NewVBox(txtc, container.NewCenter(buttons)))

	// we wont add any of the objects to the struct since we dont need them
	return page
}

//

func createWeaponGrid(mainapp *Application, fileName string, wfiles []fs.FileInfo) { // fileName is weaponName
	_, ok := WeaponGrids[fileName]
	if !ok {
		WeaponGrids[fileName] = container.NewVBox() //(len(wfiles))
	}
	WeaponGrids[fileName].Objects = []fyne.CanvasObject{}
	WeaponGrids[fileName].Add(widget.NewButton("Refresh", func() {
		createWeaponGrid(mainapp, fileName, getServerFilesAt(mainapp.net, "/"+fileName))
	}))
	for _, wf := range wfiles { // wf is fileName
		wfName := wf.Name()
		WeaponGrids[fileName].Add(widget.NewButton(wfName, func() {
			mainapp.net.Download(wfName, fileName)
			d := dialog.NewCustom("Success", "Dismiss", widget.NewLabel("File downloaded!"), mainapp.window)
			d.Show()
		}))
	}
	WeaponGrids[fileName].Add(widget.NewButton("Back", func() {
		SetPage("Download", mainapp.window)
	}))
}

func createDownloadPage(mainapp *Application) *Page {
	dhead := newHeader("Download Assets")

	// create container with buttons for all weapons, since thats all we are storing right now
	files := getServerFilesAt(mainapp.net, "/")
	fileGrid := container.NewAdaptiveGrid(len(files))
	borderc := container.NewBorder(nil, container.NewVBox(fileGrid, widget.NewLabel(""), widget.NewLabel("")), nil, nil, dhead)

	// add all weapon buttons
	for _, file := range files {
		fileName := file.Name()
		wfiles := getServerFilesAt(mainapp.net, "/"+fileName)
		createWeaponGrid(mainapp, fileName, wfiles)
		fileGrid.Add(widget.NewButton(fileName, func() {
			mainapp.window.SetContent(WeaponGrids[fileName])
		}))
	}

	// cant forget this bad boy
	but := widget.NewButton("Back", func() { SetPage("Home", mainapp.window) })
	butc := container.NewWithoutLayout(but)
	but.Resize(fyne.NewSize(100, 40))
	but.Move(fyne.NewPos((1024/2)-20, 768-40-20))
	mc := container.NewStack(butc, borderc)

	// we wont add any of the objects to the struct since we dont need them
	page := NewPageClass("Download", mainapp, mc)
	return page
}

//

// remove "file://"
func fixUid(uid string) string {
	return substr(uid, 7, len(uid))
}

func deselectButton(uid string, s *FileSelector, tree *fyne.Container) {
	didDeselect := s.Deselect(uid)
	if !didDeselect {
		fmt.Println("Could not deselect item " + uid)
		return
	}

	for i := range tree.Objects {
		if TreeArr[i].Text == uid {
			tree.Objects = append(tree.Objects[:i], tree.Objects[i+1:]...)
			TreeArr = append(TreeArr[:i], TreeArr[i+1:]...)
			break
		}
	}
}

func getPossibleOptions(mainapp *Application) []string {
	optionsFiles := getServerFilesAt(mainapp.net, "/")
	options := []string{}
	for _, opt := range optionsFiles {
		options = append(options, opt.Name())
	}
	return options
}

func selectButton(uid string, mainapp *Application, s *FileSelector, tree *fyne.Container) {
	didSelect := s.Select(uid)
	if !didSelect {
		return
	}

	path := fixUid(uid)
	isDir := isLocalFileDir(path)

	if isDir {
		s.Deselect(uid)
		return
	}

	o := getPossibleOptions(mainapp)
	gc := newDropdownSelection("", o, func(str string) { s.SetFileString(uid, str) })
	b := widget.NewButton(uid, func() { deselectButton(uid, s, tree) })
	mc := container.NewHBox(gc, b)

	tree.Objects = append(tree.Objects, mc)
	TreeArr = append(TreeArr, b)
}

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
		success := localFileCopy(dd, nd)
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
	s := NewFileSelector()
	selectFromTree := xwidget.NewFileTree(storage.NewFileURI(getLocalHomeDir()))
	selectedTree := container.NewGridWithRows(10)

	selectFromTree.OnSelected = func(uid string) {
		exists := s.GetFile(uid)
		if exists {
			deselectButton(uid, s, selectedTree)
		} else {
			selectButton(uid, mainapp, s, selectedTree)
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
