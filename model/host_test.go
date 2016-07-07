package model

import "testing"

func TestHost_MemUsage(t *testing.T) {
	host := Host{TotalMemSizeMB: 2048, MemUsedMB: 512}
	if host.MemUsage() != 25 {
		t.Errorf("Expected memory utilization to be 25%%, instead got: %f", host.MemUsage())
	}
}
