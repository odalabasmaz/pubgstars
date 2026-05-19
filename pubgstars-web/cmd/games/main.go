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
	log.Println("games handler:", event.Context["http-method"])
	switch event.Context["http-method"] {
	case "GET":
		return svc.Response{StatusCode: 200, Body: listActiveGames()}, nil
	case "PUT", "POST":
		email := svc.GetUsernameFromJwtToken(event.Params["header"]["Authorization"])
		if event.Body["id"] == nil {
			return svc.Response{StatusCode: 200, Body: addGame(event.Body, email)}, nil
		}
		return svc.Response{StatusCode: 200, Body: updateGame(event.Body, email)}, nil
	case "DELETE":
		return svc.Response{StatusCode: 200, Body: deleteGame(event.Body)}, nil
	default:
		return svc.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + event.Context["http-method"]}, nil
	}
}

func listActiveGames() []model.Game {
	filt := expression.Name("status").Equal(expression.Value("active"))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		fmt.Println(err)
	}
	games, err := svc.ListGames(tables.GAMES, expr)
	if err != nil {
		log.Println("listActiveGames error:", err)
	}
	return games
}

func deleteGame(gameMap map[string]interface{}) []model.Game {
	gameMap["status"] = "deleted"
	games, err := saveGame(gameMap)
	if err != nil {
		log.Println("deleteGame error:", err)
	}
	return games
}

func addGame(gameMap map[string]interface{}, email string) []model.Game {
	var game model.Game
	game.Id = svc.GenerateKey(10)
	game.RoomPassword = svc.GenerateKey(8)
	game.Status = "active"
	game.RegisteredUserCount = 0
	now := svc.CurrentTimeMillis()
	game.InsertedAt = now
	game.InsertedBy = email
	game.UpdatedAt = now
	game.UpdatedBy = email
	game.GameDate = gameMap["gameDate"].(string)
	game.League = gameMap["league"].(string)
	game.Type = gameMap["type"].(string)
	game.Map = gameMap["map"].(string)
	game.Price = gameMap["price"].(float64)

	if err := svc.SaveGame(game); err != nil {
		log.Println("addGame error:", err)
		return nil
	}
	return []model.Game{game}
}

func updateGame(gameMap map[string]interface{}, email string) []model.Game {
	gameMap["updatedAt"] = svc.CurrentTimeMillis()
	gameMap["updatedBy"] = email
	games, err := saveGame(gameMap)
	if err != nil {
		log.Println("updateGame error:", err)
	}
	return games
}

func saveGame(gameMap map[string]interface{}) ([]model.Game, error) {
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
