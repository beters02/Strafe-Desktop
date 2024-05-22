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

	resp, err := http.Get("https://raw.githubusercontent.com/beters02/Strafe-Desktop/main/strafedesktop/version.json")

	if err != nil {
		return ret, err
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return ret, err
	}

	jsonFile, err := os.Open("temp_version.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return ret, err
	}

	return byteValue, nil
}

func deleteVersionJson() {
	_, err := os.Stat("temp_version.json")
	if err != nil {
		return
	}
	os.Remove("temp_version.json")
}

func GetMostRecentVersion() string {
	table := map[string]interface{}{}
	bytec, _ := downloadVersionJson()
	deleteVersionJson()

	err := json.Unmarshal(bytec, &table)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	a := table["Version"].(string)
	return a
}

func GetLocalVersion() (string, error) {
	table := map[string]interface{}{}

	jsonFile, err := os.Open("version.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	err = json.Unmarshal(byteValue, &table)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	a := table["Version"].(string)
	return a, nil
}

func DownloadRecentBuild(rv string) {
	out, err := os.Create("builds/Strafe-Desktop-" + rv + ".exe")

	if err != nil {
		fmt.Printf("Could not get recent build. %v\n", err)
		return
	}

	defer out.Close()

	resp, err := http.Get("https://raw.githubusercontent.com/beters02/Strafe-Desktop/main/strafedesktop/builds/Strafe-Desktop.exe")

	if err != nil {
		fmt.Printf("Could not get recent build. %v\n", err)
		return
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		fmt.Printf("Could not get recent build. %v\n", err)
		return
	}

	bytec, _ := json.Marshal(map[string]interface{}{"Version": rv})
	os.Remove("version.json")
	os.WriteFile("version.json", bytec, 0777)
	fmt.Printf("Recent build downloaded! Please close this exe and open the new one.")
}
