// Blueprint of potential explorer class responsible for retrieving files from the users file system to upload to the server

package main

import (
	"errors"
	"os"
	"strings"
)

// NOTE: this isn't multi-Unicode-codepoint aware, like specifying skintone or
//
//	gender of an emoji: https://unicode.org/emoji/charts/full-emoji-modifiers.html
func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

//

func getDirFiles(explorer *Explorer, dir string) ([]os.DirEntry, error) {
	var farr []os.DirEntry = explorer.Cache[dir]
	var err error
	if farr == nil {
		farr, err = os.ReadDir(dir)
		if err != nil {
			explorer.Cache[dir] = farr
		}
	}
	return farr, err
}

type Explorer struct {
	CurrentDirectory string
	Cache            map[string][]os.DirEntry // this is the map of arrays type in go.... nice nice nice
}

func CreateExplorer() Explorer {
	explorer := Explorer{
		CurrentDirectory: "C:/Users/{getWindowsUserNameHere}/Desktop",
		Cache:            map[string][]os.DirEntry{},
	}
	return explorer
}

func (explorer *Explorer) Next(dir string) ([]os.DirEntry, error) {
	newDir := explorer.CurrentDirectory + dir + "/"

	farr, err := getDirFiles(explorer, newDir)
	if err != nil {
		return nil, errors.New("could not read previous directory")
	}

	explorer.CurrentDirectory = newDir
	return farr, err
}

func (explorer *Explorer) Back() ([]os.DirEntry, error) {
	slashIndex := strings.LastIndex(explorer.CurrentDirectory, "/")
	var previousDir string = substr(explorer.CurrentDirectory, 1, slashIndex)

	farr, err := getDirFiles(explorer, previousDir)
	if err != nil {
		return nil, errors.New("could not read previous directory")
	}

	explorer.CurrentDirectory = previousDir
	return farr, err
}
