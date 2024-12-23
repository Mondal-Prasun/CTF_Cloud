package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Hello ctf cloud...")

	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthCheck)

	fmt.Println("The server is running on PORT:8080")
	if err := http.ListenAndServe(":8080", timeMiddleWare(mux)); err != nil {
		log.Fatal("Something went wrong while starting server...")
	}

}
