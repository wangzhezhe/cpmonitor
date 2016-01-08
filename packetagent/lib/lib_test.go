package lib

import (
	"fmt"
	"testing"
)

func TestGetinfofromport(t *testing.T) {
	portint := 5000
	pid, pname, err := Getinfofromport(portint)
	fmt.Println(pid, pname, err)
	portint = 9999
	pid, pname, err = Getinfofromport(portint)
	//pid should be the -1 if the portint is not be listened
	fmt.Println(pid, pname, err)
}

/*
func TestGetinfofromportbylsof(t *testing.T) {
	portint := 35472
	err := Getinfofromportbylsof(portint)
	//fmt.Println(pid, pname, err)
	if err != nil {
		fmt.Println(err)
	}

}
*/
