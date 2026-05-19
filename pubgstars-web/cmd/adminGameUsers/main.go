package main

import (
	AwsUtils "../../internal"
	Model "../../model"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

func Handler(ctx context.Context, event AwsUtils.RequestEvent) (AwsUtils.Response, error) {
	log.Println("begin !!")
	//operator := AwsUtils.GetUsernameFromJwtTokenForAdmin(event.Params["header"]["Authorization"])
	httpMethod := event.Context["http-method"]
	gameId := event.Params["querystring"]["gameId"]

	switch httpMethod {
	case "GET":
		return AwsUtils.Response{StatusCode: 200, Body: listGameUsers(gameId)}, nil
	default:
		return AwsUtils.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + httpMethod}, nil
	}
}

func listGameUsers(gameId string) []Model.User {
	var users []Model.User
	userIds := AwsUtils.GetGameUsersByGameId(gameId).Users
	for _, userId := range userIds {
		log.Println(userId)
		user := AwsUtils.GetUserById(userId)
		users = append(users, user)
	}
	return users
}

func main() {
	//listGameUsers("1FPvPfFs1LW")
	lambda.Start(Handler)
}
