package test

import (
	"fmt"
	"log"

	"github.com/Ashish9738/go-ciaos/ciaos"
)

func TestAppend() {
	client, err := ciaos.Ciaos(GetConfig())

	if err != nil {
		log.Fatalf("Error initializing Ciaos client: %v", err)
	}

	folderPath := "/home/ash/Pictures/Wallpapers"

	byteStrings, err := ReadFilesFromDir(folderPath)
	fmt.Println("byte strings", byteStrings)
	fmt.Println("client", client)
	if err != nil {
		log.Fatalf("Error reading files from directory: %v", err)
	}

	response, err := client.Append("Ash here", byteStrings)
	if err != nil {
		log.Fatalf("Error during upload: %v", err)
	}

	fmt.Println("Upload successful:", response.Status)

}