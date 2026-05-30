package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
)

type App struct {
	store svc.Store
}

func (a *App) Handle(ctx context.Context, event svc.RequestEvent) (svc.Response, error) {
	email := svc.GetUsernameFromJwtToken(event.Params["header"]["Authorization"])
	user := a.store.GetUserByEmail(email)

	gameId := svc.CovertToString(event.Body["id"])
	game := a.store.GetGameById(gameId)

	if err := a.store.UnregisterUserToGame(user, game); err != nil {
		return svc.Response{StatusCode: 400, ErrorMessage: err.Error()}, nil
	}
	return svc.Response{StatusCode: 200}, nil
}

func main() {
	log.Println("unregisterToGame starting")
	app := &App{store: svc.NewDynamoStore()}
	lambda.Start(app.Handle)
}
