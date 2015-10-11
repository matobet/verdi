package db

import (
	"fmt"
	"reflect"

	"github.com/garyburd/redigo/redis"
)

type Conn redis.Conn

// Connect returns connection to redis store
func Connect() (Conn, error) {
	return redis.Dial("tcp", ":6379")
}

// Put creates or updates entity to redis store
func Put(c redis.Conn, data interface{}) error {
	id := redisID(data)
	redisType := redisType(data)
	key := redisKeyWithTypeAndID(redisType, id)
	c.Send("MULTI")
	c.Send("SADD", redisType, id)
	c.Send("HMSET", redis.Args{key}.AddFlat(data)...)
	_, err := c.Do("EXEC")
	return err
}

func PutAll(c redis.Conn, data []interface{}) error {
	if len(data) == 0 {
		return nil
	}
	redisType := redisType(data[0])
	for _, datum := range data {
		id := redisID(datum)
		key := redisKeyWithTypeAndID(redisType, id)
		c.Send("SADD", redisType, id)
		c.Send("HMSET", redis.Args{key}.AddFlat(datum)...)
	}
	return c.Flush()
}

// Get retrieves entity of given type and ID from redis store
func Get(c redis.Conn, id string, data interface{}) error {
	key := redisKeyWithID(data, id)
	values, err := redis.Values(c.Do("HGETALL", key))
	if err != nil {
		return err
	}

	return redis.ScanStruct(values, data)
}

func GetAll(c redis.Conn, dest interface{}) error {
	r := reflect.TypeOf(dest).Elem()
	p := r.Elem()
	t := p.Elem()
	redisType := redisType(reflect.New(t).Interface())
	ids, err := redis.Strings(c.Do("SMEMBERS", redisType))
	if err != nil {
		return err
	}
	for _, id := range ids {
		c.Send("HGETALL", redisKeyWithTypeAndID(redisType, id))
	}
	c.Flush()
	n := len(ids)
	destValue := reflect.ValueOf(dest).Elem()
	destValue.Set(reflect.MakeSlice(reflect.SliceOf(p), 0, n))
	for range ids {
		values, _ := redis.Values(c.Receive())
		entityValue := reflect.New(t)
		entity := entityValue.Interface()
		redis.ScanStruct(values, entity)
		destValue.Set(reflect.Append(destValue, entityValue))
	}

	return nil
}

func Delete(c redis.Conn, data interface{}) error {
	id := redisID(data)
	redisType := redisType(data)
	key := redisKeyWithTypeAndID(redisType, id)
	c.Send("MULTI")
	c.Send("SREM", redisType, id)
	c.Send("DEL", key)
	_, err := c.Do("EXEC")
	return err
}

type redisTyped interface {
	RedisType() string
}

type redisIDed interface {
	RedisID() string
}

func redisKey(data interface{}) string {
	key := redisID(data)
	return redisKeyWithID(data, key)
}

func redisKeyWithID(data interface{}, id string) string {
	redisType := redisType(data)
	return redisKeyWithTypeAndID(redisType, id)
}

func redisKeyWithTypeAndID(redisType, id string) string {
	return fmt.Sprintf("%s:%s", redisType, id)
}

func redisType(data interface{}) string {
	if typed, ok := data.(redisTyped); ok {
		return typed.RedisType()
	}

	return reflect.TypeOf(data).Elem().Name()
}

func redisID(data interface{}) string {
	if keyed, ok := data.(redisIDed); ok {
		return keyed.RedisID()
	}

	return reflect.ValueOf(data).Elem().FieldByName("ID").String()
}
