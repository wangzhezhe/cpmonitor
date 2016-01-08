package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Test struct {
	Testa string `json:"testa"`
	Testb string `json:"testb"`
	Testc int    `json:"testc"`
}

func TestGetclient(t *testing.T) {
	server := "http://127.0.0.1:9200/"
	client, err := Getclient(server)
	if err != nil {
		fmt.Println("fail to create the client:", err)
		return
	}
	fmt.Println(client)

	// Create an index
	//_, err = client.CreateIndex("testagent").Do()
	if err != nil {
		fmt.Println(err)
		return
	}

	message := &Test{Testa: "Testing messages a", Testb: "Testing message b", Testc: 100}
	//client.Push(message, "Test")
	jsonmessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("json message:", string(jsonmessage))
	client.Push(jsonmessage, "Test")

}
