package app

import (
	"fmt"
	"github.com/cpmonitor/packetagent/lib"
	//"github.com/fsouza/go-dockerclient"
	"github.com/golang/glog"
	"io/ioutil"
	"regexp"
	"runtime"
	"strconv"
)

type Hostmetricmanager struct {
}

var (
	cpuRegExp  = regexp.MustCompile(`^processor\s*:\s*([0-9]+)$`)
	coreRegExp = regexp.MustCompile(`^core id\s*:\s*([0-9]+)$`)
	nodeRegExp = regexp.MustCompile(`^physical id\s*:\s*([0-9]+)$`)
	// Power systems have a different format so cater for both
	cpuClockSpeedMHz     = regexp.MustCompile(`(?:cpu MHz|clock)\s*:\s*([0-9]+\.[0-9]+)(?:MHz)?`)
	memoryCapacityRegexp = regexp.MustCompile(`MemTotal:\s*([0-9]+) kB`)
	swapCapacityRegexp   = regexp.MustCompile(`SwapTotal:\s*([0-9]+) kB`)
)

const (
	blockDir = "/sys/block"
	cacheDir = "/sys/devices/system/cpu/cpu"
	//all net device info was stored in this dir
	netDir       = "/sys/class/net"
	dmiDir       = "/sys/class/dmi"
	ppcDevTree   = "/proc/device-tree"
	s390xDevTree = "/etc" // s390/s390x changes
)

/*
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

	the machine info do not need to store in the influxdb

*/

func (self *Hostmetricmanager) Start() {
	glog.Info("start collecting the hostmetrics")
}

func (self *Hostmetricmanager) getOsversion() (string, error) {
	//on ubuntu search /etc/issue of /proc/version
	out, err := ioutil.ReadFile("/proc/version")
	if err != nil {
		return "", err
	}
	osinfo := string(out)
	return osinfo, nil
}

func (self *Hostmetricmanager) getDockerinfo() ([]string, error) {
	endpoint := "unix:///var/run/docker.sock"
	dockerclient, err := lib.Getdockerclient(endpoint)
	info, err := dockerclient.Getinfo()
	if err != nil {
		return nil, err
	}
	return info, nil
}

//the actually core of the progress is the runtime.GOMAXPROCS(0)
func (self *Hostmetricmanager) getNumcores() (int, error) {
	//maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	//if maxProcs < numCPU {
	//return maxProcs, nil
	//}
	return numCPU, nil

}

// unit MB
func (self *Hostmetricmanager) getMemoryInfo() (int64, error) {
	out, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, err
	}

	memoryCapacity, err := parseCapacity(out, memoryCapacityRegexp)
	if err != nil {
		return 0, err
	}
	memtotal_m := memoryCapacity / 1024
	return memtotal_m, err
}

// parseCapacity matches a Regexp in a []byte, returning the resulting value in bytes.
// Assumes that the value matched by the Regexp is in KB.
func parseCapacity(b []byte, r *regexp.Regexp) (int64, error) {
	matches := r.FindSubmatch(b)
	if len(matches) != 2 {
		return -1, fmt.Errorf("failed to match regexp in output: %q", string(b))
	}
	m, err := strconv.ParseInt(string(matches[1]), 10, 64)
	if err != nil {
		return -1, err
	}

	// Convert to bytes.
	return m * 1024, err
}

func (self *Hostmetricmanager) getNetworkDevices() (string, error) {
	return "", nil
}
