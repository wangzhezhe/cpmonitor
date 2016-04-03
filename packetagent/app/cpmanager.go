package app

import ()

type CPManager struct {
	Packetmanager          Packetmanager
	Hostmetricmanager      Hostmetricmanager
	Containermetricmanager Containermetricmanager
}

func Newcpmanager() (CPManager, error) {
	packetmanager := Packetmanager{}
	hostmetricmanager := Hostmetricmanager{}
	containermanamge := Containermetricmanager{}
	cpmanager := CPManager{
		Packetmanager:          packetmanager,
		Hostmetricmanager:      hostmetricmanager,
		Containermetricmanager: containermanamge,
	}
	return cpmanager, nil

}

func (self *CPManager) Start() error {
	//start the intervals

	//start the Packetmanager
	self.Packetmanager.Start()
	return nil
}

func (self *CPManager) Stop() error {
	return nil
}
