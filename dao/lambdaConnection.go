package dao

import (
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/aws"
	"os"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type LambdaConnectionManager struct {
	LambdaClient lambdaiface.LambdaAPI
}

func initializeLambdaSessionClient() *lambda.Lambda {
	awsConfig := &aws.Config{Region: aws.String(os.Getenv("AWS_DEFAULT_REGION"))}
	sess, _ := session.NewSession(awsConfig)
	return lambda.New(sess)
}

func InitLambdaConnectionManager() *LambdaConnectionManager {
	client := new(LambdaConnectionManager)
	client.LambdaClient = initializeLambdaSessionClient()
	return client
}


