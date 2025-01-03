package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func responseWithJson(w http.ResponseWriter, statusCode int, payload interface{}) {

	data, err := json.Marshal(payload)

	if err != nil {
		fmt.Println("Cannot marshal the payload to json..")
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)

}

func responseWithError(w http.ResponseWriter, statuscode int, msg string) {

	data := struct {
		Error string `json:"err"`
	}{
		Error: msg,
	}

	responseWithJson(w, statuscode, data)
}
