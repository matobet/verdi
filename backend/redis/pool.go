package redis

import (
	"time"

	r "github.com/garyburd/redigo/redis"
	"github.com/matobet/verdi/config"
	"github.com/matobet/verdi/env"
)

type Pool struct {
	r.Pool
}

func NewPool() *Pool {
	return &Pool{r.Pool{
		Dial: func() (r.Conn, error) {
			conn, err := r.Dial("tcp", config.Conf.RedisServer)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow: func(conn r.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}}
}

func (p *Pool) Redis() env.Redis {
	return &Conn{p.Get()}
}
