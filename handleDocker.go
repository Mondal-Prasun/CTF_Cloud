package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

const (
	DEFAULT_IMAGE_TAG = "ubuntu:latest"
)

func initDockerClient() *client.Client {
	log.Println("Initilizing docker")

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		log.Fatal("Cannot initilize Docker Client")
	}

	log.Println("Docker client initilized.....")
	return dockerClient
}

func (dockerCLi *DockerConfig) initilizeDockerImage() {
	allImages, err := dockerCLi.dockerClinet.ImageList(context.Background(), image.ListOptions{})

	if err != nil {
		log.Println("Cannot fetch container images..")
		return
	}

	defaultImageExcist := false

	for _, image := range allImages {
		for _, tag := range image.RepoTags {
			if tag == DEFAULT_IMAGE_TAG {
				defaultImageExcist = true
				fmt.Printf("%s :image already exsist\n", DEFAULT_IMAGE_TAG)

			}
		}
	}

	if !defaultImageExcist {
		_, err := dockerCLi.dockerClinet.ImagePull(context.Background(), DEFAULT_IMAGE_TAG, image.PullOptions{})

		if err != nil {
			log.Fatal("Cannot pull image : ", DEFAULT_IMAGE_TAG)
		}

		log.Printf("%s is intilized... \n", DEFAULT_IMAGE_TAG)
	}

}

func (dockerCLi *DockerConfig) createDockerContainer(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		userData := struct {
			// Uuid uuid.UUID `json:"uid"`
			Name string `json:"name"`
		}{}

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&userData); err != nil {
			responseWithError(w, 300, "json is not valid..")
			return
		}

		containerConfig := &container.Config{
			Image: DEFAULT_IMAGE_TAG,
		}
		///TODO: change [userData.Name] to [userData.Uuid]
		createdContainer, err := dockerCLi.dockerClinet.ContainerCreate(context.Background(), containerConfig, nil, nil, nil, userData.Name)

		if err != nil {
			responseWithError(w, 501, fmt.Sprintf("Cannot create container : %s", err.Error()))
			return
		}

		log.Println("Created :", createdContainer.ID)

		responseWithJson(w, 202, struct {
			ContId string `json:"Id"`
		}{
			ContId: createdContainer.ID,
		})

	}
}
