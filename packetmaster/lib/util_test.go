package lib

import (
	"cpmonitor/packetmaster/lib"
	"fmt"
	"testing"
)

func TestCheckip(t *testing.T) {
	t.SkipNow()
	endpointa := "10.10.103.131:9080"
	endpointb := "10.10.103"
	resulta := lib.Checkip(endpointa)
	resultb := lib.Checkip(endpointb)
	fmt.Println(resulta, resultb)
}

func TestCheckserver(t *testing.T) {
	ip := "10.10.103.131"
	port := "8080"
	iflisten, pinfo, err := Checkserver(ip, port)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("if the server is listening:", iflisten, "the process info:", pinfo)

}
