package main

import (
	"bufio"
	"io"
	"io/fs"
	"os"
	"os/user"
	"strings"

	strafe "github.com/beters02/Strafe"
	"github.com/fatih/color"
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

func getServerFilesAt(net strafe.Net, path string) []fs.FileInfo {
	prePath := "/shr/strafe"
	fi, err := net.Client.ReadDir(prePath + path)
	if err != nil {
		panic(err)
	}
	return fi
}

func isLocalFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isLocalFileDir(path string) bool {
	isDir := false

	if substr(path, len(path), len(path)) != "/" {
		path = path + "/"
	}

	fi, err := os.Stat(path)
	if err == nil {
		isDir = fi.IsDir()
	}

	return isDir
}

func getLocalHomeDir() string {
	u, err := user.Current()
	if err != nil {
		return "/"
	}
	return u.HomeDir
}

func localCloseFile(f *os.File) {
	if err := f.Close(); err != nil {
		panic(err)
	}
}

func localFileCopy(fromPath string, toPath string, doMove bool) bool {
	success := false
	errcolor := color.New(color.FgRed)

	f, err := os.Open(fromPath)
	if err != nil {
		errcolor.Printf("\nCould not open path : %v ... Does the file exist?", fromPath)
		return success
	}

	n, err := os.Create(toPath)
	if err != nil {
		errcolor.Printf("\nCould not create file : %v", toPath)
		return success
	}
	defer localCloseFile(n)

	var r *bufio.Reader = bufio.NewReader(f)
	var w *bufio.Writer = bufio.NewWriter(n)

	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	success = true
	lsuc := true

	for {
		// read a chunk
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			errcolor.Printf("\nCould not create file : %v", toPath)
			lsuc = false
			break
		}

		if n == 0 {
			break
		}

		// write a chunk
		if _, err := w.Write(buf[:n]); err != nil {
			errcolor.Printf("\nCould not create file : %v", toPath)
			lsuc = false
			break
		}
	}

	localCloseFile(f)
	if doMove {
		os.RemoveAll(fromPath)
	}

	return success && lsuc
}

func getFileName(dir string) string {
	slashInd := strings.LastIndex(dir, "/")
	return substr(dir, slashInd+1, len(dir))
}

func secToLastIndex(s string, subs string) int {
	li := strings.LastIndex(s, subs)
	if li == -1 {
		return -1
	}

	s1 := substr(s, 0, li)
	sli := strings.LastIndex(s1, subs)

	if sli == -1 {
		return li
	}

	return sli
}

func fixUid(uid string) string {
	return substr(uid, 7, len(uid))
}

func getFileType(name string) string {
	di := strings.LastIndex(name, ".")
	if di == -1 {
		return "folder"
	}
	return substr(name, 0, di)
}
