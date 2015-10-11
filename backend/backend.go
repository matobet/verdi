package backend

import (
	"github.com/matobet/verdi/backend/cmd"
	"github.com/matobet/verdi/config"
)

func Init() (err error) {

	err = cmd.Init()
	if err != nil {
		return
	}

	go cmd.Listen(cmd.GlobalQueue)
	go cmd.Listen(cmd.QueueByClassAndID(cmd.Host, config.Conf.HostID))

	return
}
