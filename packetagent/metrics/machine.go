package metrics

import ()

type Machineinfo struct {
	//OS version
	Osversion string

	// Docker version.
	DockerVersion string `json:"docker_version"`

	// The number of cores in this machine.
	NumCores int `json:"num_cores"`

	// The amount of memory (in bytes) in this machine
	MemoryCapacity uint64 `json:"memory_capacity"`

	// Network devices
	NetworkDevices []NetInfo `json:"network_devices"`
}
