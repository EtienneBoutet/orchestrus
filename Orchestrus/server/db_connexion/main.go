package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Image struct {
	ID    string         `gorm:"PRIMARY_KEY;AUTO_INCREMENT=FALSE" json:"id"`
	IP    string         `gorm:"PRIMARY_KEY;AUTO_INCREMENT=FALSE" json:"ip"`
	Name  string         `gorm:"NOT NULL" json:"name"`
	Ports postgres.Jsonb `gorm:"NOT NULL" json:"ports"`
}

type Worker struct {
	IP     string  `gorm:"PRIMARY_KEY;AUTO_INCREMENT=FALSE" json:"ip"`
	Active bool    `gorm:"NOT NULL" json:"active"`
	Images []Image `gorm:"FOREIGNKEY:ip" json:"images"`
}

var db *gorm.DB
var err error

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func getWorkers(w http.ResponseWriter, r *http.Request) {
	var workers []Worker
	db.Find(&workers)
	for i := 0; i < len(workers); i++ {
		var images []Image
		db.Where("ip = ?", workers[i].IP).Find(&images)
		workers[i].Images = images
	}
	json.NewEncoder(w).Encode(&workers)
}

func addWorker(w http.ResponseWriter, r *http.Request) {
	var worker Worker
	json.NewDecoder(r.Body).Decode(&worker)
	db.Create(&worker)
	json.NewEncoder(w).Encode(&worker)
}

func updateWorker(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var worker Worker
	db.Where("ip = ?", params["ip"]).First(&worker)
	json.NewDecoder(r.Body).Decode(&worker)
	db.Save(&worker)
	json.NewEncoder(w).Encode(&worker)
}

func getImages(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var images []Image
	db.Where("ip = ?", params["ip"]).Find(&images)
	json.NewEncoder(w).Encode(&images)
}

func addImage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var image Image
	json.NewDecoder(r.Body).Decode(&image)
	image.IP = params["ip"]
	db.Create(&image)
	json.NewEncoder(w).Encode(&image)
}

func deleteImage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var image Image
	db.Where("ip = ? AND id = ?", params["ip"], params["id"]).First(&image)
	db.Delete(&image)
	json.NewEncoder(w).Encode(&image)
}

func handleRequests() {
	router := mux.NewRouter()

	router.HandleFunc("/workers", getWorkers).Methods("GET")
	router.HandleFunc("/workers", addWorker).Methods("POST")
	router.HandleFunc("/workers/{ip}", updateWorker).Methods("PUT")

	router.HandleFunc("/workers/{ip}/images", getImages).Methods("GET")
	router.HandleFunc("/workers/{ip}/images", addImage).Methods("POST")
	router.HandleFunc("/workers/{ip}/images/{id}", deleteImage).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":"+getEnv("PORT", "1234"), router))
}

func connectDB(connectionString string) *gorm.DB {
	for i := 1; i <= 10; i++ {
		db, err = gorm.Open("postgres", connectionString)
		if err != nil {
			time.Sleep(100 * time.Duration(i) * time.Millisecond)
		} else {
			return db
		}
	}
	return nil
}

func main() {
	var connectionString = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		getEnv("POSTGRES_HOST", "127.0.0.1"),
		getEnv("POSTGRES_PORT", "5432"),
		getEnv("POSTGRES_USER", "postgres"),
		getEnv("POSTGRES_DB", "orchestrus"),
		getEnv("POSTGRES_PASSWORD", "postgres"),
	)
	db := connectDB(connectionString)

	if db == nil {
		panic("failed to connect database")
	}

	defer db.Close()
	db.AutoMigrate(&Worker{}, &Image{})
	fmt.Println("Db_connexion API started")
	handleRequests()
}
