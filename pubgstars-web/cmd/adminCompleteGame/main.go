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
	operator := AwsUtils.GetUsernameFromJwtTokenForAdmin(event.Params["header"]["Authorization"])
	httpMethod := event.Context["http-method"]
	log.Println(event.Body)
	switch httpMethod {
	case "POST":
		return AwsUtils.Response{StatusCode: 200, Body: completeGame(event.Body, operator)}, nil
	default:
		return AwsUtils.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + httpMethod}, nil
	}
}

func completeGame(gameMap map[string]interface{}, operator string) interface{} {
	gameId := gameMap["gameId"].(string)
	first := gameMap["firstWinner"].(string)
	second := gameMap["secondWinner"].(string)
	third := gameMap["thirdWinner"].(string)
	game := AwsUtils.GetGameById(gameId)
	game.UpdatedAt = AwsUtils.CurrentTimeMillis()
	game.UpdatedBy = operator
	game.Status = "completed"
	game.Winner1st = first
	game.Winner2nd = second
	game.Winner3rd = third

	firstWinner := AwsUtils.GetUserById(first)
	firstWinner.Balance += game.Award1st
	firstWinner.Gain += game.Award1st

	secondWinner := AwsUtils.GetUserById(second)
	secondWinner.Balance += game.Award2nd
	secondWinner.Gain += game.Award2nd

	thirdWinner := AwsUtils.GetUserById(third)
	thirdWinner.Balance += game.Award3rd
	thirdWinner.Gain += game.Award3rd

	err := DataService.CompleteGame(operator, game, firstWinner, secondWinner, thirdWinner)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
