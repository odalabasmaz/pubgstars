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
			games, err := listGamesWithStatus(status)
			if err != nil {
				log.Println("listGamesWithStatus error:", err)
				return svc.Response{StatusCode: 500, ErrorMessage: "Failed to list games"}, nil
			}
			return svc.Response{StatusCode: 200, Body: games}, nil
		}
		games, err := listGames()
		if err != nil {
			log.Println("listGames error:", err)
			return svc.Response{StatusCode: 500, ErrorMessage: "Failed to list games"}, nil
		}
		return svc.Response{StatusCode: 200, Body: games}, nil
	case "PUT", "POST":
		if event.Body["id"] == nil {
			games, err := addGame(event.Body, operator)
			if err != nil {
				log.Println("addGame error:", err)
				return svc.Response{StatusCode: 500, ErrorMessage: "Failed to add game"}, nil
			}
			return svc.Response{StatusCode: 200, Body: games}, nil
		}
		games, err := updateGame(event.Body, operator)
		if err != nil {
			log.Println("updateGame error:", err)
			return svc.Response{StatusCode: 500, ErrorMessage: "Failed to update game"}, nil
		}
		return svc.Response{StatusCode: 200, Body: games}, nil
	case "DELETE":
		return deleteGame(event.Body, operator)
	default:
		return svc.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + event.Context["http-method"]}, nil
	}
}

func listGames() ([]model.Game, error) {
	expr, err := expression.NewBuilder().Build()
	if err != nil {
		return nil, fmt.Errorf("build expression: %w", err)
	}
	return svc.ListGamesWithPassword(tables.GAMES, expr)
}

func listGamesWithStatus(status string) ([]model.Game, error) {
	filt := expression.Name("status").Equal(expression.Value(status))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, fmt.Errorf("build expression: %w", err)
	}
	return svc.ListGamesWithPassword(tables.GAMES, expr)
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

func addGame(gameMap map[string]interface{}, operator string) ([]model.Game, error) {
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

func updateGame(gameMap map[string]interface{}, operator string) ([]model.Game, error) {
	gameMap["updatedAt"] = svc.CurrentTimeMillis()
	gameMap["updatedBy"] = operator
	return saveGame(gameMap)
}

func saveGame(gameMap map[string]interface{}) ([]model.Game, error) {
	gameId, ok := gameMap["id"].(string)
	if !ok {
		return nil, fmt.Errorf("game id is missing or not a string")
	}
	status, ok := gameMap["status"].(string)
	if !ok {
		return nil, fmt.Errorf("game status is missing or not a string")
	}

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
