package main

import (
	AwsUtils "../../internal"
	Model "../../model"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

var (
	db = AwsUtils.GetDynamoDbClient("eu-central-1")
)

func Handler(ctx context.Context, event AwsUtils.RequestEvent) (AwsUtils.Response, error) {
	log.Println("begin !!")
	//operator := AwsUtils.GetUsernameFromJwtTokenForAdmin(event.Params["header"]["Authorization"])
	httpMethod := event.Context["http-method"]
	switch httpMethod {
	case "GET":
		return AwsUtils.Response{StatusCode: 200, Body: listUsers()}, nil
	default:
		return AwsUtils.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + httpMethod}, nil
	}
}

func listUsers() []Model.User {
	users, err := AwsUtils.ListUsers()
	if err != nil {
		log.Println("Error occurred.")
		log.Println(err)
	}
	return users
}

func main() {
	lambda.Start(Handler)
}
