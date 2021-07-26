package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo"
	"go.knocknote.io/rapidash"
)

type User struct {
	ID        int    `json:"id"`
	Nickname  string `json:"nickname"`
	LoginName string `json:"loginname"`
	PassHash  string `json:"passhash"`
}

func (u *User) DecodeRapidash(decoder rapidash.Decoder) error {
	u.ID = decoder.Int("id")
	u.Nickname = decoder.String("nickname")
	u.LoginName = decoder.String("login_name")
	u.PassHash = decoder.String("pass_hash")
	return decoder.Error()
}

func (u User) RapidashStruct() *rapidash.Struct {
	return rapidash.NewStruct("users").
		FieldInt("id").
		FieldString("nickname").
		FieldString("login_name").
		FieldString("pass_hash")
}

type UserSlice []*User

func (e *UserSlice) DecodeRapidash(decoder rapidash.Decoder) error {
	*e = make(UserSlice, decoder.Len())
	for i := 0; i < decoder.Len(); i++ {
		var user User
		if err := user.DecodeRapidash(decoder.At(i)); err != nil {
			return err
		}
		(*e)[i] = &user
	}
	return decoder.Error()
}

func main() {
	// mysql
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/sample_db_go")
	defer db.Close()
	judgePanic(err)

	// rapidash
	// redis へのアドレスを指定して rapidash インスタンス作成
	cache, err := rapidash.New(rapidash.ServerAddrs([]string{"127.0.0.1:6379"}))
	judgePanic(err)
	err = cache.WarmUp(db, new(User).RapidashStruct(), true)
	judgePanic(err)
	// *sql.Tx インスタンスの作成
	txConn, err := db.Begin()
	judgePanic(err)
	cacheTx, err := cache.Begin(txConn)
	judgePanic(err)

	defer func() {
		if err := cacheTx.RollbackUnlessCommitted(); err != nil {
			log.Println(err)
		}
	}()

	var users UserSlice

	if err := cacheTx.FindByQueryBuilder(rapidash.NewQueryBuilder("users").Eq("id", 1), &users); err != nil {
		log.Println(err)
	}
	// Commitしてキャッシュを更新する
	if err := cacheTx.Commit(); err != nil {
		log.Println(err)
	}

	cacheTx.Commit()

	// api-server
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, users)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func judgePanic(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func show(db *sql.DB) User {
	id := 1
	var user User
	err := db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user.ID, &user.Nickname, &user.LoginName, &user.PassHash)
	judgePanic(err)

	return user
}

func ConnectionRedis() redis.Conn {
	const Addr = "127.0.0.1:6379"

	c, err := redis.Dial("tcp", Addr)
	judgePanic(err)

	return c
}
