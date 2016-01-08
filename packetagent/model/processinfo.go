package model

import (
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

type ProcessDetail struct {
	*process.Process
}

func Memoryinfo() error {
	v, _ := mem.VirtualMemory()
	fmt.Println(v)
	return nil
}
func (pdetail *ProcessDetail) Cmdinfo() (string, error) {
	cmd, err := pdetail.Cmdline()
	if err != nil {
		return "", err
	}
	return cmd, nil
}

//TODO capsulate more functions if it is needed
