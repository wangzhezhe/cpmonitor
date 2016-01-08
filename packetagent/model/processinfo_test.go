package model

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"testing"
)

func TestProcessinfo(t *testing.T) {
	Pdetail := &ProcessDetail{
		Process: &process.Process{Pid: 22637},
	}
	cmd, err := Pdetail.Cmdinfo()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cmd)
}
