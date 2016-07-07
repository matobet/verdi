package backend

import (
	"fmt"
	"log"
	"time"

	"github.com/matobet/verdi/backend/cmd"
	"github.com/matobet/verdi/backend/redis"
	"github.com/matobet/verdi/backend/scheduler"
	"github.com/matobet/verdi/backend/virt"
	"github.com/matobet/verdi/config"
	"github.com/matobet/verdi/env"
	"github.com/matobet/verdi/model"
	"github.com/shirou/gopsutil/mem"
)

type backend struct {
	redisPool *redis.Pool
	virt      *virt.Conn
}

func Init() (env.Backend, error) {

	virt, err := virt.NewConn()
	if err != nil {
		return nil, fmt.Errorf("backend: Error connecting to libvirt: '%s'", err)
	}

	b := &backend{
		redisPool: redis.NewPool(),
		virt:      virt,
	}

	err = b.Redis().LoadScripts()
	if err != nil {
		return nil, err
	}

	go cmd.Listen(b, cmd.GlobalQueue)
	go cmd.Listen(b, cmd.QueueByClassAndID(cmd.Host, config.Conf.HostID.String()))

	go scheduler.Listen(b, model.GlobalClusterID)

	go b.monitor()

	return b, nil
}

func (b *backend) Redis() env.Redis {
	return b.redisPool.Redis()
}

func (b *backend) Virt() env.Virt {
	return b.virt
}

func (b *backend) Run(command string, params map[string]interface{}) (map[string]interface{}, error) {
	return cmd.Run(b, command, params)
}

func (b *backend) monitor() {
	redis := b.Redis()
	hostStats := model.HostStats{ID: config.Conf.HostID}
	for {
		time.Sleep(1 * time.Second)
		memStats, err := mem.VirtualMemory()
		if err != nil {
			log.Printf("Error gathering memory stats: %s", err)
			continue
		}
		log.Printf("%+v", memStats)
		hostStats.TotalMemSizeMB = memStats.Total
		hostStats.MemUsedMB = memStats.Used
		tx := redis.Tx().Begin()
		tx.Put(hostStats)
		tx.Commit()
	}
}
