package main

import (
	AwsUtils "../../internal"
	Model "../../model"
	Tables "../../model/tables"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"log"
)

func Handler(ctx context.Context, event AwsUtils.RequestEvent) (AwsUtils.Response, error) {
	log.Println("begin gamesHistory!!")

	httpMethod := event.Context["http-method"]
	switch httpMethod {
	case "GET":
		return AwsUtils.Response{StatusCode: 200, Body: listGamesHistory()}, nil
	default:
		return AwsUtils.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + httpMethod}, nil
	}
}

func listGamesHistory() []Model.Game {
	filt := expression.Name("status").Equal(expression.Value("completed"))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		fmt.Println(err)
	}
	games, err := AwsUtils.ListGames(Tables.GAMES, expr)
	if err != nil {
		fmt.Println(err)
		return games
	}
	return games
}

func main() {
	//listGamesHistory()
	lambda.Start(Handler)
}
