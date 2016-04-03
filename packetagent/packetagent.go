package main

import (
	"flag"
	"github.com/cpmonitor/packetagent/app"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"net/http"
)

// This example shows the minimal code needed to get a restful.WebService working.
//
// GET http://localhost:8080/hello

var (
	Listenport = "9990"
	Backend    string
)

func main() {
	flag.Parse()

	flag.StringVar(&Backend, "influx", "http://localhost:8086", "server of the influxdb")
	app.Influxserver = Backend
	wsContainer := restful.NewContainer()
	if glog.V(1) {
		glog.Info("start listen packet")
	}
	app.Register(wsContainer)

	manager, err := app.Newcpmanager()
	if err != nil {
		glog.Error(err)
	}
	manager.Start()
	server := &http.Server{Addr: ":" + Listenport, Handler: wsContainer}
	server.ListenAndServe()
	manager.Stop()
}
