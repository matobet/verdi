package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/matobet/verdi/config"
)

// Lock locks the given key in redis
func (conn *Conn) Lock(lock string) (acquired bool, err error) {
	ret, err := redis.String(conn.Do("SET", lock,
		config.Conf.HostID, "EX", config.Conf.CommandTimeout, "NX"))

	return ret != "", err
}

// Unlock unlocks the given key in reds provided we still own it
func (conn *Conn) Unlock(lock string) (released bool, err error) {
	ret, err := redis.Int(conn.DoScript("unlock", lock, config.Conf.HostID))
	return ret == 1, err
}
