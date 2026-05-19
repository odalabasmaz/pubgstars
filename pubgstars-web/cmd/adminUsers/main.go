package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
)

func Handler(ctx context.Context, event svc.RequestEvent) (svc.Response, error) {
	switch event.Context["http-method"] {
	case "GET":
		return svc.Response{StatusCode: 200, Body: listUsers()}, nil
	default:
		return svc.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + event.Context["http-method"]}, nil
	}
}

func listUsers() []model.User {
	users, err := svc.ListUsers()
	if err != nil {
		log.Println("listUsers error:", err)
	}
	return users
}

func main() {
	lambda.Start(Handler)
}
