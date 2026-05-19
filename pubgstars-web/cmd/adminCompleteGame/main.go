package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
)

func Handler(ctx context.Context, event svc.RequestEvent) (svc.Response, error) {
	operator := svc.GetUsernameFromJwtTokenForAdmin(event.Params["header"]["Authorization"])
	switch event.Context["http-method"] {
	case "POST":
		return svc.Response{StatusCode: 200, Body: completeGame(event.Body, operator)}, nil
	default:
		return svc.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + event.Context["http-method"]}, nil
	}
}

func completeGame(gameMap map[string]interface{}, operator string) interface{} {
	gameId := gameMap["gameId"].(string)
	game := svc.GetGameById(gameId)
	game.UpdatedAt = svc.CurrentTimeMillis()
	game.UpdatedBy = operator
	game.Status = "completed"
	game.Winner1st = gameMap["firstWinner"].(string)
	game.Winner2nd = gameMap["secondWinner"].(string)
	game.Winner3rd = gameMap["thirdWinner"].(string)

	firstWinner := svc.GetUserById(game.Winner1st)
	firstWinner.Balance += game.Award1st
	firstWinner.Gain += game.Award1st

	secondWinner := svc.GetUserById(game.Winner2nd)
	secondWinner.Balance += game.Award2nd
	secondWinner.Gain += game.Award2nd

	thirdWinner := svc.GetUserById(game.Winner3rd)
	thirdWinner.Balance += game.Award3rd
	thirdWinner.Gain += game.Award3rd

	if err := svc.CompleteGame(operator, game, firstWinner, secondWinner, thirdWinner); err != nil {
		log.Printf("completeGame error: %v", err)
		return err.Error()
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
