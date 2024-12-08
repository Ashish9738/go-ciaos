package test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Ashish9738/go-ciaos/ciaos"
)

func TestPut() {
	client, err := ciaos.Ciaos(GetConfig())
	if err != nil {
		log.Fatalf("Error initializing Ciaos client: %v", err)
	}

	folderPath := "/home/ash/Pictures/wallpaper/"

	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	key := "Ash here"
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(folderPath, file.Name())

			response, err := client.Put(filePath, key)
			if err != nil {
				log.Fatalf("Error during upload of %s: %v", file.Name(), err)
			}

			fmt.Printf("Upload successful for %s: %s\n", key, response.Status)
		}
	}
}

func TestPutAlt() {
	client, err := ciaos.Ciaos(GetConfig())
	if err != nil {
		log.Fatalf("Error initializing Ciaos client: %v", err)
	}

	folderPath := "/home/ash/Pictures/wallpaper/"

	baseKey := "wallpaper_"

	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	for i, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(folderPath, file.Name())

			key := fmt.Sprintf("%s%d", baseKey, i)

			fmt.Printf("Uploading file: %s with key: %s\n", file.Name(), key)

			response, err := client.Put(filePath, key)
			if err != nil {
				log.Fatalf("Error during upload of %s: %v", file.Name(), err)
			}

			fmt.Printf("Upload successful for key %s: %s\n", key, response.Status)
		}
	}
}

// func test() {
// 	// TestPut()
// 	TestPut()
// }
