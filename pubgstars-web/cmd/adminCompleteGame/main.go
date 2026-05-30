package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
)

type App struct {
	store svc.Store
}

func (a *App) Handle(ctx context.Context, event svc.RequestEvent) (svc.Response, error) {
	operator := svc.GetUsernameFromJwtTokenForAdmin(event.Params["header"]["Authorization"])
	switch event.Context["http-method"] {
	case "POST":
		if err := a.completeGame(event.Body, operator); err != nil {
			log.Printf("completeGame error: %v", err)
			return svc.Response{StatusCode: 400, ErrorMessage: err.Error()}, nil
		}
		return svc.Response{StatusCode: 200}, nil
	default:
		return svc.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + event.Context["http-method"]}, nil
	}
}

func (a *App) completeGame(gameMap map[string]interface{}, operator string) error {
	gameId, ok := gameMap["gameId"].(string)
	if !ok {
		return fmt.Errorf("gameId is missing or not a string")
	}
	firstWinnerId, ok := gameMap["firstWinner"].(string)
	if !ok {
		return fmt.Errorf("firstWinner is missing or not a string")
	}
	secondWinnerId, ok := gameMap["secondWinner"].(string)
	if !ok {
		return fmt.Errorf("secondWinner is missing or not a string")
	}
	thirdWinnerId, ok := gameMap["thirdWinner"].(string)
	if !ok {
		return fmt.Errorf("thirdWinner is missing or not a string")
	}

	game := a.store.GetGameById(gameId)
	game.UpdatedAt = svc.CurrentTimeMillis()
	game.UpdatedBy = operator
	game.Status = "completed"
	game.Winner1st = firstWinnerId
	game.Winner2nd = secondWinnerId
	game.Winner3rd = thirdWinnerId

	firstWinner := a.store.GetUserById(game.Winner1st)
	firstWinner.Balance += game.Award1st
	firstWinner.Gain += game.Award1st

	secondWinner := a.store.GetUserById(game.Winner2nd)
	secondWinner.Balance += game.Award2nd
	secondWinner.Gain += game.Award2nd

	thirdWinner := a.store.GetUserById(game.Winner3rd)
	thirdWinner.Balance += game.Award3rd
	thirdWinner.Gain += game.Award3rd

	return a.store.CompleteGame(operator, game, firstWinner, secondWinner, thirdWinner)
}

func main() {
	app := &App{store: svc.NewDynamoStore()}
	lambda.Start(app.Handle)
}
