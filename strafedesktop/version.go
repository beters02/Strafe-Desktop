package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func downloadLinkAsFile(where string, link string) (bool, error) {
	out, err := os.Create(where)
	if err != nil {
		return false, err
	}
	defer out.Close()

	resp, err := http.Get(link)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return false, err
	}
	return true, nil
}

func jsonFileToTable(where string) (map[string]interface{}, error) {
	table := map[string]interface{}{}

	jsonFile, err := os.Open(where)
	if err != nil {
		return table, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return table, err
	}

	err = json.Unmarshal(byteValue, &table)
	if err != nil {
		return table, err
	}

	return table, nil
}

func GetMostRecentVersion() (string, error) {
	didDownload, err := downloadLinkAsFile("temp_version.json", "https://raw.githubusercontent.com/beters02/Strafe-Desktop/main/strafedesktop/version.json")
	if !didDownload {
		return "", err
	}

	table, err := jsonFileToTable("temp_version.json")
	if err != nil {
		return "", err
	}

	os.Remove("temp_version.json")
	a := table["Version"].(string)
	return a, nil
}

func GetLocalVersion() (string, error) {
	table, err := jsonFileToTable("version.json")
	if err != nil {
		return "", err
	}
	a := table["Version"].(string)
	return a, nil
}

func CheckForNewerVersion(installedVersion string, latestVersion string) (isNewer bool) {
	installedTotal := getVersionTotal(installedVersion)
	latestTotal := getVersionTotal(latestVersion)
	return latestTotal > installedTotal
}

func UpdatePrompt(recentVersion string) (didUpdate bool) {
	var s string
	fmt.Println("New version found! Would you like to update? y N")
	fmt.Scan(&s)
	if s == "y" || s == "Y" {
		DownloadRecentBuild(recentVersion)
		return true
	}

	fmt.Println("Not updating")
	return false
}

func DownloadRecentBuild(rv string) {
	_, err := downloadLinkAsFile("./Strafe-Desktop-"+rv+".exe", "https://raw.githubusercontent.com/beters02/Strafe-Desktop/main/strafedesktop/builds/Strafe-Desktop.exe")

	if err != nil {
		fmt.Printf("Could not get recent build. %v\n", err)
		return
	}

	bytec, _ := json.Marshal(map[string]interface{}{"Version": rv})
	os.Remove("version.json")
	os.WriteFile("version.json", bytec, 0777)

	fmt.Printf("Recent build downloaded! Please close this exe and open the new one.")
	fmt.Printf("strafedesktop/Strafe-Desktop-" + rv + ".exe")
}

func getVersionTotal(vstr string) int {
	arr := strings.Split(vstr, ".")
	num := 0
	for _, v := range arr {
		i, _ := strconv.Atoi(v)
		num += i
	}
	return num
}
