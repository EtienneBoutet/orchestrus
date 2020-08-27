package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

var workerURL = os.Getenv("WORKER_URL")
var httpdURL = os.Getenv("HTTPD_URL")

func CheckHTTPDStarted() bool {
	for i := 1; i <= 20; i++ {
		resp, err := http.Get(httpdURL)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				return true
			}
		}
		time.Sleep(100 * time.Duration(i) * time.Millisecond)
	}
	return false
}

func TestEverything(t *testing.T) {
	// Check the status works
	resp, err := http.Get(workerURL + "/")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 204 {
		t.Fatal("Worker didn't reply with 204 for GET /")
	}
	// Start HTTPD
	jsonValue, _ := json.Marshal(Image{Name: "httpd", Ports: map[string]string{"8080": "80"}})
	resp, err = http.Post(workerURL+"/images", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal("Worker didn't reply with 200: " + string(body))
	}
	imageID := string(body)
	if imageID == "" || strings.Contains(imageID, " ") {
		t.Fatal("Worker didn't reply with an image ID: " + imageID)
	}
	// Check HTTPD is started
	httpdStarted := CheckHTTPDStarted()
	if !httpdStarted {
		t.Fatal("HTTPD didn't start in time")
	}
	// Stop HTTPD
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", workerURL+"/images/"+imageID, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 204 {
		t.Fatal("Worker didn't reply with 204 for DELETE /images/{id}")
	}
}
