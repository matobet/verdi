package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/matobet/verdi/env"
)

type Conn struct {
	redis.Conn
}

var _ env.Redis = (*Conn)(nil)

func (c *Conn) Put(data interface{}) {
	id := redisID(data)
	redisType := redisType(data)
	key := redisKeyWithTypeAndID(redisType, id)
	c.Conn.Send("SADD", redisType, id)
	c.Conn.Send("HMSET", redis.Args{key}.AddFlat(data)...)
	if indexer, ok := data.(redisIndexer); ok {
		for field, value := range indexer.RedisIndexes() {
			c.Conn.Send("SET", redisIndexByTypeFieldAndValue(redisType, field, value), id)
		}
	}
}

func (c *Conn) Delete(data interface{}) {
	id := redisID(data)
	redisType := redisType(data)
	key := redisKeyWithTypeAndID(redisType, id)
	c.Send("SREM", redisType, id)
	c.Send("DEL", key)
	if indexer, ok := data.(redisIndexer); ok {
		for field, value := range indexer.RedisIndexes() {
			c.Conn.Send("DEL", redisIndexByTypeFieldAndValue(redisType, field, value))
		}
	}
}

func (conn *Conn) GetString(key string) (string, error) {
	return redis.String(conn.Do("GET", key))
}

func (conn *Conn) HGetString(key, field string) (string, error) {
	return redis.String(conn.Do("HGET", key, field))
}

func (conn *Conn) Tx() env.RedisTx {
	return &Tx{*conn}
}
