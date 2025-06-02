package main

import (
	"learning_go/internal/router"
	"log"
	"net/http"
)

func main() {
	r := router.New()
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
