package main

import (
	"bartenderAsFunctionServer/model"
	"github.com/aws/aws-lambda-go/events"
	"encoding/json"
	"bartenderAsFunctionServer/dao"
	"github.com/aws/aws-lambda-go/lambda"
)

var DataConnectionManager dao.RegistryConnectionInterface

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := model.User{}
	json.Unmarshal([]byte(request.Body), &body)
	saveCommandError := DataConnectionManager.SaveRegistry(body)
	if saveCommandError != nil {
		return events.APIGatewayProxyResponse{StatusCode: 404, Body: "{\" error\":\"error\"}"}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: request.Body}, nil

}

func main() {
	DataConnectionManager = dao.CreateRegistryConnection()
	lambda.Start(Handler)
}
