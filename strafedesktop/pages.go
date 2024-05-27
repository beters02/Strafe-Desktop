package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
	dirToZip   bool
	currFiles  map[string]string
	currFilesN int
	maxFiles   int
}

func NewFileSelector(dirToZip bool) *FileSelector {
	return &FileSelector{
		currFiles:  map[string]string{},
		currFilesN: 0,
		maxFiles:   10,
		dirToZip:   dirToZip,
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

// Upload/Download Shared Functions

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

func selectButton(uid string, mainapp *Application, s *FileSelector, tree *fyne.Container, dropDown bool) {
	didSelect := s.Select(uid)
	if !didSelect {
		return
	}

	path := fixUid(uid)
	isDir := isLocalFileDir(path)

	if isDir {
		if !s.dirToZip {
			s.Deselect(uid)
			return
		}
	}

	b := widget.NewButton(uid, func() { deselectButton(uid, s, tree) })
	var mc *fyne.Container

	if dropDown {
		o := getPossibleOptions(mainapp)
		gc := newDropdownSelection("", o, func(str string) { s.SetFileString(uid, str) })
		mc = container.NewHBox(gc, b)
	} else {
		mc = container.NewHBox(b)
	}

	tree.Objects = append(tree.Objects, mc)
	TreeArr = append(TreeArr, b)
}

//
