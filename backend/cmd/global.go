package cmd

import (
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/matobet/verdi/config"
	"github.com/matobet/verdi/model"
	maps "github.com/mitchellh/mapstructure"
)

type IDParams struct {
	ID model.GUID `mapstructure:"id"`
}

type AddVmParams struct {
	Name      string     `mapstructure:"name"`
	ClusterID model.GUID `mapstructure:"cluster_id"`
}

func addVM(params map[string]interface{}) (result interface{}, err error) {
	var p AddVmParams
	if err = maps.Decode(params, &p); err != nil {
		return
	}

	conn := redisPool.Get()
	defer conn.Close()

	lock, err := redis.String(conn.Do("SET", "lock:vm:name:"+p.Name,
		config.Conf.HostID, "EX", config.Conf.CommandTimeout, "NX"))

	if err != nil || lock == "" {
		return nil, fmt.Errorf("VM with name '%s' is already being created", p.Name)
	}

	existing, err := redis.String(conn.Do("GET", "query:VM:by:name:"+p.Name))
	if existing != "" {
		return nil, fmt.Errorf("VM with name '%s' already exists", p.Name)
	}

	vm := &model.VM{
		ID:        model.NewGUID(),
		Name:      p.Name,
		ClusterID: p.ClusterID,
	}

	conn.Send("MULTI")
	conn.Send("SADD", "VM", vm.ID)
	conn.Send("HMSET", redis.Args{"VM:" + vm.ID}.AddFlat(vm)...)
	conn.Send("SET", "query:VM:by:name:"+p.Name, vm.ID)
	_, err = conn.Do("EXEC")

	return "Created", err
}

func runVM(params map[string]interface{}) (result interface{}, err error) {
	return nil, errors.New("Not implemented")
}

func stopVM(params map[string]interface{}) (result interface{}, err error) {
	return nil, errors.New("Not implemented")

}

func removeVM(params map[string]interface{}) (result interface{}, err error) {
	var p IDParams
	if err = maps.Decode(params, &p); err != nil {
		return
	}

	conn := redisPool.Get()
	defer conn.Close()

	name, err := redis.String(conn.Do("HGET", "VM:"+p.ID, "name"))
	if err == redis.ErrNil {
		return nil, fmt.Errorf("VM with ID '%s' does not exist", p.ID)
	}
	if err != nil {
		return
	}

	conn.Send("MULTI")
	conn.Send("SREM", "VM", p.ID)
	conn.Send("DEL", "VM:"+p.ID, "query:VM:by:name:"+name)
	_, err = conn.Do("EXEC")

	return "Removed", err
}
