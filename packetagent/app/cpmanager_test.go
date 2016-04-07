package app

import (
	"fmt"
	"testing"
	"time"
)

func TestProcessinfo(t *testing.T) {
	t.SkipNow()
	c, err := Newcpmanager()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c)
}

func TestStart(t *testing.T) {
	//test the start and the stop
	manager, err := Newcpmanager()
	if err != nil {
		fmt.Println(err)
	}
	manager.Start()
	time.Sleep(time.Second * 10)
	manager.Stop()
	time.Sleep(time.Second * 10)
}
