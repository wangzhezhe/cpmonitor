package app

import (
	"fmt"
	"testing"
)

func TestProcessinfo(t *testing.T) {
	c, err := Newcpmanager()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c)
}
