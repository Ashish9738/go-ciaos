package test

import (
	"fmt"
	"log"

	"github.com/Ashish9738/go-ciaos/ciaos"
)

func TestUpadateKey() {
	client, err := ciaos.Ciaos(GetConfig())
	if err != nil {
		log.Fatalf("Error initializing Ciaos client: %v", err)
	}

	oldKey := "Ash here"
	newKey := "Ash Up"

	resp, err := client.UpdateKey(oldKey, newKey)
	if err != nil {
		log.Fatalf("Error during updating key: %v", err)
	}

	fmt.Printf("updating key successful for key %s: %s\n", oldKey, newKey)
	fmt.Println("response", resp)

}
