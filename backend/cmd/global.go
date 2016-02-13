package cmd

import (
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/matobet/verdi/env"
	"github.com/matobet/verdi/model"
)

type IDParams struct {
	ID model.GUID
}

type AddVmParams struct {
	Name      string
	ClusterID model.GUID
}

func addVM(backend env.Backend, params *AddVmParams) (result interface{}, err error) {
	redis := backend.Redis()
	defer redis.Close()

	nameLock := vmNameLock(params.Name)
	nameLocked, err := redis.Lock(nameLock)
	if err != nil || !nameLocked {
		return nil, fmt.Errorf("VM with name '%s' is already being created", params.Name)
	}
	defer redis.Unlock(nameLock)

	existing, err := redis.GetString("q:VM:name:" + params.Name)
	if existing != "" {
		return nil, fmt.Errorf("VM with name '%s' already exists", params.Name)
	}

	vm := &model.VM{
		ID:        model.NewGUID(),
		Name:      params.Name,
		ClusterID: params.ClusterID,
	}

	tx := redis.Tx().Begin()
	tx.Put(vm)
	err = tx.Commit()

	return "Created", err
}

type UpdateVmParams struct {
	ID   model.GUID
	Name string
}

func updateVM(backend env.Backend, params *UpdateVmParams) (result interface{}, err error) {
	conn := backend.Redis()
	defer conn.Close()

	ok, err := conn.Exists("VM:" + params.ID.String())
	if err != nil {
		return
	}
	if !ok {
		return nil, fmt.Errorf("VM with ID '%s' does not exist", params.ID)
	}
	vm := &model.VM{
		ID:   params.ID,
		Name: params.Name,
	}

	// FIXME: needs generic lua script for proper update of indexes
	tx := conn.Tx().Begin()
	tx.Put(vm)
	err = tx.Commit()

	return "Updated", err
}

func runVM(backend env.Backend, params map[string]interface{}) (result interface{}, err error) {
	return nil, errors.New("Not implemented")
}

func stopVM(backend env.Backend, params map[string]interface{}) (result interface{}, err error) {
	return nil, errors.New("Not implemented")
}

func removeVM(backend env.Backend, params IDParams) (result interface{}, err error) {
	conn := backend.Redis()
	defer conn.Close()

	name, err := conn.HGetString("VM:"+params.ID.String(), "name")
	if err == redis.ErrNil {
		return nil, fmt.Errorf("VM with ID '%s' does not exist", params.ID)
	}
	if err != nil {
		return
	}

	tx := conn.Tx().Begin()
	tx.Delete(&model.VM{ID: params.ID, Name: name})
	err = tx.Commit()

	return "Removed", err
}

func vmNameLock(name string) string {
	return "lock:vm:name:" + name
}
