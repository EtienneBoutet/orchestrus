package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Image .
type Image struct {
	ID    string            `json:"id"`
	IP    string            `json:"ip"`
	Name  string            `json:"name"`
	Ports map[string]string `json:"ports"`
}

// Worker .
type Worker struct {
	IP     string  `json:"ip"`
	Active bool    `json:"active"`
	Images []Image `json:"images"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func startImage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "42")
}

func stopImage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage)
	router.HandleFunc("/images", startImage).Methods("POST")
	router.HandleFunc("/images/{id}", stopImage).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":5000", router))
}

func main() {
	handleRequests()
}
