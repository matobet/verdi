package backend

import (
	"github.com/matobet/verdi/backend/cmd"
	"github.com/matobet/verdi/backend/redis"
	"github.com/matobet/verdi/backend/scheduler"
	"github.com/matobet/verdi/backend/virt"
	"github.com/matobet/verdi/config"
	"github.com/matobet/verdi/env"
	"github.com/matobet/verdi/model"

	"github.com/fatih/structs"
)

type backend struct {
	redisPool *redis.Pool
	virt      virt.Conn
}

func Init() (env.Backend, error) {

	//	virt, err := virt.NewConn()
	//	if err != nil {
	//		return nil, fmt.Errorf("backend: Error connecting to libvirt: '%s'", err)
	//	}

	b := &backend{
		redisPool: redis.NewPool(),
		virt:      nil, //virt,
	}

	go cmd.Listen(b, cmd.GlobalQueue)
	go cmd.Listen(b, cmd.QueueByClassAndID(cmd.Host, config.Conf.HostID))

	go scheduler.Listen(b, model.GlobalClusterID)

	return b, nil
}

func (back *backend) Redis() env.Redis {
	return back.redisPool.Redis()
}

func (back *backend) Run(command string, params interface{}) (map[string]interface{}, error) {
	return cmd.Run(back, command, structs.Map(params))
}
