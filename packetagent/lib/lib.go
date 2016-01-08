package lib

import (
	"bufio"
	"fmt"
	"github.com/golang/glog"
	"io"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

func Systemexec(s string) error {
	cmd := exec.Command("/bin/sh", "-c", s)
	if glog.V(1) {
		glog.Info(s)
	}
	out, err := cmd.StdoutPipe()
	go func() {
		o := bufio.NewReader(out)
		for {
			line, _, err := o.ReadLine()
			if err == io.EOF {
				break
			} else {
				fmt.Println(string(line))
				//if glog.V(1) {
				//	glog.Info(line)
				//}
			}
		}
	}()
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

//get the ipv4 addr of local machine
func CheckLocalip(iface string) (string, error) {
	ifaceobj, err := net.InterfaceByName(iface)
	if err != nil {
		return "", err
	}
	addrarry, err := ifaceobj.Addrs()
	if err != nil {
		return "", err
	}
	var localip = ""
	if glog.V(1) {
		glog.Info(addrarry)
	}
	for _, ip := range addrarry {
		IP := ip.String()
		if strings.Contains(IP, "/24") {
			localip = strings.TrimSuffix(IP, "/24")
		}
	}

	return localip, nil
}

//serch info from local machine
func Getinfofromport(portip int) (int, string, error) {
	//using gonetstat to get the pid of this app
	p, err := Tcpfromport(portip)
	// format header
	//fmt.Printf("Proto %16s %20s %14s %24s\n", "Local Adress", "Foregin Adress","State", "Pid/Program")

	if err != nil {
		return -1, "", err
	} else {
		//ip_port := fmt.Sprintf("%v:%v", p.Ip, p.Port)
		//fip_port := fmt.Sprintf("%v:%v", p.ForeignIp, p.ForeignPort)
		//pid_program := fmt.Sprintf("%v/%v", p.Pid, p.Name)
		//fmt.Printf("tcp %16v %20v %16v %20v\n", ip_port, fip_port,p.State, pid_program)
		pid, _ := strconv.Atoi(p.Pid)
		//fmt.Printf("detail %v \n", p)
		pinfo := "pname:" + p.Name + " pid:" + p.Pid + " exe:" + p.Exe + " pstat:" + p.State
		return pid, pinfo, nil
	}
}

//using sys operation to get the pid according to the command lsof
func Getinfofromportbylsof(portip int) error {
	cmd := "lsof -i :" + strconv.Itoa(portip) + " > null"
	err := Systemexec(cmd)
	if err != nil {
		return err
	}
	return nil
}
