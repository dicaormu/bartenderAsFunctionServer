package dao

import (
	"github.com/aws/aws-sdk-go/aws"
	"os"
	"strings"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"bartenderAsFunctionServer/model"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type RegistryConnection struct {
	DynamoConnection *dynamodb.DynamoDB
}

type RegistryConnectionInterface interface {
	SaveRegistry(command model.User) error
	GetRegistry() ([]model.User, error)
	GetRegistryById(id string) model.User
}

func  (con *RegistryConnection) GetRegistryById(name string) model.User {
	tableName := os.Getenv("TABLE_REGISTRY")
	result, err := con.DynamoConnection.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(name),
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		return model.User{}
	}
	command := model.User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &command)

	if err != nil {
		fmt.Println(err)
		return model.User{}
	}

	if command.Name == "" {
		return model.User{}
	}
	return command
}


func (con *RegistryConnection) SaveRegistry(user model.User) error {
	tableName := os.Getenv("TABLE_REGISTRY")
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Create item in table
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err = con.DynamoConnection.PutItem(input)
	if err != nil {
		fmt.Println("error saving", err)
		return err
	}
	return nil
}

func (con *RegistryConnection) GetRegistry() ([]model.User, error) {
	tableName := os.Getenv("TABLE_REGISTRY")
	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	result, err := con.DynamoConnection.Scan(params)
	if err != nil {
		fmt.Println("error getting dynamo", err)
		return nil, err
	}
	var users []model.User
	for _, v := range result.Items {
		var item model.User
		err = dynamodbattribute.UnmarshalMap(v, &item)
		users = append(users, item)
		fmt.Printf("getting %d users", len(users))
	}
	return users, nil
}

func CreateRegistryConnection() RegistryConnectionInterface {
	return &RegistryConnection{initializeDynamoDBClient()}
}

func initializeDynamoDBClient() *dynamodb.DynamoDB {
	localEnv := os.Getenv("AWS_SAM_LOCAL")
	dynamoUrl := os.Getenv("dynamoUrl")
	awsConfig := &aws.Config{Region: aws.String(os.Getenv("AWS_DEFAULT_REGION"))}
	if len(localEnv) > 0 && strings.ToLower(localEnv) == "true" {
		if dynamoUrl == "" {
			dynamoUrl = "http://docker.for.mac.localhost:8000"
		}
		awsConfig.Endpoint = aws.String(dynamoUrl)
	}
	sessionVar, err := session.NewSession(awsConfig)
	if err != nil {
		fmt.Println("error connecting", err)
		os.Exit(1)
	}
	return dynamodb.New(sessionVar, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody))
}
