package main

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
)

func Handler(ctx context.Context, event svc.RequestEvent) (svc.Response, error) {
	email := svc.GetUsernameFromJwtToken(event.Params["header"]["Authorization"])
	gameId := svc.CovertToString(event.Body["id"])
	game := svc.GetGameById(gameId)

	location, _ := time.LoadLocation("Europe/Istanbul")
	now := time.Now().In(location)
	gameDate, err := time.ParseInLocation("200601021504", game.GameDate, location)
	if err != nil {
		return svc.Response{StatusCode: 400, ErrorMessage: "Oyun tarihi geçersiz"}, nil
	}
	log.Printf("gameDate: %v, now: %v", gameDate, now)

	if gameDate.After(now.Add(1 * time.Hour)) {
		return svc.Response{StatusCode: 400, ErrorMessage: "Oyun saatine 1 saatten fazla süre bulunmaktadır. Şifre alınamaz!"}, nil
	}

	user := svc.GetUserByEmail(email)
	userGame := svc.GetUserGamesByUserId(user.Id)
	for _, gid := range userGame.Games {
		if game.Id == gid {
			return svc.Response{StatusCode: 200, Body: game}, nil
		}
	}
	return svc.Response{StatusCode: 400, ErrorMessage: "Şifreyi görebilmek için oyuna kayıtlı olmanız gerekmektedir"}, nil
}

func main() {
	lambda.Start(Handler)
}
