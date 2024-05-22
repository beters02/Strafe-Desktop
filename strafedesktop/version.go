package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func downloadVersionJson() ([]byte, error) {
	ret := []byte{}

	out, err := os.Create("temp_version.json")

	if err != nil {
		return ret, err
	}

	defer out.Close()

	resp, err := http.Get("https://raw.githubusercontent.com/beters02/Strafe-Desktop/main/version.json")

	if err != nil {
		return ret, err
	}

	defer resp.Body.Close()
	fmt.Println(resp.Body)

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return ret, err
	}

	fmt.Println("Copied!")
	return ret, nil
}

func GetMostRecentVersion() {
	downloadVersionJson()
}

func GetLocalVersion() (map[string]interface{}, error) {
	table := map[string]interface{}{}

	jsonFile, err := os.Open("version.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return table, err
	}

	err = json.Unmarshal(byteValue, &table)
	if err != nil {
		fmt.Println(err)
		return table, err
	}

	fmt.Println("Getting version")
	fmt.Println(table["Version"])
	return table, err
}
