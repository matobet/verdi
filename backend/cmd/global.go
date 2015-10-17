package cmd

import (
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/matobet/verdi/env"
	"github.com/matobet/verdi/model"
	maps "github.com/mitchellh/mapstructure"
)

type IDParams struct {
	ID model.GUID `mapstructure:"id" structs:"id"`
}

type AddVmParams struct {
	Name      string     `structs:"name" mapstructure:"name"`
	ClusterID model.GUID `structs:"cluster_id" mapstructure:"cluster_id"`
}

func addVM(backend env.Backend, params map[string]interface{}) (result interface{}, err error) {
	var p AddVmParams
	if err = maps.Decode(params, &p); err != nil {
		return
	}

	conn := backend.Redis()
	defer conn.Close()

	nameLock := "lock:vm:name:" + p.Name
	nameLocked, err := conn.Lock(nameLock)
	if err != nil || !nameLocked {
		return nil, fmt.Errorf("VM with name '%s' is already being created", p.Name)
	}
	defer conn.Unlock(nameLock)

	existing, err := conn.GetString("q:VM:name:" + p.Name)
	if existing != "" {
		return nil, fmt.Errorf("VM with name '%s' already exists", p.Name)
	}

	vm := &model.VM{
		ID:        model.NewGUID(),
		Name:      p.Name,
		ClusterID: p.ClusterID,
	}

	tx := conn.Tx().Begin()
	tx.Put(vm)
	err = tx.Commit()

	return "Created", err
}

func runVM(backend env.Backend, params map[string]interface{}) (result interface{}, err error) {
	return nil, errors.New("Not implemented")
}

func stopVM(backend env.Backend, params map[string]interface{}) (result interface{}, err error) {
	return nil, errors.New("Not implemented")
}

func removeVM(backend env.Backend, params map[string]interface{}) (result interface{}, err error) {
	var p IDParams
	if err = maps.Decode(params, &p); err != nil {
		return
	}

	conn := backend.Redis()
	defer conn.Close()

	name, err := conn.HGetString("VM:"+p.ID.String(), "name")
	if err == redis.ErrNil {
		return nil, fmt.Errorf("VM with ID '%s' does not exist", p.ID)
	}
	if err != nil {
		return
	}

	tx := conn.Tx().Begin()
	tx.Delete(&model.VM{ID: p.ID, Name: name})
	err = tx.Commit()

	return "Removed", err
}
