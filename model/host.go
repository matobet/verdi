package model

type Host struct {
	ID             GUID   `redis:"id"`
	Name           string `redis:"name"`
	TotalMemSizeMB uint64 `redis:"total_mem_size_mb"`
	MemUsedMB      uint64 `redis:"mem_used_mb"`
}

type HostStats struct {
	ID             GUID   `redis:"id"`
	TotalMemSizeMB uint64 `redis:"total_mem_size_mb"`
	MemUsedMB      uint64 `redis:"mem_used_mb"`
}

func (host *Host) MemUsage() float32 {
	return 100.0 * float32(host.MemUsedMB) / float32(host.TotalMemSizeMB)
}

func (hs *HostStats) RedisType() string {
	return "Host"
}
