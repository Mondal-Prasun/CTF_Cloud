package main

import (
	"net/http"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		data := struct {
			Message string `json:"msg"`
		}{
			Message: "This is a health check",
		}

		responseWithJson(w, 201, data)
	}

}
