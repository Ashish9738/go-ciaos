package test

import (
	"fmt"
	"log"

	"github.com/Ashish9738/go-ciaos/ciaos"
)

func TestGet() {
	client, err := ciaos.Ciaos(GetConfig())

	if err != nil {
		log.Fatalf("Error initializing Ciaos client: %v", err)
	}

	response, err := client.Get("abcd")
	if err != nil {
		log.Fatalf("Error during upload: %v", err)
	}

	fmt.Println("Upload successful:", response)

}
