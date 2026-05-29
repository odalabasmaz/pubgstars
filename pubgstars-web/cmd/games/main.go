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
		games, err := listActiveGames()
		if err != nil {
			log.Println("listActiveGames error:", err)
			return svc.Response{StatusCode: 500, ErrorMessage: "Failed to list games"}, nil
		}
		return svc.Response{StatusCode: 200, Body: games}, nil
	case "PUT", "POST":
		email := svc.GetUsernameFromJwtToken(event.Params["header"]["Authorization"])
		if event.Body["id"] == nil {
			games, err := addGame(event.Body, email)
			if err != nil {
				log.Println("addGame error:", err)
				return svc.Response{StatusCode: 400, ErrorMessage: err.Error()}, nil
			}
			return svc.Response{StatusCode: 200, Body: games}, nil
		}
		games, err := updateGame(event.Body, email)
		if err != nil {
			log.Println("updateGame error:", err)
			return svc.Response{StatusCode: 500, ErrorMessage: "Failed to update game"}, nil
		}
		return svc.Response{StatusCode: 200, Body: games}, nil
	case "DELETE":
		games, err := deleteGame(event.Body)
		if err != nil {
			log.Println("deleteGame error:", err)
			return svc.Response{StatusCode: 500, ErrorMessage: "Failed to delete game"}, nil
		}
		return svc.Response{StatusCode: 200, Body: games}, nil
	default:
		return svc.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + event.Context["http-method"]}, nil
	}
}

func listActiveGames() ([]model.Game, error) {
	filt := expression.Name("status").Equal(expression.Value("active"))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, fmt.Errorf("build expression: %w", err)
	}
	return svc.ListGames(tables.GAMES, expr)
}

func deleteGame(gameMap map[string]interface{}) ([]model.Game, error) {
	gameMap["status"] = "deleted"
	return saveGame(gameMap)
}

func addGame(gameMap map[string]interface{}, email string) ([]model.Game, error) {
	gameDate, ok := gameMap["gameDate"].(string)
	if !ok {
		return nil, fmt.Errorf("gameDate is missing or not a string")
	}
	league, ok := gameMap["league"].(string)
	if !ok {
		return nil, fmt.Errorf("league is missing or not a string")
	}
	gameType, ok := gameMap["type"].(string)
	if !ok {
		return nil, fmt.Errorf("type is missing or not a string")
	}
	gameMap2, ok := gameMap["map"].(string)
	if !ok {
		return nil, fmt.Errorf("map is missing or not a string")
	}
	price, ok := gameMap["price"].(float64)
	if !ok {
		return nil, fmt.Errorf("price is missing or not a number")
	}

	now := svc.CurrentTimeMillis()
	game := model.Game{
		Id:           svc.GenerateKey(10),
		RoomPassword: svc.GenerateKey(8),
		Status:       "active",
		InsertedAt:   now,
		InsertedBy:   email,
		UpdatedAt:    now,
		UpdatedBy:    email,
		GameDate:     gameDate,
		League:       league,
		Type:         gameType,
		Map:          gameMap2,
		Price:        price,
	}

	if err := svc.SaveGame(game); err != nil {
		return nil, err
	}
	return []model.Game{game}, nil
}

func updateGame(gameMap map[string]interface{}, email string) ([]model.Game, error) {
	gameMap["updatedAt"] = svc.CurrentTimeMillis()
	gameMap["updatedBy"] = email
	return saveGame(gameMap)
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
