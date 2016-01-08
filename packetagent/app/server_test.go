package app

import (
	"flag"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	"testing"
)

func TestRegister(t *testing.T) {
	t.SkipNow()
	flag.Parse()
	wsContainer := restful.NewContainer()
	Register(wsContainer)
	log.Println("start listening")
	server := &http.Server{Addr: ":" + "9998", Handler: wsContainer}
	server.ListenAndServe()
	/*
		client := &http.Client{}
		url := "http://127.0.0.1:9998/packet/9998"
		log.Println("send the request")
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err.Error())
		}

		response, _ := client.Do(request)
		log.Println(response.Body)
	*/
}

func TestGetprocessinfo(t *testing.T) {
	portnum := 8080
	pinfo, err := Getprocessinfo(portnum)
	if err != nil {
		log.Println(err)
	}
	log.Println(pinfo)
	/*
		client := &http.Client{}
		url := "http://127.0.0.1:9998/packet/9998"
		log.Println("send the request")
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err.Error())
		}

		response, _ := client.Do(request)
		log.Println(response.Body)
	*/
}
