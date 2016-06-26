package model

type VM struct {
	ID        GUID   `redis:"id"`
	Name      string `redis:"name"`
	Status    Status `redis:"status"`
	MemSizeMB uint64 `redis:"mem_size_mb"`
}

func (vm *VM) RedisIndexes() map[string]string {
	return map[string]string{
		"name": vm.Name,
	}
}
