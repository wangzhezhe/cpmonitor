package influxdbbackend

import (
	"fmt"
	"github.com/cpmonitor/packetagent/metrics"
	"testing"
	"time"
)

func TestProcessinfo(t *testing.T) {
	c, err := Getinfluxclient("http://127.0.0.1:8086", "wangzhe", "123456", "testb")
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
	infotype := "type_packet"
	measurement := "packetabcde"
	httpinstance := &metrics.HttpTransaction{
		Srcip:       "8082",
		Srcport:     "8080",
		Destip:      "8081",
		Respondtime: 0.123456,

		Packetdetail: metrics.Packetdetail{Requestdetail: "abbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", Responddetail: "b"},
	}
	err = c.AddStats(infotype, measurement, httpinstance)
	if err != nil {
		fmt.Println(err)
	}

}
