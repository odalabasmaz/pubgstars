package main

import (
	AwsUtils "../../internal"
	GameUtils "../../internal"
	Model "../../model"
	Tables "../../model/tables"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"log"
)

func Handler(ctx context.Context, event AwsUtils.RequestEvent) (AwsUtils.Response, error) {
	log.Println("begin !!")
	email := AwsUtils.GetUsernameFromJwtToken(event.Params["header"]["Authorization"])
	log.Println("username found: " + email)

	httpMethod := event.Context["http-method"]
	switch httpMethod {
	case "GET":
		return AwsUtils.Response{StatusCode: 200, Body: listUserGames(email)}, nil
	default:
		return AwsUtils.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + httpMethod}, nil
	}
}

func listUserGames(email string) []Model.Game {
	filt := expression.Name("status").Equal(expression.Value("active"))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		fmt.Println(err)
	}
	games, err := AwsUtils.ListGames(Tables.GAMES, expr)

	if err != nil {
		log.Println("Error occurred.")
		log.Println(err)
	} else if len(games) != 0 {
		user := AwsUtils.GetUserByEmail(email)
		userGame := AwsUtils.GetUserGamesByUserId(user.Id)

		//TODO: pass by reference??
		for i, game := range games {
			for _, gameId := range userGame.Games {
				if game.Id == gameId {
					games[i].Registered = true
					games[i].ShowPassword = GameUtils.IsGameInLastHour(game)
					break
				}
			}
		}
	}

	return games
}

func main() {
	//listUserGames("odalabasmaz+pg1@gmail.com")
	lambda.Start(Handler)
}
