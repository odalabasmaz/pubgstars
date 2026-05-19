package main

import (
	AwsUtils "../../internal"
	DataService "../../internal"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"log"
)

func Handler(ctx context.Context, event AwsUtils.RequestEvent) (AwsUtils.Response, error) {
	log.Println("begin !!")
	operator := AwsUtils.GetUsernameFromJwtTokenForAdmin(event.Params["header"]["Authorization"])
	httpMethod := event.Context["http-method"]
	log.Println(event.Body)
	switch httpMethod {
	case "GET":
		return listMessages(), nil
	case "POST":
		return updateMessage(event.Body, operator), nil
	default:
		return AwsUtils.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + httpMethod}, nil
	}
}

func listMessages() AwsUtils.Response {
	expr, err := expression.NewBuilder().Build()
	if err != nil {
		fmt.Println(err)
	}
	if messages, err := AwsUtils.ListMessages(expr); err != nil {
		return AwsUtils.Response{StatusCode: 400, ErrorMessage: err.Error()}
	} else {
		return AwsUtils.Response{StatusCode: 200, Body: messages}
	}
}

func updateMessage(messageMap map[string]interface{}, operator string) AwsUtils.Response {
	log.Println("updateMessages...")
	currentTimeMillis := AwsUtils.CurrentTimeMillis()
	messageMap["updatedAt"] = currentTimeMillis
	messageMap["updatedBy"] = operator
	if err := DataService.SaveMessage(messageMap); err != nil {
		return AwsUtils.Response{StatusCode: 400, ErrorMessage: err.Error()}
	} else {
		return AwsUtils.Response{StatusCode: 200, Body: "ok"}
	}
}

func main() {
	lambda.Start(Handler)
}
