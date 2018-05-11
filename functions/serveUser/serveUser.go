package main

import (
	"bartenderAsFunctionServer/model"
	"github.com/aws/aws-lambda-go/lambda"
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
)

type statistics struct {
	User   string `json:"user"`
	Billed int    `json:"billed"`
}

func Handler(event model.User) error {
	fmt.Println("invoking for user ", event.Name)
	fmt.Println("invoking for url ", event.Url)
	response, err := http.Get(event.Url + "/command")
	if err != nil {
		fmt.Println(err)
	}
	var commands []model.Command
	data, _ := ioutil.ReadAll(response.Body)
	fmt.Println("data:", data)
	json.Unmarshal([]byte(data), &commands)

	for counter, command := range commands {
		if len(command.Beer) > 0 {
			for _, beerItem := range command.Beer {
				ii, _ := json.Marshal(beerItem)
				http.Post(event.Url+"/command/"+command.IdCommand+"/beer/serve", "application/json", bytes.NewBuffer(ii))
			}
		}
		if len(command.Food) > 0 {
			for _, foodItem := range command.Food {
				ii, _ := json.Marshal(foodItem)
				http.Post(event.Url+"/command/"+command.IdCommand+"/food/serve", "application/json", bytes.NewBuffer(ii))
			}
		}
		if counter > 10 {
			break
		}
	}
	//get bills and count
	responseBill, err := http.Get(event.Url + "/bill")
	if err != nil {
		fmt.Println(err)
	}
	items := model.CommandRequest{}

	dataBill, _ := ioutil.ReadAll(responseBill.Body)
	fmt.Println("dataBill:", data)
	json.Unmarshal([]byte(dataBill), &items)

	amt := 0
	for _, beer := range items.Beer {
		amt += beer.Amount
	}
	for _, food := range items.Food {
		amt += food.Amount * 2
	}
	stats := statistics{User: event.Name, Billed: amt}

	fmt.Println("stats: ", stats)
	return nil
}

func main() {
	lambda.Start(Handler)
}
