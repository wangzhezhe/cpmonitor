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

package model

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"os"
	"sync"
	"time"
)

type influxdbStorage struct {
	Influxclient    client.Client
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

func queryDB(clnt client.Client, cmd string, dbname string) (res []client.Result, err error) {
	q := client.Query{
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

func Newinfluxclient(influxserver string, username string, password string, dbname string) (*influxdbStorage, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
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

	influxc := &influxdbStorage{
		Influxclient: c,
		MachineName:  hostname,
		Database:     dbname,
		LastWrite:    time.Now(),
	}

	return influxc, nil
}

func (self *influxdbStorage) AddStats(infotype string) error {
	if infotype == Infotypepacket {

	}

	if infotype == Infotypemetric {

	}

	return nil
}

func (self *influxdbStorage) Close() error {
	err := self.Influxclient.Close()
	if err != nil {
		return err
	}
	return nil
}
