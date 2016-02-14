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
		log.Fatal("Failed to load configuration: ", err)
	}

	backend, err := backend.Init()
	if err != nil {
		log.Fatal("Failed to initialize backend: ", err)
	}

	reply, err := backend.Run("RemoveVM", &cmd.IDParams{
		ID: "4392ce9e-5a9b-442e-9037-91764e14b129",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}
