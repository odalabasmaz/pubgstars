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
	gameMap := event.Body
	gameId := AwsUtils.CovertToString(gameMap["id"])

	switch httpMethod {
	case "POST":
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
		var userTemp Model.User
		userTemp.Username = user.Username
		users = append(users, userTemp)
	}
	return users
}

func main() {
	//listGameUsers("6SyTSCalGW")
	lambda.Start(Handler)
}
