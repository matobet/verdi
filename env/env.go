package env

import (
	"github.com/garyburd/redigo/redis"
	"github.com/matobet/verdi/model"
)

type (
	RedisPool interface {
		Redis() Redis
	}

	Redis interface {
		redis.Conn

		LoadScripts() error

		RedisReader
		Redlock
		Tx() RedisTx
	}

	RedisReader interface {
		Exists(key string) (bool, error)
		Get(data interface{}) error
		GetString(key string) (string, error)
		HGetString(key, field string) (string, error)
	}

	RedisWriter interface {
		Put(data interface{})
		Delete(data interface{})
	}

	Redlock interface {
		Lock(lock string) (acquired bool, err error)
		Unlock(lock string) (released bool, err error)
	}

	RedisTx interface {
		RedisWriter

		Begin() RedisTx
		Commit() error
	}

	Commander interface {
		Run(cmd string, params map[string]interface{}) (map[string]interface{}, error)
	}

	Virter interface {
		Virt() Virt
	}

	Virt interface {
		StartVM(vm *model.VM) error
		StopVM(vm *model.VM) error
	}

	Backend interface {
		RedisPool
		Virter
		Commander
	}
)
