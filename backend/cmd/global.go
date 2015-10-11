package cmd

import (
	"errors"

	maps "github.com/mitchellh/mapstructure"
)

type AddVmParams struct {
	Name      string `mapstructure:"name"`
	ClusterID string `mapstructure:"cluster_id"`
}

func addVM(params map[string]interface{}) (result interface{}, err error) {
	var p AddVmParams
	if err = maps.Decode(params, &p); err != nil {
		return
	}

	// lock:vm:name:`name`
	return nil, errors.New("Not implemented")
}

func runVM(params map[string]interface{}) (result interface{}, err error) {
	return nil, errors.New("Not implemented")
}

func stopVM(params map[string]interface{}) (result interface{}, err error) {
	return nil, errors.New("Not implemented")

}

func deleteVM(params map[string]interface{}) (result interface{}, err error) {
	return nil, errors.New("Not implemented")
}
