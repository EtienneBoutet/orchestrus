package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gorilla/mux"
)

// Image is the params given to startImage
type Image struct {
	Name  string            `json:"name"`
	Ports map[string]string `json:"ports"`
}

// Started is a map of docker container ids to Image started by this server
var Started map[string]Image

// Docker is the global docker client
var Docker *client.Client

func canonizeImageName(name string) string {
	if !strings.Contains(name, "/") {
		return "docker.io/library/" + name
	}
	return "docker.io" + name
}

func mapPorts(ports map[string]string) nat.PortMap {
	portMap := nat.PortMap{}

	for hostPort, containerPort := range ports {
		hostBinding := nat.PortBinding{
			HostIP:   "",
			HostPort: hostPort,
		}
		containerPort, err := nat.NewPort("tcp", containerPort)
		if err != nil {
			panic("Container port is invalid")
		}

		portMap[containerPort] = []nat.PortBinding{hostBinding}
	}

	return portMap
}

func startContainer(image Image) string {
	ctx := context.Background()
	imageName := canonizeImageName(image.Name)

	// Pull the image
	out, err := Docker.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()
	_, err = ioutil.ReadAll(out)
	if err != nil {
		panic(err)
	}

	// Create the container
	resp, err := Docker.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageName,
		},
		&container.HostConfig{
			PortBindings: mapPorts(image.Ports),
		}, nil, "")
	if err != nil {
		panic(err)
	}

	// Start the container
	err = Docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		stopContainer(resp.ID)
		panic(err)
	}

	return resp.ID
}

func stopContainer(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Docker.ContainerStop(ctx, id, nil)
	if err != nil {
		panic("Failed to stop the container '" + id[:10] + "'!")
	}

	err = Docker.ContainerRemove(ctx, id, types.ContainerRemoveOptions{})
	if err != nil {
		panic("Failed to remove the container '" + id[:10] + "'!")
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func listImages(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Started)
}

func startImage(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var image Image
	json.Unmarshal(reqBody, &image)

	// Capture panics and return a 500
	defer func() {
		if r := recover(); r != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, r)
		}
	}()
	fmt.Printf("Starting %+v\n", image)
	id := startContainer(image)

	Started[id] = image
	fmt.Fprintf(w, id)
}

func stopImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, found := Started[id]; found {
		delete(Started, id)
		stopContainer(id)
		fmt.Println("Stopped the container " + id[:10])
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage)
	router.HandleFunc("/images", listImages).Methods("GET")
	router.HandleFunc("/images", startImage).Methods("POST")
	router.HandleFunc("/images/{id}", stopImage).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":5000", router))
}

func main() {
	Started = make(map[string]Image)
	client, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}
	Docker = client

	fmt.Println("Worker Rest API started")
	handleRequests()
}
