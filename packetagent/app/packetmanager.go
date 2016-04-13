package app

import (
	"github.com/golang/glog"
)

type Packetmanager struct {
}

func (self *Packetmanager) Start() {
	glog.Info("start the packetmanager")
}

func (self *Packetmanager) Stop() {

}
