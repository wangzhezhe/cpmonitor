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
	setMaxProcs()
	flag.Parse()
	flag.StringVar(&Backend, "influx", "http://127.0.0.1:8086", "server of the influxdb")
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

func setMaxProcs() {
	// TODO(vmarmol): Consider limiting if we have a CPU mask in effect.
	// Allow as many threads as we have cores unless the user specified a value.
	var numProcs int
	if *maxProcs < 1 {
		numProcs = runtime.NumCPU()
	} else {
		numProcs = *maxProcs
	}
	runtime.GOMAXPROCS(numProcs)

	// Check if the setting was successful.
	actualNumProcs := runtime.GOMAXPROCS(0)
	if actualNumProcs != numProcs {
		glog.Warningf("Specified max procs of %v but using %v", numProcs, actualNumProcs)
	}
}
