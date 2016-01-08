package lib

import (
	"io/ioutil"
	"log"
	"net/http"
	//"github.com/golang/glog"
	"encoding/json"
	"regexp"
)

var (
	agentpath   string = "/packet/checkport/"
	agentlisten string = "9990"
	client             = &http.Client{}
)

func Checkip(endpoint string) bool {
	checkok := false
	endpointregex := "^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}:[0-9]{2,5}$"
	match, _ := regexp.MatchString(endpointregex, endpoint)
	if match {
		checkok = true
	}
	return checkok
}

//check the ip:port is listenging
func Checkserver(ip string, port string) (bool, string, error) {
	agenturl := "http://" + ip + ":" + agentlisten + agentpath + port
	log.Println("the check url :", agenturl)
	request, err := http.NewRequest("GET", agenturl, nil)
	if err != nil {
		return false, "", err
	}
	response, err := client.Do(request)
	if err != nil {
		return false, "", err
	}
	responsebody, _ := ioutil.ReadAll(response.Body)
	//statmap["ifserver"] = xxx
	//statmap["psummary"] = xxx
	var processstat map[string]interface{}
	json.Unmarshal(responsebody, &processstat)
	log.Printf("return data: %+v:", processstat)
	log.Println(processstat["ifserver"])
	iflisten, ok := processstat["ifserver"].(string)
	if ok && iflisten == "true" {
		pinfo := processstat["psummary"].(string)
		return true, pinfo, nil
	} else {
		return false, "", nil
	}

}
