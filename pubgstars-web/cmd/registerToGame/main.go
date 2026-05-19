package main

import (
	AwsUtils "../../internal"
	DataService "../../internal"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

func Handler(ctx context.Context, event AwsUtils.RequestEvent) (AwsUtils.Response, error) {
	log.Println("begin !!")

	/// TX Begin
	email := AwsUtils.GetUsernameFromJwtToken(event.Params["header"]["Authorization"])
	user := AwsUtils.GetUserByEmail(email)

	gameMap := event.Body
	gameId := AwsUtils.CovertToString(gameMap["id"])
	game := AwsUtils.GetGameById(gameId)

	e := DataService.RegisterUserToGame(user, game)
	if e != nil {
		return AwsUtils.Response{StatusCode: 400, ErrorMessage: e.Error()}, nil
	}
	/// TX End

	return AwsUtils.Response{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
