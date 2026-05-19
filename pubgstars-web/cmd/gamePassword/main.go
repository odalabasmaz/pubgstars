package main

import (
	AwsUtils "../../internal"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"time"
)

func Handler(ctx context.Context, event AwsUtils.RequestEvent) (AwsUtils.Response, error) {
	log.Println("begin show password !!")

	email := AwsUtils.GetUsernameFromJwtToken(event.Params["header"]["Authorization"])
	gameMap := event.Body
	gameId := AwsUtils.CovertToString(gameMap["id"])

	game := AwsUtils.GetGameById(gameId)
	location, _ := time.LoadLocation("Europe/Istanbul")
	now := time.Now().In(location)
	log.Print(now)
	gameDate, err := time.ParseInLocation("200601021504", game.GameDate, location)
	log.Print(gameDate)

	if err != nil {
		return AwsUtils.Response{StatusCode: 400, ErrorMessage: "Oyun tarihi geçersiz"}, err
	}

	if gameDate.After(now.Add(1 * time.Hour)) {
		return AwsUtils.Response{StatusCode: 400, ErrorMessage: "Oyun saatine 1 saatten fazla süre bulunmaktadır. Şifre alınamaz!"}, err
	}

	user := AwsUtils.GetUserByEmail(email)
	userId := user.Id
	userGame := AwsUtils.GetUserGamesByUserId(userId)
	found := false
	for _, gameId := range userGame.Games {
		if game.Id == gameId {
			found = true
			break
		}
	}
	if !found {
		return AwsUtils.Response{StatusCode: 400, ErrorMessage: "Şifreyi Görebilmek için oyuna kayıtlı olanız gerekmektedir"}, err

	}

	return AwsUtils.Response{StatusCode: 200, Body: game}, nil
}

func main() {
	lambda.Start(Handler)
}
