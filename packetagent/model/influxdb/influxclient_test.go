package model

import (
	"fmt"
	"testing"
	"time"
)

func TestProcessinfo(t *testing.T) {
	c, err := Newinfluxclient("http://localhost:8086", "wangzhe", "123456", "test")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c)
	duration := time.Duration(time.Second)
	Duration, info, err := c.Influxclient.Ping(duration)
	c.Influxclient.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
	fmt.Println(Duration)
}
