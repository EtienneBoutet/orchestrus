package main

import (
	"encoding/json"
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

func getWorkers(w http.ResponseWriter, r *http.Request) {
	workers := []Worker{Worker{IP: "localhost", Active: true, Images: []Image{}}}
	json.NewEncoder(w).Encode(&workers)
}

func addWorker(w http.ResponseWriter, r *http.Request) {
	worker := Worker{IP: "localhost", Active: true, Images: []Image{}}
	json.NewEncoder(w).Encode(&worker)
}

func addImage(w http.ResponseWriter, r *http.Request) {
	image := Image{ID: "42", IP: "localhost", Name: "httpd", Ports: map[string]string{"8080": "80"}}
	json.NewEncoder(w).Encode(&image)
}

func deleteImage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func handleRequests() {
	router := mux.NewRouter()

	router.HandleFunc("/workers", getWorkers).Methods("GET")
	router.HandleFunc("/workers", addWorker).Methods("POST")
	router.HandleFunc("/workers/{ip}/images", addImage).Methods("POST")
	router.HandleFunc("/workers/{ip}/images/{id}", deleteImage).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":1234", router))
}

func main() {
	handleRequests()
}
