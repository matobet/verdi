package cmd

import (
	"errors"

	"github.com/matobet/verdi/backend/scheduler"
	"github.com/matobet/verdi/config"
	"github.com/matobet/verdi/env"
	"github.com/matobet/verdi/model"
	maps "github.com/mitchellh/mapstructure"
)

func pingHost(_backend env.Backend, _params map[string]interface{}) (result interface{}, err error) {
	return "PONG", nil
}

type ClusterHostParams struct {
	HostID    model.GUID `structs:"host_id" mapstructure:"host_id"`
	ClusterID model.GUID `structs:"cluster_id" mapstructure:"cluster_id"`
}

var ErrMissingClusterID = errors.New("'cluster_id' must be specified")

func (params *ClusterHostParams) validate() error {
	if params.ClusterID == "" {
		return ErrMissingClusterID
	}
	return nil
}

func addHostToCluster(backend env.Backend, params map[string]interface{}) (result interface{}, err error) {
	var p ClusterHostParams
	maps.Decode(params, &p)

	if err = p.validate(); err != nil {
		return nil, err
	}

	conn := backend.Redis()
	defer conn.Close()

	hostID := config.Conf.HostID

	conn.Send("MULTI")
	conn.Send("SADD", "Cluster:"+p.ClusterID, hostID)
	conn.Send("SADD", "Host:"+hostID+":clusters", p.ClusterID)
	_, err = conn.Do("EXEC")

	go scheduler.Listen(backend, p.ClusterID)

	return "Added", err
}

func removeHostFromCluster(backend env.Backend, params map[string]interface{}) (result interface{}, err error) {
	var p ClusterHostParams
	maps.Decode(params, &p)

	if err = p.validate(); err != nil {
		return nil, err
	}

	conn := backend.Redis()
	defer conn.Close()

	hostID := config.Conf.HostID

	conn.Send("MULTI")
	conn.Send("SREM", "Cluster:"+p.ClusterID, hostID)
	conn.Send("SREM", "Host:"+hostID+":clusters", p.ClusterID)
	_, err = conn.Do("EXEC")

	scheduler.StopListen(backend, p.ClusterID)

	return "Removed", err
}
