package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model/tables"
)

func Handler(ctx context.Context, event svc.RequestEvent) (svc.Response, error) {
	switch event.Context["http-method"] {
	case "GET":
		return svc.Response{StatusCode: 200, Body: listGamesHistory()}, nil
	default:
		return svc.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + event.Context["http-method"]}, nil
	}
}

func listGamesHistory() []model.Game {
	filt := expression.Name("status").Equal(expression.Value("completed"))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		fmt.Println(err)
	}
	games, err := svc.ListGames(tables.GAMES, expr)
	if err != nil {
		log.Println("listGamesHistory error:", err)
	}
	return games
}

func main() {
	lambda.Start(Handler)
}
