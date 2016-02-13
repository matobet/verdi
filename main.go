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

	reply, err := backend.Run("UpdateVM", &cmd.UpdateVmParams{
		ID:   "09950252-456e-49fd-9c6d-993ea961bf08",
		Name: "xxxx",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}
