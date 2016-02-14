package redis

import (
	"fmt"
	"reflect"
)

type redisTyper interface {
	RedisType() string
}

type redisIDer interface {
	RedisID() string
}

type redisIndexer interface {
	RedisIndexes() map[string]string
}

func redisKeyWithTypeAndID(redisType, id string) string {
	return fmt.Sprintf("%s:%s", redisType, id)
}

func redisType(data interface{}) string {
	if typed, ok := data.(redisTyper); ok {
		return typed.RedisType()
	}

	return reflect.TypeOf(data).Elem().Name()
}

func redisID(data interface{}) string {
	if keyed, ok := data.(redisIDer); ok {
		return keyed.RedisID()
	}

	return reflect.ValueOf(data).Elem().FieldByName("ID").String()
}
