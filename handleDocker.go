package main

import (
	"context"
	"fmt"
	"log"

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
