package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID   int
	Name string
}

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/sample_db_go")
	defer db.Close()

	judgePanic(err)

	fmt.Println("Show all：")
	show(db)

	fmt.Println("Insert：")
	insert(db)

	fmt.Println("Update")
	update(db)
	show(db)

	fmt.Println("Delete")
	delete(db)
}

func judgePanic(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func show(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM users")

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name)

		judgePanic(err)

		fmt.Println(user.ID, user.Name)
	}

	err = rows.Err()

	judgePanic(err)
}

func insert(db *sql.DB) {
	stmtInsert, err := db.Prepare("INSERT INTO users(name) VALUES(?)")

	if err != nil {
		panic(err.Error())
	}

	defer stmtInsert.Close()

	result, err := stmtInsert.Exec("近　海斗")
	judgePanic(err)

	lastInsertID, err := result.LastInsertId()
	judgePanic(err)

	fmt.Println(lastInsertID)
}

func update(db *sql.DB) {
	stmtUpdate, err := db.Prepare("UPDATE users SET name=? WHERE id=?")
	judgePanic(err)

	defer stmtUpdate.Close()

	result, err := stmtUpdate.Exec("漱石", 1)
	judgePanic(err)

	rowsAffect, err := result.RowsAffected()
	judgePanic(err)

	fmt.Println(rowsAffect)
}

func delete(db *sql.DB) {
	stmtDelete, err := db.Prepare("DELETE FROM users WHERE id=?")
	judgePanic(err)

	defer stmtDelete.Close()

	result, err := stmtDelete.Exec(3)
	judgePanic(err)

	rowsAffect, err := result.RowsAffected()
	judgePanic(err)

	fmt.Println(rowsAffect)
}
