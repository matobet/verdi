package env

import (
	"github.com/garyburd/redigo/redis"
	"github.com/matobet/verdi/model"
)

type RedisPool interface {
	Redis() Redis
}

type Redis interface {
	redis.Conn

	LoadScripts() error

	RedisReader
	Redlock
	Tx() RedisTx
}

type RedisReader interface {
	Exists(key string) (bool, error)
	Get(data interface{}) error
	GetString(key string) (string, error)
	HGetString(key, field string) (string, error)
}

type RedisWriter interface {
	Put(data interface{})
	Delete(data interface{})
}

type Redlock interface {
	Lock(lock string) (acquired bool, err error)
	Unlock(lock string) (released bool, err error)
}

type RedisTx interface {
	RedisWriter

	Begin() RedisTx
	Commit() error
}

type Commander interface {
	Run(cmd string, params map[string]interface{}) (map[string]interface{}, error)
}

type Virter interface {
	Virt() Virt
}

type Virt interface {
	StartVM(vm *model.VM) error
	StopVM(vm *model.VM) error
}

type Backend interface {
	RedisPool
	Virter
	Commander
}
