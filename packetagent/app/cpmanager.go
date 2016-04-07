package app

import (
	"github.com/golang/glog"
	"time"
)

type CPManager struct {
	Packetmanager          Packetmanager
	Hostmetricmanager      Hostmetricmanager
	Containermetricmanager Containermetricmanager
	Interval               time.Duration
	Quitglobalchannel      chan error //should be an array at last
}

func Newcpmanager() (CPManager, error) {
	packetmanager := Packetmanager{}
	hostmetricmanager := Hostmetricmanager{}
	containermanamge := Containermetricmanager{}
	cpmanager := CPManager{
		Packetmanager:          packetmanager,
		Hostmetricmanager:      hostmetricmanager,
		Containermetricmanager: containermanamge,
		//the frequency of collection  default 1s
		Interval:          time.Duration(time.Second),
		Quitglobalchannel: make(chan error),
	}
	return cpmanager, nil

}

func (self *CPManager) Start() {
	//start the intervals

	//start the Packetmanager start the housekeeping
	//start other components
	//quitGlobalHousekeeping := make(chan error)
	self.Packetmanager.Start()
	go self.globalhousekeeping(self.Quitglobalchannel)
}

func (self *CPManager) globalhousekeeping(quit chan error) {
	ticker := time.Tick(self.Interval)
	for {
		select {
		case t := <-ticker:
			start := time.Now()
			glog.Info("the time", start)
			glog.Info("unix time", t.Unix())
			//start collecting the data
			//send the heartbeat to the service
		case <-quit:
			//stop collecting
			glog.Info("Accept the signal , existing the current thread")
			return
		}
	}
	return
}

//send the kill signal to the manager
func (self *CPManager) Stop() error {
	self.Quitglobalchannel <- nil
	return nil
}

func (self *CPManager) detectcontainer() error {
	return nil
}

func (self *CPManager) detectsubcontainer() error {
	return nil
}
