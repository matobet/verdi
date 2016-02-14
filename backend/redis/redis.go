package redis

import (
	"fmt"

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
		fmt.Println(indexer.RedisIndexes())
		updateIndexesScript.Send(c.Conn, redis.Args{redisType, id}.AddFlat(indexer.RedisIndexes())...)
	}
}

var updateIndexesScript = redis.NewScript(2, `
	local type = KEYS[1]
	local id = KEYS[2]
	local newIndexValues = {}
	local indexes = {}
	local nextKey
	for i, v in ipairs(ARGV) do
		if i % 2 == 1 then
			table.insert(indexes, v)
			nextKey = v
		else
			newIndexValues[nextKey] = v
		end
	end
	local old = redis.call('HMGET', string.format('%s:%s', type, id), unpack(indexes))
	local updates = 0
	for i, index in ipairs(indexes) do
		if old[i] and old[i] ~= newIndexValues[index] then
			redis.call('DEL', string.format('q:%s:%s:%s', type, index, old[i]))
		end
		if (not old[i] or old[i] ~= newIndexValues[index]) and newIndexValues[index] ~= '' then
			redis.call('SET', string.format('q:%s:%s:%s', type, index, newIndexValues[index]), id)
			updates = updates + 1
		end
	end
	return updates
`)
