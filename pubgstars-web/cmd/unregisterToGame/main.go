package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
)

func Handler(ctx context.Context, event svc.RequestEvent) (svc.Response, error) {
	email := svc.GetUsernameFromJwtToken(event.Params["header"]["Authorization"])
	user := svc.GetUserByEmail(email)

	gameId := svc.CovertToString(event.Body["id"])
	game := svc.GetGameById(gameId)

	if err := svc.UnregisterUserToGame(user, game); err != nil {
		return svc.Response{StatusCode: 400, ErrorMessage: err.Error()}, nil
	}
	return svc.Response{StatusCode: 200}, nil
}

func main() {
	log.Println("unregisterToGame starting")
	lambda.Start(Handler)
}
