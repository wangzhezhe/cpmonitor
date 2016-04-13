package app

import (
	"fmt"
	"testing"
	//"time"
)

func TestGetMemoryInfo(t *testing.T) {

	hostmanager := &Hostmetricmanager{}

	memoinfo, err := hostmanager.getMemoryInfo()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("the memcapacity %d KB", memoinfo)
}

func TestGetNumcores(t *testing.T) {
	hostmanager := &Hostmetricmanager{}
	numcore, err := hostmanager.getNumcores()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("the nums of cores: %d", numcore)
}

func TestGetOsVersion(t *testing.T) {
	hostmanager := &Hostmetricmanager{}
	numcore, err := hostmanager.getOsversion()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("the os info %v:", numcore)
}

func TestGetDockerinfo(t *testing.T) {
	hostmanager := &Hostmetricmanager{}
	info, err := hostmanager.getDockerinfo()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("the docker info %v", info)
}
