package main

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const (
	MAIN_DATABASE_NAME = "user.db"
)

type User struct {
	Uid      uuid.UUID `json:"uid"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
}

func initDataBase() (database *sql.DB) {
	db, err := sql.Open("sqlite3", MAIN_DATABASE_NAME)

	if err != nil {
		log.Fatal("Cannot create or connect to sqlite3 driver: ", err)
	}

	createUserTable := `CREATE TABLE IF NOT EXISTS Users (
		id TEXT PRIMARY KEY NOT NULL,
        name TEXT NOT NULL UNIQUE,          
        password TEXT NOT NULL 
	);`
	_, err = db.Exec(createUserTable)

	if err != nil {
		log.Fatal("Cannot create table: ", err)
	}
	log.Println("Conneted to the sqlite3 database....")

	return db
}

func insertUser(db *sql.DB, user *User) (res sql.Result, error error) {
	inUser := `INSERT INTO Users (id, name, password) VALUES (?, ?, ?)`
	data, err := db.Exec(inUser, user.Uid, user.Name, user.Password)

	if err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func getUserDetails(db *sql.DB, uName string) (userId string, userName string, passwords string, error error) {
	selectUser := `SELECT id AS uid,name,password FROM Users WHERE name = ?`

	data := db.QueryRow(selectUser, uName)

	var id string
	var name string
	var password string

	err := data.Scan(&id, &name, &password)

	if err != nil {
		return "", "", "", err
	} else {
		return id, name, password, nil
	}

}

// func getUsers(db *sql.DB) (rows interface{}, error error) {
// 	selectUser := `SELECT * from Users`

// 	res, err := db.Query(selectUser)

// 	type Data struct {
// 		Id   int    `json:"id"`
// 		Name string `json:"name"`
// 		Age  int    `json:"age"`
// 	}

// 	if err != nil {
// 		return "", err
// 	} else {
// 		var alldata []Data
// 		for res.Next() {
// 			var id, age int
// 			var name string

// 			if err := res.Scan(&id, &name, &age); err != nil {
// 				log.Println("Cannot find all the data")
// 				return
// 			}
// 			alldata = append(alldata, Data{
// 				Id:   id,
// 				Name: name,
// 				Age:  age,
// 			})

// 		}

// 		return alldata, nil
// 	}

// }
