package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
)

func Handler(ctx context.Context, event svc.RequestEvent) (svc.Response, error) {
	operator := svc.GetUsernameFromJwtTokenForAdmin(event.Params["header"]["Authorization"])
	switch event.Context["http-method"] {
	case "GET":
		return listMessages(), nil
	case "POST":
		return updateMessage(event.Body, operator), nil
	default:
		return svc.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + event.Context["http-method"]}, nil
	}
}

func listMessages() svc.Response {
	expr, err := expression.NewBuilder().Build()
	if err != nil {
		fmt.Println(err)
	}
	messages, err := svc.ListMessages(expr)
	if err != nil {
		return svc.Response{StatusCode: 400, ErrorMessage: err.Error()}
	}
	return svc.Response{StatusCode: 200, Body: messages}
}

func updateMessage(messageMap map[string]interface{}, operator string) svc.Response {
	log.Println("updateMessages...")
	messageMap["updatedAt"] = svc.CurrentTimeMillis()
	messageMap["updatedBy"] = operator
	if err := svc.SaveMessage(messageMap); err != nil {
		return svc.Response{StatusCode: 400, ErrorMessage: err.Error()}
	}
	return svc.Response{StatusCode: 200, Body: "ok"}
}

func main() {
	lambda.Start(Handler)
}
