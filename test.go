package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Ashish9738/go-ciaos/ciaos"
	"github.com/Ashish9738/go-ciaos/config"
)

func testConfig() config.Config {
	return config.Config{
		APIURL:        "http://localhost:8080",
		UserId:        "45678",
		UserAccessKey: "testaccesskey",
	}
}

func readFilesFromDirectory(folderPath string) ([][]byte, error) {
	var byteStrings [][]byte

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		filePath := folderPath + "/" + file.Name()
		if !file.IsDir() {
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %v", filePath, err)
			}
			byteStrings = append(byteStrings, fileData)
		}
	}
	return byteStrings, nil
}

func testPost() {
	client, err := ciaos.Ciaos(testConfig())
	if err != nil {
		log.Fatalf("Error initializing Ciaos client: %v", err)
	}

	folderPath := "/home/ash/Pictures/wallpaper"

	byteStrings, err := readFilesFromDirectory(folderPath)
	if err != nil {
		log.Fatalf("Error reading files from directory: %v", err)
	}

	response, err := client.PutBinary("abcd", byteStrings)
	if err != nil {
		log.Fatalf("Error during upload: %v", err)
	}

	fmt.Println("Upload successful:", response.Status)
}

func main() {
	testPost()
}
