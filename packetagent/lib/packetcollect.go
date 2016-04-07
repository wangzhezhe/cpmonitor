package lib

import (
	//"encoding/json"
	"github.com/golang/glog"

	"net"

	"strings"

	"container/list"
	"github.com/cpmonitor/packetagent/metrics"
	"github.com/cpmonitor/packetagent/model"
	"github.com/cpmonitor/packetagent/model/influxdbbackend"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"sync"
	"time"
)

var (
	snapshotLen      int32  = 65535
	promiscuous      bool   = false
	ESSERVER         string = "http://10.10.105.86:9200/"
	err              error
	timeout          time.Duration = 2 * time.Second
	handle           *pcap.Handle
	localip          string
	httpinstancelist *list.List
	ESClient         *model.ESClient
	Activeflag       bool = false //if flag is true , the agent is collecting the data
	Flagmutex             = &sync.Mutex{}

	influxserver = "http://10.10.105.33:8086"
	username     = "wangzhe"
	password     = "123456"
	dbname       = "test"
	Influxclient *influxdbbackend.InfluxdbStorage
)

func init() {
	glog.Info("init  the ESClient")
	//ESClient, err = model.Getclient(ESSERVER)
	Influxclient, err = influxdbbackend.Getinfluxclient(influxserver, username, password, dbname)
	if err != nil {
		glog.Info("fail to create the client !!!:", err)
		return
	}
}
func checkLocalip(iface string) (string, error) {
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

//detect the http packet return the info
func detectHttp(packet gopacket.Packet) (bool, []byte) {
	applicationLayer := packet.ApplicationLayer()
	if applicationLayer != nil {
		if strings.Contains(string(applicationLayer.Payload()), "HTTP") {

			if glog.V(1) {
				glog.Info("HTTP found!")
			}
			return true, applicationLayer.LayerContents()
		} else {
			return false, nil
		}
	} else {
		return false, nil
	}
}

//if it is the output stream from local machine
func outputStream(packet gopacket.Packet, Srcaddr *metrics.Address, Destaddr *metrics.Address) {
	ishttp, httpcontent := detectHttp(packet)
	if httpcontent != nil {
		if glog.V(1) {
			//glog.Info("the content of packet sent:", string(httpcontent))
		}
	}

	if ishttp {
		sendtime := time.Now()
		//iphandler := packet.Layer(layers.LayerTypeIPv4)
		reqdetail := string(packet.ApplicationLayer().LayerContents())
		httpinstance := &metrics.HttpTransaction{
			Srcip:        Srcaddr.IP,
			Srcport:      Srcaddr.PORT,
			Destip:       Destaddr.IP,
			Destport:     Destaddr.PORT,
			Timesend:     sendtime,
			Packetdetail: metrics.Packetdetail{Requestdetail: reqdetail, Responddetail: ""},
		}
		//put the httpinstance into a list
		if glog.V(1) {
			glog.Infof("store the instance:%v\n", httpinstance)
		}
		httpinstancelist.PushBack(httpinstance)
		if glog.V(2) {
			glog.Infof("the length of the list :", httpinstancelist.Len())
		}
	}

}

//adjust if this is the response of the packet
func ifreverse(httpinstance *metrics.HttpTransaction, Srcaddr *metrics.Address, Destaddr *metrics.Address) bool {
	if httpinstance.Srcip == Destaddr.IP && httpinstance.Destip == Srcaddr.IP {
		if httpinstance.Srcport == Destaddr.PORT && httpinstance.Destport == Srcaddr.PORT {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

//if it is the input stream to local machine
func inputStream(packet gopacket.Packet, Srcaddr *metrics.Address, Destaddr *metrics.Address) {
	//get the instance from the list which has the reverse srcaddr and the destaddr
	respdetail := string(packet.Data())
	if glog.V(1) {
		glog.Info("the length of the list before extract element:", httpinstancelist.Len())
	}

	for element := httpinstancelist.Front(); element != nil; element = element.Next() {
		httpinstance := element.Value.(*metrics.HttpTransaction)
		isreverse := ifreverse(httpinstance, Srcaddr, Destaddr)
		if isreverse {
			httpinstance.Timereceive = time.Now()
			//the units of time is ms
			responsetime := httpinstance.Timereceive.Sub(httpinstance.Timesend).Seconds() * 1000
			httpinstance.Respondtime = responsetime
			//store the respond detail
			glog.Info("respond info:", respdetail)
			httpinstance.Packetdetail.Responddetail = respdetail
			if glog.V(0) {
				glog.Infof("Respond duration:%vms", responsetime)
				glog.Infof("Get the response: %v", httpinstance)
				//glog.Info("detail:", httpinstance.Packetdetail)

			}
			httpinstancelist.Remove(element)
			//??how to use generic to realize the push function of different type
			//jsoninfo, _ := json.Marshal(httpinstance)
			//the type should be the ip of this machine
			//the first parameter is index the second one is type
			//err := ESClient.Push(jsoninfo, "packetagent", localip)
			measurement := "respondtime"
			err := Influxclient.AddStats("type_packet", measurement, httpinstance)
			if err != nil {
				glog.Info("error to push")
			}
			break
		}
	}

}

//every time get a new packet
func processPacketInfo(packet gopacket.Packet) {
	//get the specified layer
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		if glog.V(2) {
			glog.Info("TCP layer is detected.")
		}

		tcphandler, _ := tcpLayer.(*layers.TCP)
		srcport := tcphandler.SrcPort
		destport := tcphandler.DstPort
		//get the specified layer
		iplayer := packet.Layer(layers.LayerTypeIPv4)
		httphandler, _ := iplayer.(*layers.IPv4)
		srcip := httphandler.SrcIP
		destip := httphandler.DstIP
		//log.Println(srcip.String())
		//send the packet from local machine
		Srcaddr := &metrics.Address{IP: srcip.String(), PORT: srcport.String()}
		Destaddr := &metrics.Address{IP: destip.String(), PORT: destport.String()}
		if glog.V(2) {
			glog.Infof("srcaddr %v destaddr %v \n", Srcaddr, Destaddr)
		}
		var mutex = &sync.Mutex{}

		if srcip.String() == localip {
			mutex.Lock()
			outputStream(packet, Srcaddr, Destaddr)
			mutex.Unlock()
		}
		//get the packet from the local machine
		if destip.String() == localip {

			mutex.Lock()
			inputStream(packet, Srcaddr, Destaddr)
			mutex.Unlock()
		}

	}
}

func Startcollect(port int, device string, timesignal <-chan time.Time) {
	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		glog.Info(err.Error())
	}

	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	templocalip, err := checkLocalip(device)
	localip = templocalip
	if glog.V(0) {
		glog.Info(localip)
	}
	httpinstancelist = list.New()
	if err != nil {
		glog.Info(err.Error())
	}
A:
	for packet := range packetSource.Packets() {
		select {
		case <-timesignal:
			//stop the falg
			Flagmutex.Lock()
			Activeflag = false
			Flagmutex.Unlock()
			break A
		default:
			processPacketInfo(packet)

		}

	}
}
