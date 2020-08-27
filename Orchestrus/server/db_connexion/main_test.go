package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

func compareWorkersList(workers1 []Worker, workers2 []Worker) bool {
	if len(workers1) != len(workers2) {
		return false
	}

	for i := range workers1 {
		if !compareWorker(workers1[i], workers2[i]) {
			return false
		}
	}

	return true
}

func compareWorker(worker1 Worker, worker2 Worker) bool {
	if worker1.IP != worker2.IP {
		return false
	}
	if worker1.Active != worker2.Active {
		return false
	}

	return compareImagesList(worker1.Images, worker2.Images)
}

func compareImagesList(images1 []Image, images2 []Image) bool {
	if len(images1) != len(images2) {
		return false
	}

	for i := range images1 {
		if !compareImage(images1[i], images2[i]) {
			return false
		}
	}
	return true
}

func compareImage(image1 Image, image2 Image) bool {
	if image1.ID != image2.ID {
		return false
	}
	if image1.IP != image2.IP {
		return false
	}
	if image1.Name != image2.Name {
		return false
	}

	if !cmp.Equal(image1.Ports, image2.Ports) {
		return false
	}

	return true
}

var dbConnexionURL = os.Getenv("DB_CONNEXION_URL")

func TestEverything(t *testing.T) {
	client := &http.Client{}
	// Connect to database
	var db *gorm.DB
	var connectionString = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		getEnv("POSTGRES_HOST", "postgres"),
		getEnv("POSTGRES_PORT", "5432"),
		getEnv("POSTGRES_USER", "postgres"),
		getEnv("POSTGRES_DB", "orchestrus"),
		getEnv("POSTGRES_PASSWORD", "postgres"),
	)

	db, err = gorm.Open("postgres", connectionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Test of adding a worker
	expectedAddedWorker := Worker{IP: "192.168.0.1", Active: true, Images: []Image{}}

	jsonValue, _ := json.Marshal(expectedAddedWorker)
	_, err := http.Post(dbConnexionURL+"/workers", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatal(err)
	}

	var foundAddedWorker Worker
	db.First(&foundAddedWorker, "ip = ?", "192.168.0.1")

	if !compareWorker(expectedAddedWorker, foundAddedWorker) {
		t.Fatal("POST /workers didn't return expected value")
	}

	// Test of modifying active status of worker
	expectedModifiedWorker := Worker{IP: "192.168.0.1", Active: false, Images: []Image{}}
	jsonValue, _ = json.Marshal(expectedModifiedWorker)

	putReq, err := http.NewRequest("PUT", dbConnexionURL+"/workers/192.168.0.1", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(putReq)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	var foundModifiedWorker Worker
	db.First(&foundModifiedWorker, "ip = ?", "192.168.0.1")

	if !compareWorker(expectedModifiedWorker, foundModifiedWorker) {
		t.Fatal("PUT /workers/{id} didn't return expected value")
	}

	// Test of listing workers
	resp, err = http.Get(dbConnexionURL + "/workers")
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var foundWorkersList []Worker
	json.Unmarshal(body, &foundWorkersList)

	expectedWorkersList := []Worker{Worker{IP: "192.168.0.1", Active: false}}

	if !compareWorkersList(expectedWorkersList, foundWorkersList) {
		t.Fatal("GET /workers didn't return expected value")
	}

	// Test of adding an image to a worker
	ports := json.RawMessage(`{"8080": "80"}`)
	var jsonbPorts postgres.Jsonb
	json.Unmarshal(ports, &jsonbPorts)

	expectedImage := Image{IP: "192.168.0.1", Name: "test", ID: "12345", Ports: jsonbPorts}
	jsonValue, _ = json.Marshal(expectedImage)

	_, err = http.Post(dbConnexionURL+"/workers/192.168.0.1/images", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatal(err)
	}

	var foundImage Image
	db.First(&foundImage, "id = ?", "12345")

	if !compareImage(expectedImage, foundImage) {
		t.Fatal("POST /workers/{ip}/images didn't return expected value")
	}

	// Test of removing an image from a worker
	deleteReq, err := http.NewRequest("DELETE", dbConnexionURL+"/workers/192.168.0.1/images/12345", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(deleteReq)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	var expectedDeletedImage Image
	var foundDeletedImage Image
	db.First(&foundDeletedImage, "id = ?", "12345")

	if !compareImage(foundDeletedImage, expectedDeletedImage) {
		t.Fatal("DELETE /workers/{ip}/images/{id} didn't return expected value")
	}
}
