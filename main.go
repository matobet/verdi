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

	backend, err := backend.Init()
	if err != nil {
		log.Fatal("Failed to initialize backend", err)
	}

	reply, err := backend.Run("RemoveVM", &cmd.IDParams{
		ID: "6bf2ca64-f2a0-4777-9c8f-731f9fe53a64",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}
