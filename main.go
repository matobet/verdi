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

	reply, err := cmd.Run("AddVM", map[string]interface{}{
		"host": "7dacd585-309e-4a67-976b-ff06f8960377",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}
