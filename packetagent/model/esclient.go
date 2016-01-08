package model

import (
	"github.com/golang/glog"
	"gopkg.in/olivere/elastic.v2"
	"time"
)

type ESClient struct {
	*elastic.Client
}

func Getclient(server string) (*ESClient, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(server))
	if err != nil {
		return nil, err

	}
	esclient := &ESClient{Client: client}
	return esclient, nil
}

//the input info should be the json byte[]
func (client *ESClient) Push(info []byte, indexstr string, instancetype string) error {
	infostr := string(info)
	id := time.Now().String()
	returnmessage, err := client.Index().Index(indexstr).Type(instancetype).Id(id).BodyJson(infostr).Do()
	if err != nil {
		glog.Info(returnmessage)
		return err
	}
	glog.Info("upload ok")
	return nil
}
