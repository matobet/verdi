package main

import (
	"fmt"
	"log"

	"github.com/matobet/verdi/backend"
	"github.com/matobet/verdi/backend/cmd"
	"github.com/matobet/verdi/config"
)

func main() {

	err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration", err)
	}

	err = backend.Init()
	if err != nil {
		log.Fatal("Failed to initialize backend", err)
	}

	reply, err := cmd.Run("RemoveVM", map[string]interface{}{
		"id": "bbe08913-c1ff-44a6-b0ba-423f23725135",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}
