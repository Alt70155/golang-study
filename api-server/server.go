package main

import (
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type User struct {
	ID        int    `json:"id"`
	Nickname  string `json:"nickname"`
	LoginName string `json:"loginname"`
	PassHash  string `json:"passhash"`
}

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/sample_db_go")
	defer db.Close()
	judgePanic(err)

	e := echo.New()
	user := show(db)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, user)
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
