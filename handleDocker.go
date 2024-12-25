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
	DEFAULT_IMAGE_TAG = "pwjcw/ctf"
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
	type ContainerDet struct {
		ContainerID string `json:"ctnID"`
		Message     string `json:"msg"`
	}

	if r.Method == http.MethodPost {

		containerChan := make(chan ContainerDet)

		go func(containrChan chan<- ContainerDet) {

			userData := struct {
				// Uuid uuid.UUID `json:"uid"`
				Name string `json:"name"`
			}{}

			decoder := json.NewDecoder(r.Body)

			if err := decoder.Decode(&userData); err != nil {
				responseWithError(w, 300, "json is not valid..")
				return
			}

			///TODO: change [userData.Name] to [userData.Uuid]
			createdContainer, err := dockerCLi.dockerClinet.ContainerCreate(context.Background(), &container.Config{
				Image: DEFAULT_IMAGE_TAG,
				Cmd: []string{
					"tail", "-f", "/dev/null",
				},
			}, &container.HostConfig{
				RestartPolicy: container.RestartPolicy{
					Name: container.RestartPolicyUnlessStopped,
				},
			},
				nil, nil, userData.Name)

			if err != nil {
				responseWithError(w, 501, fmt.Sprintf("Cannot create container : %s", err.Error()))
				return
			}

			containrChan <- ContainerDet{ContainerID: createdContainer.ID, Message: "Container created successFully"}

			log.Println("Created :", createdContainer.ID)

		}(containerChan)

		res := <-containerChan

		responseWithJson(w, 202, res)
	}
}

func (dockerCLi *DockerConfig) startDockerContainer(w http.ResponseWriter, r *http.Request) {
	type ContainerStartDet struct {
		Start   bool   `json:"cntStarted"`
		Message string `json:"msg"`
	}

	if r.Method == http.MethodPost {
		containerStartChan := make(chan ContainerStartDet)
		go func(startChan chan<- ContainerStartDet) {
			data := struct {
				CntId string `json:"containerId"`
			}{}

			dec := json.NewDecoder(r.Body)

			err := dec.Decode(&data)

			if err != nil {
				startChan <- ContainerStartDet{Start: false, Message: err.Error()}
				return
			}

			err = dockerCLi.dockerClinet.ContainerStart(context.Background(), data.CntId, container.StartOptions{})

			if err != nil {
				startChan <- ContainerStartDet{Start: false, Message: fmt.Sprintf("Cannot start container: %s", err.Error())}
				return
			}

			startChan <- ContainerStartDet{Start: true, Message: "Container started..."}
		}(containerStartChan)

		start := <-containerStartChan

		responseWithJson(w, 202, start)

	}
}

// func startDockerContainer2(w http.ResponseWriter, r *http.Request) {

// 	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

// 	if err != nil {
// 		log.Println(err.Error())
// 		return
// 	}

// 	crCtn, err := cli.ContainerCreate(context.Background(), &container.Config{
// 		Image: DEFAULT_IMAGE_TAG,
// 		Cmd: []string{
// 			"tail", "-f", "/dev/null",
// 		},
// 	}, &container.HostConfig{
// 		RestartPolicy: container.RestartPolicy{
// 			Name: container.RestartPolicyUnlessStopped,
// 		},
// 	}, nil, nil, "test")

// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	if err := cli.ContainerStart(context.Background(), crCtn.ID, container.StartOptions{}); err != nil {
// 		log.Println(err.Error())
// 	}

// 	log.Println("Test container strated")

// }
