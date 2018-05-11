package main

import (
	"github.com/aws/aws-lambda-go/events"
	"bartenderAsFunctionServer/dao"
	"github.com/aws/aws-lambda-go/lambda"
	"fmt"
	"bartenderAsFunctionServer/model"
	"github.com/aws/aws-sdk-go/aws"
	serviceLambda "github.com/aws/aws-sdk-go/service/lambda"
	"os"
	"encoding/json"
)

var DataConnectionManager dao.RegistryConnectionInterface
var lambdaName = os.Getenv("LAMBDA_SERVE_USER")

func Handler(event events.CodePipelineEvent) error {
	fmt.Println("executing handler for serving")
	users, err := DataConnectionManager.GetRegistry()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("invoking for %d users", len(users))
	for _, user := range users {
		invokeLambdaForUser(user)
	}
	return nil

}
func invokeLambdaForUser(user model.User) {
	lambdaConnectionManager := dao.InitLambdaConnectionManager()
	bytes, _ := json.Marshal(user)
	lambdaConnectionManager.LambdaClient.Invoke(&serviceLambda.InvokeInput{FunctionName: aws.String(lambdaName), Payload: bytes})
}

func main() {
	DataConnectionManager = dao.CreateRegistryConnection()
	lambda.Start(Handler)
}
