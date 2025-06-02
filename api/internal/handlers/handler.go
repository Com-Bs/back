package handler

import (
	"log"
	"net/http"
)

func signUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func logIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func logOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// Default GET function
func GetLogs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("GetingLogs")
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// Default GET function
func GetFullCompile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
