package main

import (
	"github.com/gomodule/redigo/redis"
)

func Connection() redis.Conn {
	const Addr = "127.0.0.1:6379"

	c, err := redis.Dial("tcp", Addr)
	judgePanic(err)

	return c
}

// データの登録（Redis: SET key value）
func Set(key, value string, c redis.Conn) string {
	res, err := redis.String(c.Do("SET", key, value))
	judgePanic(err)

	return res
}

func judgePanic(err error) {
	if err != nil {
		panic(err)
	}
}
