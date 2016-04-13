// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
inserting:Time series data is also called points are written to the database using batch inserts.
A series is a combination of a measurement (time/values) and a set of tags.


*/

package influxdbbackend

import (
	"errors"
	"fmt"
	"github.com/cpmonitor/packetagent/metrics"
	influxpackage "github.com/influxdata/influxdb/client/v2"
	"os"
	"sync"
	"time"
)

type InfluxdbStorage struct {
	Influxclient    influxpackage.Client
	MachineName     string
	Database        string
	RetentionPolicy string
	BufferDuration  time.Duration
	LastWrite       time.Time
	//points          []*influxdb.Point
	Lock         sync.Mutex
	ReadyToFlush func() bool
}

const (
	Infotypepacket string = "type_packet"
	Infotypemetric string = "type_metric"
)

func queryDB(clnt influxpackage.Client, cmd string, dbname string) (res []influxpackage.Result, err error) {
	q := influxpackage.Query{
		Command:  cmd,
		Database: dbname,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func Getinfluxclient(influxserver string, username string, password string, dbname string) (*InfluxdbStorage, error) {
	client, err := Newinfluxclient(influxserver, username, password, dbname)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func Newinfluxclient(influxserver string, username string, password string, dbname string) (*InfluxdbStorage, error) {
	c, err := influxpackage.NewHTTPClient(influxpackage.HTTPConfig{
		Addr:     influxserver,
		Username: username,
		Password: password,
	})

	if err != nil {
		return nil, err
	}

	_, err = queryDB(c, fmt.Sprintf("CREATE DATABASE %s", dbname), dbname)
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	influxc := &InfluxdbStorage{
		Influxclient: c,
		MachineName:  hostname,
		Database:     dbname,
		LastWrite:    time.Now(),
	}

	return influxc, nil
}

func (self *InfluxdbStorage) AddStats(infotype string, measurement string, content interface{}) error {
	if infotype == Infotypepacket {
		//transfer interface into HttpTransaction
		httpinstance, ok := content.(*metrics.HttpTransaction)
		if !ok {
			return errors.New("fail in transformation")
		}
		influxclient := self.Influxclient
		//create bp point write bp point

		bp, _ := influxpackage.NewBatchPoints(influxpackage.BatchPointsConfig{
			Database:  self.Database,
			Precision: "us",
		})
		fields := map[string]interface{}{
			"respondtime": httpinstance.Respondtime,
		}
		tags := map[string]string{
			"Srcip":    httpinstance.Srcip,
			"Srcport":  httpinstance.Srcport,
			"Destip":   httpinstance.Destip,
			"Destport": httpinstance.Destport,
			// problems in processing nesting problem
			// "Requestdetail": httpinstance.Packetdetail.Requestdetail,
			// "Responddetail": httpinstance.Packetdetail.Responddetail,
		}
		fmt.Println("**the measurement**", measurement)
		point, err := influxpackage.NewPoint(measurement, tags, fields, time.Now())
		if err != nil {
			return err
		}
		fmt.Println("the point name:", point.Name())
		bp.AddPoint(point)
		influxclient.Write(bp)
	}
	//common metric info
	if infotype == Infotypemetric {

	}

	return nil
}

func (self *InfluxdbStorage) Close() error {
	err := self.Influxclient.Close()
	if err != nil {
		return err
	}
	return nil
}
