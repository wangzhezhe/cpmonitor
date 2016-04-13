package lib

import (
	"github.com/fsouza/go-dockerclient"
)

type Dockerclient struct {
	Docker_client *docker.Client
}

func Getdockerclient(endpoint string) (*Dockerclient, error) {
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return nil, err
	}
	dockerclient := &Dockerclient{Docker_client: client}
	return dockerclient, nil
}

//use docker version api and the docker info api
func (self *Dockerclient) Getinfo() ([]string, error) {
	dockerinfo := []string{""}
	info, err := self.Docker_client.Info()
	version, err := self.Docker_client.Version()

	dockerinfo = append(dockerinfo, []string(*info)...)
	dockerinfo = append(dockerinfo, []string(*version)...)
	//infotest := append(dockerinfo, "abcd")
	if err != nil {
		return nil, err
	}
	return dockerinfo, nil
}
