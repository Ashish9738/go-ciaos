package test

import (
	"fmt"
	"log"
	"os"

	"github.com/Ashish9738/go-ciaos/ciaos"
)

func GetConfig() ciaos.Config {
	return ciaos.Config{
		APIURL:        "http://localhost:8080",
		UserId:        "45678",
		UserAccessKey: "testaccesskey",
	}
}

func ReadFilesFromDir(folderPath string) ([][]byte, error) {
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

func TestPutBinary() {
	client, err := ciaos.Ciaos(GetConfig())

	if err != nil {
		log.Fatalf("Error initializing Ciaos client: %v", err)
	}

	folderPath := "/home/ash/Pictures/wallpaper"

	byteStrings, err := ReadFilesFromDir(folderPath)
	fmt.Println("byte strings", byteStrings)
	fmt.Println("client", client)
	if err != nil {
		log.Fatalf("Error reading files from directory: %v", err)
	}

	response, err := client.PutBinary("jsdusv", byteStrings)
	if err != nil {
		log.Fatalf("Error during upload: %v", err)
	}

	fmt.Println("Upload successful:", response.Status)
}

// func test() {
// 	TestPutBinary()
// }
