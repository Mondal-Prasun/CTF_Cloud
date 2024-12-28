package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

const (
	DEFAULT_IMAGE_TAG = "kali-vnc:latest"
	DEFAULT_CTF_TAG   = "bkimminich/juice-shop:latest"
)

type UserSession struct {
	KaliUrl string `json:"Kali-vnc-url"`
	CtfUrl  string `json:"ctf-url"`
}

var allSessions = make(map[string]UserSession)

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
	defaultCtfImageExcist := false

	for _, image := range allImages {
		for _, tag := range image.RepoTags {

			if tag == DEFAULT_CTF_TAG {
				defaultCtfImageExcist = true
				fmt.Printf("%s :image already exsist\n", DEFAULT_CTF_TAG)
			}

			if tag == DEFAULT_IMAGE_TAG {
				defaultImageExcist = true
				fmt.Printf("%s :image already exsist\n", DEFAULT_IMAGE_TAG)

			}
		}
	}

	if !defaultCtfImageExcist {
		_, err := dockerCLi.dockerClinet.ImagePull(context.Background(), DEFAULT_CTF_TAG, image.PullOptions{})

		if err != nil {
			log.Fatal("Cannot pull image : ", DEFAULT_IMAGE_TAG)
		}

		log.Printf("%s is intilized... \n", DEFAULT_IMAGE_TAG)
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
		KaliContainerID string `json:"kaliCtnId"`
		KaliUrl         string `json:"kaliUrl"`
		CtfContainerID  string `json:"ctfCtnId"`
		CtfUrl          string `json:"ctfUrl"`
		UserNetwork     string `json:"netId"`
		Message         string `json:"msg"`
		Created         bool   `json:"created"`
	}

	if r.Method == http.MethodPost {

		containerChan := make(chan ContainerDet)

		type Data struct {
			Uuid uuid.UUID `json:"uid"`
		}

		userData := Data{}

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&userData); err != nil {
			responseWithError(w, 300, err.Error())
			return
		}

		go func(containrChan chan<- ContainerDet, userData *Data) {

			// private network for indevidual user for isolation
			userNetWorkId := fmt.Sprintf("network-%s", userData.Uuid)
			userNet, err := dockerCLi.dockerClinet.NetworkCreate(context.Background(), userNetWorkId, network.CreateOptions{})

			kaliCntPort := fmt.Sprintf("%d", (8000 + len(allSessions)*2 + 1))
			ctfCntPort := fmt.Sprintf("%d", (6000 + len(allSessions)*2 + 1))

			if err != nil {
				containrChan <- ContainerDet{KaliContainerID: "", CtfContainerID: "", UserNetwork: "", Message: err.Error(), Created: false}
				return
			}

			///TODO: change [userData.Name] to [userData.Uuid]
			createdKaliContainer, err := dockerCLi.dockerClinet.ContainerCreate(context.Background(), &container.Config{
				Image: DEFAULT_IMAGE_TAG,
				Cmd:   []string{"sh", "-c", "dbus-launch vncserver :1 && websockify --web=/usr/share/novnc/ 6969 localhost:5901"},
				ExposedPorts: nat.PortSet{
					"6969/tcp": struct{}{},
				},
			}, &container.HostConfig{
				PortBindings: nat.PortMap{
					"6969/tcp": []nat.PortBinding{
						{
							HostIP:   "0.0.0.0",
							HostPort: kaliCntPort,
						},
					},
				},
				NetworkMode: container.NetworkMode(userNetWorkId),
				RestartPolicy: container.RestartPolicy{
					Name: container.RestartPolicyUnlessStopped,
				},
			},
				nil, nil,
				fmt.Sprintf("%s-kali", userData.Uuid))

			if err != nil {
				containrChan <- ContainerDet{KaliContainerID: "", CtfContainerID: "", UserNetwork: "", Message: err.Error(), Created: false}
				return
			}

			createdCtfContainer, err := dockerCLi.dockerClinet.ContainerCreate(context.Background(),
				&container.Config{
					Image: DEFAULT_CTF_TAG,
					ExposedPorts: nat.PortSet{
						"3000/tcp": struct{}{},
					},
				},
				&container.HostConfig{
					PortBindings: nat.PortMap{
						"3000/tcp": []nat.PortBinding{
							{
								HostIP:   "0.0.0.0",
								HostPort: ctfCntPort,
							},
						},
					},
					NetworkMode: container.NetworkMode(userNetWorkId),
				}, nil, nil,
				fmt.Sprintf("%s-ctf", userData.Uuid))

			if err != nil {
				dockerCLi.dockerClinet.ContainerRemove(context.Background(), createdKaliContainer.ID, container.RemoveOptions{
					Force: true,
				})
				containrChan <- ContainerDet{KaliContainerID: "", CtfContainerID: "", UserNetwork: "", Message: err.Error(), Created: false}
				return
			}

			containrChan <- ContainerDet{
				KaliContainerID: createdKaliContainer.ID,
				KaliUrl:         fmt.Sprintf("192.168.0.101:%s/vnc.html", kaliCntPort),
				CtfContainerID:  createdCtfContainer.ID,
				CtfUrl:          fmt.Sprintf("192.168.0.101:%s/", ctfCntPort),
				UserNetwork:     userNet.ID,
				Message:         "Containers created....",
				Created:         true}

			log.Println("Created :", createdKaliContainer.ID)

		}(containerChan, &userData)

		res := <-containerChan
		if res.Created {
			allSessions[userData.Uuid.String()] = UserSession{
				KaliUrl: res.KaliUrl,
				CtfUrl:  res.CtfUrl,
			}
			responseWithJson(w, 202, res)
		} else {
			responseWithError(w, 503, res.Message)
		}

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
				KaliCntId string `json:"kaliContainerId"`
				CtfCntId  string `json:"ctfContainerId"`
			}{}

			dec := json.NewDecoder(r.Body)

			err := dec.Decode(&data)

			if err != nil {
				startChan <- ContainerStartDet{Start: false, Message: err.Error()}
				return
			}

			err = dockerCLi.dockerClinet.ContainerStart(context.Background(), data.KaliCntId, container.StartOptions{})

			if err != nil {
				startChan <- ContainerStartDet{Start: false, Message: fmt.Sprintf("Cannot start container: %s", err.Error())}
				return
			}

			err = dockerCLi.dockerClinet.ContainerStart(context.Background(), data.CtfCntId, container.StartOptions{})

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

//MARK:Todo:terminate
