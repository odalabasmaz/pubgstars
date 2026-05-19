package main

import (
	"context"
	"log"
	"sort"

	"github.com/aws/aws-lambda-go/lambda"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
)

func Handler(ctx context.Context, event svc.RequestEvent) (svc.Response, error) {
	switch event.Context["http-method"] {
	case "GET":
		return svc.Response{StatusCode: 200, Body: listGamesLeaderboard()}, nil
	default:
		return svc.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + event.Context["http-method"]}, nil
	}
}

func listGamesLeaderboard() []model.User {
	users, err := svc.ListUsers()
	if err != nil {
		log.Println("listGamesLeaderboard error:", err)
		return nil
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].Gain > users[j].Gain
	})
	if len(users) > 3 {
		return users[:3]
	}
	return users
}

func main() {
	lambda.Start(Handler)
}
