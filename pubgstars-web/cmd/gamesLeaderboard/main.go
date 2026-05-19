package main

import (
	AwsUtils "../../internal"
	Model "../../model"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"sort"
)

func Handler(ctx context.Context, event AwsUtils.RequestEvent) (AwsUtils.Response, error) {
	log.Println("begin games leaderboard!!")

	httpMethod := event.Context["http-method"]
	switch httpMethod {
	case "GET":
		return AwsUtils.Response{StatusCode: 200, Body: listGamesLeaderboard()}, nil
	default:
		return AwsUtils.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + httpMethod}, nil
	}
}

func listGamesLeaderboard() []Model.User {
	users, err := AwsUtils.ListUsers()
	if err != nil {
		return users
	}

	// sort
	sort.Slice(users, func(i, j int) bool {
		return users[i].Gain > users[j].Gain
	})

	// take top 3
	if len(users) > 3 {
		return users[:3]
	} else {
		return users
	}
}

func main() {
	//listGamesLeaderboard()
	lambda.Start(Handler)
}
