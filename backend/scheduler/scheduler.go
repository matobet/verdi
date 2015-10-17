package scheduler

import (
	"log"

	"github.com/matobet/verdi/env"
	"github.com/matobet/verdi/model"
)

func Listen(backend env.Backend, clusterID model.GUID) {
	log.Println("started listening on cluster", clusterID, "commands ...")

}

func StopListen(backend env.Backend, clusterID model.GUID) {
	log.Println("stopped listening on cluster", clusterID, "commands ...")
}
