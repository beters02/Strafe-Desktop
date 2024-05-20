package main

import (
	"io/fs"

	strafe "github.com/beters02/Strafe"
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
