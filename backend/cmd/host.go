package cmd

import "errors"

func pingHost(_params map[string]interface{}) (result interface{}, err error) {
	return "PONG", nil
}

func addHostToCluster(params map[string]interface{}) (result interface{}, err error) {
	return nil, errors.New("Not implemented")
}

func removeHostFromCluster(params map[string]interface{}) (result interface{}, err error) {
	return nil, errors.New("Not implemented")
}
