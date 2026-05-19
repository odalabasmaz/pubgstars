package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
)

func Handler(ctx context.Context, event svc.RequestEvent) (svc.Response, error) {
	gameId := event.Params["querystring"]["gameId"]
	switch event.Context["http-method"] {
	case "GET":
		return svc.Response{StatusCode: 200, Body: listGameUsers(gameId)}, nil
	default:
		return svc.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + event.Context["http-method"]}, nil
	}
}

func listGameUsers(gameId string) []model.User {
	var users []model.User
	for _, userId := range svc.GetGameUsersByGameId(gameId).Users {
		log.Println(userId)
		users = append(users, svc.GetUserById(userId))
	}
	return users
}

func main() {
	lambda.Start(Handler)
}
