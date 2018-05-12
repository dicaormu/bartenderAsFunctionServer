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
	fmt.Println("data:", string(data))
	json.Unmarshal(data, &commands)

	for counter, command := range commands {
		if len(command.Beer) > 0 {
			serveItem(command.Beer, event.Url, command.IdCommand, "beer")
		}
		if len(command.Food) > 0 {
			serveItem(command.Food, event.Url, command.IdCommand, "food")
		}
		if counter > 10 {
			break
		}
	}
	//get bills and count
	responseBill, err := http.Get(event.Url + "/bill")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	items := model.CommandRequest{}

	dataBill, _ := ioutil.ReadAll(responseBill.Body)
	fmt.Println("dataBill:", string(dataBill))
	json.Unmarshal(dataBill, &items)

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

func serveItem(items []model.Item, url string, idcommand string, typeItem string) {
	for _, item := range items {
		ii, _ := json.Marshal(item)
		_, err := http.Post(url+"/command/"+idcommand+"/"+typeItem+"/serve", "application/json", bytes.NewBuffer(ii))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func main() {
	lambda.Start(Handler)
}
