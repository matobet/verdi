package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/matobet/verdi/config"
)

func (conn *Conn) Lock(lock string) (acquired bool, err error) {
	ret, err := redis.String(conn.Do("SET", lock,
		config.Conf.HostID, "EX", config.Conf.CommandTimeout, "NX"))

	return ret != "", err
}

func (conn *Conn) Unlock(lock string) (released bool, err error) {
	ret, err := redis.Int(unlockScript.Do(conn, lock, config.Conf.HostID))
	return ret == 1, err
}

var unlockScript = redis.NewScript(1, `
	if redis.call('GET', KEYS[1]) == ARGV[1] then
		return redis.call('DEL', KEYS[1])
	end
	return 0
`)
