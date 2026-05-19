package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model/tables"
)

var db = svc.GetDynamoDbClient("eu-central-1")

func Handler(ctx context.Context, event svc.RequestEvent) (svc.Response, error) {
	operator := svc.GetUsernameFromJwtTokenForAdmin(event.Params["header"]["Authorization"])
	switch event.Context["http-method"] {
	case "GET":
		status := event.Params["querystring"]["status"]
		if status != "" {
			return svc.Response{StatusCode: 200, Body: listGamesWithStatus(status)}, nil
		}
		return svc.Response{StatusCode: 200, Body: listGames()}, nil
	case "PUT", "POST":
		if event.Body["id"] == nil {
			return svc.Response{StatusCode: 200, Body: addGame(event.Body, operator)}, nil
		}
		return svc.Response{StatusCode: 200, Body: updateGame(event.Body, operator)}, nil
	case "DELETE":
		return deleteGame(event.Body, operator)
	default:
		return svc.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + event.Context["http-method"]}, nil
	}
}

func listGames() []model.Game {
	expr, err := expression.NewBuilder().Build()
	if err != nil {
		fmt.Println(err)
	}
	games, err := svc.ListGamesWithPassword(tables.GAMES, expr)
	if err != nil {
		log.Println("listGames error:", err)
	}
	return games
}

func listGamesWithStatus(status string) []model.Game {
	filt := expression.Name("status").Equal(expression.Value(status))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		fmt.Println(err)
	}
	games, err := svc.ListGamesWithPassword(tables.GAMES, expr)
	if err != nil {
		log.Println("listGamesWithStatus error:", err)
	}
	return games
}

func deleteGame(gameMap map[string]interface{}, operator string) (svc.Response, error) {
	gameMap["status"] = "deleted"
	gameMap["updatedBy"] = operator
	games, err := saveGame(gameMap)
	if err != nil {
		return svc.Response{StatusCode: 400, ErrorMessage: err.Error()}, nil
	}
	return svc.Response{StatusCode: 200, Body: games}, nil
}

func addGame(gameMap map[string]interface{}, operator string) []model.Game {
	now := svc.CurrentTimeMillis()
	gameMap["id"] = svc.GenerateKey(10)
	gameMap["roomPassword"] = svc.GenerateKey(8)
	gameMap["status"] = "active"
	gameMap["registeredUserCount"] = 0
	gameMap["insertedAt"] = now
	gameMap["insertedBy"] = operator
	gameMap["updatedAt"] = now
	gameMap["updatedBy"] = operator
	return updateGame(gameMap, operator)
}

func updateGame(gameMap map[string]interface{}, operator string) []model.Game {
	gameMap["updatedAt"] = svc.CurrentTimeMillis()
	gameMap["updatedBy"] = operator
	games, err := saveGame(gameMap)
	if err != nil {
		log.Println("updateGame error:", err)
	}
	return games
}

func saveGame(gameMap map[string]interface{}) ([]model.Game, error) {
	gameId := gameMap["id"].(string)
	status := gameMap["status"].(string)

	if status == "cancelled" || status == "deleted" {
		game := svc.GetGameById(gameId)
		gameUsers := svc.GetGameUsersByGameId(gameId)
		for _, userId := range gameUsers.Users {
			user := svc.GetUserById(userId)
			if err := svc.UnregisterUserToGame(user, game); err != nil {
				log.Printf("saveGame: cannot unregister user %s: %v", user.Username, err)
				return nil, err
			}
		}
	}

	av, err := dynamodbattribute.MarshalMap(gameMap)
	if err != nil {
		return nil, fmt.Errorf("marshal game: %w", err)
	}
	_, err = db.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tables.GAMES),
	})
	if err != nil {
		return nil, fmt.Errorf("PutItem failed: %w", err)
	}
	var game model.Game
	if err = dynamodbattribute.UnmarshalMap(av, &game); err != nil {
		return nil, fmt.Errorf("unmarshal game: %w", err)
	}
	return []model.Game{game}, nil
}

func main() {
	lambda.Start(Handler)
}
