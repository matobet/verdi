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

	err := b.Redis().LoadScripts()
	if err != nil {
		return nil, err
	}

	go cmd.Listen(b, cmd.GlobalQueue)
	go cmd.Listen(b, cmd.QueueByClassAndID(cmd.Host, config.Conf.HostID))

	go scheduler.Listen(b, model.GlobalClusterID)

	return b, nil
}

func (b *backend) Redis() env.Redis {
	return b.redisPool.Redis()
}

func (b *backend) Run(command string, params interface{}) (map[string]interface{}, error) {
	return cmd.Run(b, command, structs.Map(params))
}
