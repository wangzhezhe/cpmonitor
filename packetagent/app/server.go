package app

import (
	"encoding/json"
	"errors"
	"github.com/cpmonitor/packetagent/lib"
	"github.com/cpmonitor/packetagent/model"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/shirou/gopsutil/process"
	"strconv"
	"time"
)

// This example shows the minimal code needed to get a restful.WebService working.
// GET http://localhost:8080/hello

var (
	Device       string        = "eth0"
	Defaulttime  time.Duration = 30
	Influxserver string
)

func Register(container *restful.Container) {
	if glog.V(1) {
		glog.Info("start regist")
	}
	ws := new(restful.WebService)
	ws.Path("/packet").
		Doc("control packet listening").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.GET("/startlisten/{port-id}").
		To(StartListen).
		Doc("start collect the packet from the specified port").
		Operation("findUser").
		Param(ws.PathParameter("port-id", "identifier of the listening port").DataType("string")))

	ws.Route(ws.GET("/checkport/{port-id}").
		To(Checkserver).
		Doc("check if the server is listening the port").
		Operation("checkport").
		Param(ws.PathParameter("port-id", "identifier of the listening port").DataType("string")))

	container.Add(ws)
}

//Check if the port is listening on this machine
func Checkserver(request *restful.Request, response *restful.Response) {
	portstring := request.PathParameter("port-id")
	glog.Info("get the port number", portstring)
	portint, err := strconv.Atoi(portstring)
	if err != nil {
		response.WriteError(500, err)
		return
	}
	statmap := make(map[string]string)
	statmap["ifserver"] = "false"
	statmap["psummary"] = "null"
	err = lib.Getinfofromportbylsof(portint)
	if err != nil {
		glog.Info(err)
		jstr, _ := json.Marshal(statmap)
		response.Write(jstr)
		return
	} else {
		//get the process info
		statmap["ifserver"] = "true"
		processinfo, _ := Getprocessinfo(portint)
		statmap["psummary"] = processinfo
		jstr, _ := json.Marshal(statmap)
		response.Write(jstr)
		return
	}
}

//when get the pid of that port
//start collecting the packet for 60s
//the topology only need to be check temporaray
func StartListen(request *restful.Request, response *restful.Response) {
	portstring := request.PathParameter("port-id")
	glog.Info("get the port number", portstring)
	portint, err := strconv.Atoi(portstring)
	if err != nil {
		response.WriteError(500, err)
		return
	}
	pid, _, err := lib.Getinfofromport(portint)

	if pid == -1 {
		response.WriteError(500, errors.New("the port is not be listend in this machine ( /proc/net/tcp and /proc/net/tcp6)"))
		return
	}
	if err != nil {
		response.WriteError(500, err)
		return
	}
	//start listen to specific ip:port for 60s and send the data to es
	timesignal := time.After(time.Second * Defaulttime)
	//start collect and check the timesignal every one minutes
	if !lib.Activeflag {
		go lib.Startcollect(portint, Device, timesignal)
		lib.Flagmutex.Lock()
		lib.Activeflag = true
		response.Write([]byte("activated"))
		lib.Flagmutex.Unlock()
	} else {
		response.Write([]byte("the server is already been activatied"))
	}
}

//if the port is listend by the server , collecting it's process info
func Getprocessinfo(portnum int) (string, error) {

	pid, pinfo, err := lib.Getinfofromport(portnum)
	//create the process instance and get the detail info of specified pid
	Pdetail := &model.ProcessDetail{
		Process: &process.Process{Pid: int32(pid)},
	}
	cmd, err := Pdetail.Cmdinfo()
	if err != nil {
		return "", err
	}
	glog.Info(cmd)
	pinfo = pinfo + " cmd:" + cmd
	return pinfo, nil
}
