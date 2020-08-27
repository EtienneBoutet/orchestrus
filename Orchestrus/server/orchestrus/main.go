package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

var dbConnexionURL = getEnv("DB_CONNEXION_URL", "http://localhost:1234")
var hosts []string

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

func listWorkers(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(dbConnexionURL + "/workers")
	if err != nil {
		http.Error(w, "The DB module can't be contacted.", http.StatusInternalServerError)
		return
	}

	reqBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	w.Write(reqBody)
}

func addWorker(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var worker Worker
	err = json.Unmarshal(reqBody, &worker)
	if err != nil {
		panic(err)
	}

	// Check if worker is active
	var isValid = false
	resp, err := http.Get("http://" + worker.IP + ":5000/")
	if err != nil {
		http.Error(w, "The host can't be contacted.", http.StatusInternalServerError)
		return
	}
	if resp.StatusCode == 204 {
		isValid = true
	}

	// Add worker IP to the running hosts
	hosts = append(hosts, worker.IP)
	worker.Active = isValid

	// Add worker to DB
	jsonValue, err := json.Marshal(worker)
	if err != nil {
		panic(err)
	}

	resp, err = http.Post(dbConnexionURL+"/workers", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		http.Error(w, "The DB module can't be contacted.", http.StatusInternalServerError)
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	w.Write(respBody)
}

func startImage(w http.ResponseWriter, r *http.Request) {
	if len(hosts) == 0 {
		http.Error(w, "There is no running hosts", http.StatusInternalServerError)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var image Image
	err = json.Unmarshal(reqBody, &image)
	if err != nil {
		panic(err)
	}

	// Choose random host
	randomIndex := rand.Intn(len(hosts))
	host := hosts[randomIndex]
	image.IP = host

	// Send image to host worker
	jsonValue, err := json.Marshal(image)
	if err != nil {
		panic(err)
	}
    fmt.Println(image)
	resp, err := http.Post("http://"+image.IP+":5000/images", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		http.Error(w, "The host can't be contacted", http.StatusInternalServerError)
		return
	}

	// Extract returned image ID
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	image.ID = string(body)

    fmt.Println(string(body))
	// Add image to DB
	jsonValue, err = json.Marshal(image)
	if err != nil {
		panic(err)
	}
	resp, err = http.Post(dbConnexionURL+"/workers/"+host+"/images", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		http.Error(w, "The DB module can't be contacted", http.StatusInternalServerError)
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	w.Write(respBody)
}

func stopImage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ip := params["ip"]
	id := params["id"]

	client := &http.Client{}

	// Stop image of worker
	req, err := http.NewRequest("DELETE", "http://"+ip+":5000/images/"+id, nil)
	if err != nil {
		http.Error(w, "The host can't be contacted", http.StatusInternalServerError)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Remove image from DB
	req, err = http.NewRequest("DELETE", dbConnexionURL+"/workers/"+ip+"/images/"+id, nil)
	if err != nil {
		panic(err)
	}

	resp, err = client.Do(req)
	if err != nil {
		http.Error(w, "The DB module can't be contacted", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(204)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/workers", listWorkers).Methods("GET")
	router.HandleFunc("/workers", addWorker).Methods("POST")
	router.HandleFunc("/images", startImage).Methods("POST")
	router.HandleFunc("/workers/{ip}/images/{id}", stopImage).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":1235",
		handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(router)))
}

func main() {
	fmt.Println("Orchestrus API started")
	handleRequests()
}
