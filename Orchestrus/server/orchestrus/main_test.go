package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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

func compareWorker(worker1 Worker, worker2 Worker) bool {
	if worker1.IP != worker2.IP {
		return false
	}
	if worker1.Active != worker2.Active {
		return false
	}

	return compareImagesList(worker1.Images, worker2.Images)
}

var serverURL = os.Getenv("SERVER_URL")

func TestEverything(t *testing.T) {
	cmd := exec.Command("./test_script.sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	duration := time.Duration(3) * time.Second
	time.Sleep(duration)

	// Test of adding a worker
	expectedWorker := Worker{IP: "localhost", Active: true, Images: []Image{}}
	jsonValue, _ := json.Marshal(expectedWorker)

	resp, err := http.Post(serverURL+"/workers", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatal(err)
	}

	reqBody, _ := ioutil.ReadAll(resp.Body)
	var foundWorker Worker
	err = json.Unmarshal(reqBody, &foundWorker)
	if err != nil {
		t.Fatal(err)
	}

	if !compareWorker(expectedWorker, foundWorker) {
		t.Fatal("addWorker didn't return expected value")
	}

	// Test of getting workers
	resp, err = http.Get(serverURL + "/workers")
	if err != nil {
		t.Fatal(err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var foundWorkers []Worker
	json.Unmarshal(respBody, &foundWorkers)

	expectedWorkers := []Worker{Worker{IP: "localhost", Active: true, Images: []Image{}}}

	if !compareWorkersList(expectedWorkers, foundWorkers) {
		t.Fatal("listWorkers didn't return expected value")
	}

	// Test of adding an image to a worker
	expectedImage := Image{ID: "42", IP: "localhost", Name: "httpd", Ports: map[string]string{"8080": "80"}}
	jsonValue, _ = json.Marshal(expectedImage)

	resp, err = http.Post(serverURL+"/images", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatal(err)
	}

	reqBody, _ = ioutil.ReadAll(resp.Body)
	var foundImage Image
	err = json.Unmarshal(reqBody, &foundImage)
	if err != nil {
		t.Fatal(err)
	}

	if !compareImage(expectedImage, foundImage) {
		t.Fatal("startImage didn't return expected value")
	}

	// Test of deleting an image of a worker
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", serverURL+"/workers/localhost/images/42", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if resp.StatusCode != 204 {
		t.Fatal("stopImage didn't return expected value")
	}

	log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		syscall.Kill(-pgid, 15)
	}

	cmd.Wait()
}
