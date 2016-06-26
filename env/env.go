package env

import "github.com/garyburd/redigo/redis"

type RedisPool interface {
	Redis() Redis
}

type Redis interface {
	redis.Conn

	LoadScripts() error

	Exists(key string) (bool, error)
	GetString(key string) (string, error)
	HGetString(key, field string) (string, error)

	Tx() RedisTx

	Lock(lock string) (acquired bool, err error)

	Unlock(lock string) (released bool, err error)
}

type RedisWriter interface {
	Put(data interface{})
	Delete(data interface{})
}

type RedisTx interface {
	RedisWriter

	Begin() RedisTx

	Commit() error
}

type Commander interface {
	Run(cmd string, params interface{}) (map[string]interface{}, error)
}

type Backend interface {
	RedisPool
	Commander
}
