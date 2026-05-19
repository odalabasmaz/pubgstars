package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
)

func Handler(ctx context.Context, event svc.RequestEvent) (svc.Response, error) {
	gameId := svc.CovertToString(event.Body["id"])
	switch event.Context["http-method"] {
	case "POST":
		return svc.Response{StatusCode: 200, Body: listGameUsers(gameId)}, nil
	default:
		return svc.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + event.Context["http-method"]}, nil
	}
}

func listGameUsers(gameId string) []model.User {
	var users []model.User
	for _, userId := range svc.GetGameUsersByGameId(gameId).Users {
		log.Println(userId)
		user := svc.GetUserById(userId)
		users = append(users, model.User{Username: user.Username})
	}
	return users
}

func main() {
	lambda.Start(Handler)
}
