package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func initDataBase() (database *sql.DB) {
	db, err := sql.Open("sqlite3", "test.db")

	if err != nil {
		log.Fatal("Cannot create or connect to sqlite3 driver: ", err)
	}

	createUser := `CREATE TABLE IF NOT EXISTS Users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER
	);`
	_, err = db.Exec(createUser)

	if err != nil {
		log.Fatal("Cannot create table: ", err)
	}
	log.Println("Conneted to the sqlite3 database....")

	return db
}

func insertUser(db *sql.DB, user *User) (res interface{}, error error) {
	inUser := `INSERT INTO users (name, age) VALUES (?, ?)`
	res, err := db.Exec(inUser, user.Name, user.Age)
	if err != nil {
		return "", err
	} else {
		return res, nil
	}

}

func getUsers(db *sql.DB) (rows interface{}, error error) {
	selectUser := `SELECT * from Users`

	res, err := db.Query(selectUser)

	type Data struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	if err != nil {
		return "", err
	} else {
		var alldata []Data
		for res.Next() {
			var id, age int
			var name string

			if err := res.Scan(&id, &name, &age); err != nil {
				log.Println("Cannot find all the data")
				return
			}
			alldata = append(alldata, Data{
				Id:   id,
				Name: name,
				Age:  age,
			})

		}

		return alldata, nil
	}

}
