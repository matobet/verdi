package util

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/matobet/verdi/config"
)

func NewRedisPool() *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", config.Conf.RedisServer)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
}
