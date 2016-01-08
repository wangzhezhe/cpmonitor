package controllers

import (
	"cpmonitor/packetmaster/lib"
	"github.com/astaxie/beego"
	"log"
	"strings"
	"sync/atomic"
)

type MasterController struct {
	beego.Controller
}

var (
	index    string = "packetagent"
	esserver string = "http://127.0.0.1:9200"
)

type Endpoint struct {
	Ip             string
	Port           string
	Iflisten       bool
	Processsummary string
}

// @Title get associated app
// @router /assapp/:srcendpoint [get]
func (m *MasterController) Get() {
	srcendpoint := m.Ctx.Input.Param(":srcendpoint")
	ok := lib.Checkip(srcendpoint)
	if !ok {
		lib.Handleerr("unsuported endpoint format, should be the ip:port", &m.Ctx.ResponseWriter, nil)
		return

	}
	str_split := strings.Split(srcendpoint, ":")
	ipaddr := str_split[0]
	port := str_split[1]
	log.Println("addr:", ipaddr, "port:", port)

	//send the command to the agent let it start listening

	esclient, err := lib.Getclient(esserver)
	if err != nil {
		lib.Handleerr("fail to create the client", &m.Ctx.ResponseWriter, err)
		return
	}
	//returnlist, err := esclient.Aggregationterm_indirect(index, ipaddr, facetterm, port)
	queryname := "Srcport"
	queryvalue := port
	facetterm := "Destip"
	returnlist, err := esclient.Aggregationterm_direct(index, ipaddr, facetterm, queryname, queryvalue)
	if err != nil {
		lib.Handleerr("fail aggregation", &m.Ctx.ResponseWriter, err)
		return
	}

	//glog.Info(returnlist.Facets["tags"])
	tagmap := returnlist.Facets["tags"].Terms
	if len(tagmap) == 0 {
		m.Data["json"] = []Endpoint{}
		m.ServeJson()
		return
	}
	Associatetsock := []Endpoint{}
	for _, eachitem := range tagmap {
		ip := eachitem.Term.(string)

		facetterm = "Destport"
		queryname = "Destip"
		queryvalue = ip
		log.Println(ip)
		returnlist, err = esclient.Aggregationterm_direct(index, ipaddr, facetterm, queryname, queryvalue)
		log.Println(returnlist)
		portlist := returnlist.Facets["tags"].Terms
		if len(portlist) == 0 {
			m.Data["json"] = []Endpoint{}
			m.ServeJson()
			return
		}

		for _, portitem := range portlist {
			port := portitem.Term.(string)
			//this operation mybe is time-consuming
			//iflisten, _ := lib.Checkserver(ip, port)
			//Associatetsock = append(Associatetsock, Endpoint{Ip: ip, Port: port, Iflisten: iflisten})
			Associatetsock = append(Associatetsock, Endpoint{Ip: ip, Port: port})
		}

		log.Printf("returnlist:%v", Associatetsock)

	}
	var Returncount int32 = int32(len(Associatetsock))
	Indexcount := Returncount
	Signalreturn := make(chan string)
	var i int32
	for i = 0; i < Indexcount; i++ {
		go func(i int32) {
			log.Println("index:", i)
			iflisten, pinfo, _ := lib.Checkserver(Associatetsock[i].Ip, Associatetsock[i].Port)
			log.Println(iflisten)
			Associatetsock[i].Iflisten = iflisten
			Associatetsock[i].Processsummary = pinfo
			atomic.AddInt32(&Returncount, -1)
			log.Println("return count:", Returncount)
			log.Println(Associatetsock[i])
			if Returncount == 0 {
				Signalreturn <- "ok"
			}
		}(i)
	}
A:
	for {
		select {
		case <-Signalreturn:
			break A
		}
	}
	// TODO:send checkrequest to the agent to check if the port is Listened
	// time consuming to check all the endpoint now!!!!
	//lib.Checkserver(Associatetsock)

	m.Data["json"] = Associatetsock
	m.ServeJson()

}
