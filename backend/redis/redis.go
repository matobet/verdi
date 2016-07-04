package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/matobet/verdi/env"
)

type Conn struct {
	redis.Conn
}

var _ env.Redis = (*Conn)(nil)

func (c *Conn) Get(data interface{}) error {
	id := redisID(data)
	redisType := redisType(data)
	key := redisKeyWithTypeAndID(redisType, id)
	values, err := redis.Values(c.Do("HGETALL", key))
	if err != nil {
		return err
	}
	return redis.ScanStruct(values, data)
}

func (c *Conn) Put(data interface{}) {
	id := redisID(data)
	redisType := redisType(data)
	key := redisKeyWithTypeAndID(redisType, id)
	c.updateIndexes(redisType, id, data)
	c.Send("SADD", redisType, id)
	c.Send("HMSET", redis.Args{key}.AddFlat(data)...)
}

func (c *Conn) Delete(data interface{}) {
	id := redisID(data)
	redisType := redisType(data)
	key := redisKeyWithTypeAndID(redisType, id)
	c.updateIndexes(redisType, id, data)
	c.Send("SREM", redisType, id)
	c.Send("DEL", key)
}

func (conn *Conn) Exists(key string) (bool, error) {
	return redis.Bool(conn.Do("EXISTS", key))
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

func (c *Conn) updateIndexes(redisType, id string, data interface{}) {
	if indexer, ok := data.(redisIndexer); ok {
		c.SendScript("update-indexes", redis.Args{redisType, id}.AddFlat(indexer.RedisIndexes())...)
	}
}
