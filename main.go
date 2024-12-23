package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

type dbConfig struct {
	sqlDb *sql.DB
}

func main() {
	fmt.Println("Hello ctf cloud...")

	mux := http.NewServeMux()

	db := initDataBase()

	dbCfg := dbConfig{
		sqlDb: db,
	}

	defer db.Close()

	mux.HandleFunc("/health", healthCheck)
	mux.HandleFunc("/addUser", dbCfg.insertUserHandler)
	mux.HandleFunc("/getUsers", dbCfg.getAllData)

	fmt.Println("The server is running on PORT:8080")
	if err := http.ListenAndServe(":8080", timeMiddleWare(mux)); err != nil {
		log.Fatal("Something went wrong while starting server...")
	}

}
