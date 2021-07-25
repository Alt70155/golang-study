package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func judgePanic(err error) {
	if err != nil {
		panic(err)
	}
}

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

// データの取得（Redis: GET key）
func Get(key string, c redis.Conn) string {
	res, err := redis.String(c.Do("GET", key))
	judgePanic(err)

	return res
}

type Data struct {
	Key   string
	Value string
}

// 複数のデータの登録（Redis: MSET key [key...]）
func Mset(dataList []Data, c redis.Conn) {
	var query []interface{}

	for _, v := range dataList {
		query = append(query, v.Key, v.Value)
	}

	fmt.Println(query)

	c.Do("MSET", query...)
}

// 複数の値を取得（Redis: MGET key [key...]）
func Mget(keys []string, c redis.Conn) []string {
	var query []interface{}

	for _, v := range keys {
		query = append(query, v)
	}

	fmt.Println("MGET query:", query) // [key1, key2]

	res, err := redis.Strings(c.Do("MGET", query...))
	judgePanic(err)

	return res
}

func Expire(key string, ttl int, c redis.Conn) {
	c.Do("EXPIRE", key, ttl)
}

func main() {
	// 接続
	c := Connection()
	defer c.Close()

	// データの登録
	res_set := Set("sample-key", "sample-value", c)
	fmt.Println(res_set) // OK

	// データの取得
	res_get := Get("sample-key", c)
	fmt.Println(res_get) // sample-value

	// 複数データの登録
	dataList := []Data{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	}

	Mset(dataList, c)

	// 複数データの取得
	keys := []string{"key1", "key2"}
	res_mget := Mget(keys, c)
	fmt.Println(res_mget)

	// TTLの設定
	Expire("key1", 10, c)
}
