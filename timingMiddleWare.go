package main

import (
	"log"
	"net/http"
	"time"
)

func timeMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)

		duration := time.Since(startTime)

		log.Printf("%s took %v unit time ", r.URL.Path, duration)

	})

}
