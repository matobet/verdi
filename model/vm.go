package model

type VM struct {
	ID        GUID   `redis:"id"`
	Name      string `redis:"name"`
	Status    Status `redis:"status"`
	ClusterID GUID   `redis:"cluster_id"`
}

func (vm *VM) RedisIndexes() map[string]string {
	return map[string]string{
		"name": vm.Name,
	}
}
