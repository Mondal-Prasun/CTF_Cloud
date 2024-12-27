package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/docker/docker/client"
)

type DbConfig struct {
	sqlDb *sql.DB
}

type DockerConfig struct {
	dockerClinet *client.Client
}

func main() {
	fmt.Println("Hello ctf cloud...")

	mux := http.NewServeMux()

	db := initDataBase()
	defer db.Close()
	dbCfg := &DbConfig{
		sqlDb: db,
	}

	dcli := initDockerClient()
	defer dcli.Close()

	dockerCli := &DockerConfig{
		dockerClinet: dcli,
	}

	dockerCli.initilizeDockerImage()

	mux.HandleFunc("/health", healthCheck)
	mux.HandleFunc("/signIn", dbCfg.insertUserHandler)
	// mux.HandleFunc("/getUsers", dbCfg.getAllData)
	mux.HandleFunc("/logIn", dbCfg.getUserDetails)

	mux.HandleFunc("/createContainer", dockerCli.createDockerContainer)
	mux.HandleFunc("/startContainer", dockerCli.startDockerContainer)

	fmt.Println("The server is running on PORT:8080")
	if err := http.ListenAndServe(":8080", timeMiddleWare(mux)); err != nil {
		log.Fatal("Something went wrong while starting server...")
	}

}
