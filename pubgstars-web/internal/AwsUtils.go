package internal

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model/tables"
)

type RequestEvent struct {
	Body    map[string]interface{}       `json:"body-json"`
	Params  map[string]map[string]string `json:"params"`
	Context map[string]string            `json:"context"`
}

type Response struct {
	StatusCode        int                 `json:"statusCode"`
	Headers           map[string]string   `json:"headers"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	ErrorMessage      string              `json:"errorMessage"`
	Body              interface{}         `json:"body"`
	IsBase64Encoded   bool                `json:"isBase64Encoded,omitempty"`
}

func GetDynamoDbClient(region string) *dynamodb.DynamoDB {
	profile := os.Getenv("AWS_PROFILE")
	if profile != "" {
		sess, _ := session.NewSessionWithOptions(session.Options{
			Config:  aws.Config{Region: aws.String(region)},
			Profile: profile,
		})
		return dynamodb.New(sess)
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("DynamoDB client could not be created: %v\n", err)
	}
	return dynamodb.New(sess)
}

func CurrentTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}

func CovertToString(datum interface{}) string {
	return fmt.Sprintf("%v", datum)
}

func ConvertMillisToString(millis int64) string {
	return time.Unix(millis/1000, 0).Format("02.01.2006 15:04")
}

func GetUsernameFromJwtToken(jwtToken string) string {
	parts := strings.Split(jwtToken, ".")
	if len(parts) < 2 {
		log.Printf("malformed JWT token")
		return ""
	}
	bytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Printf("failed to decode JWT: %v", err)
		return ""
	}
	datum := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &datum); err != nil {
		log.Printf("failed to unmarshal JWT payload: %v", err)
		return ""
	}
	return CovertToString(datum["email"])
}

func GetUsernameFromJwtTokenForAdmin(jwtToken string) string {
	parts := strings.Split(jwtToken, ".")
	if len(parts) < 2 {
		log.Printf("malformed JWT token")
		return ""
	}
	bytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Printf("failed to decode JWT: %v", err)
		return ""
	}
	datum := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &datum); err != nil {
		log.Printf("failed to unmarshal JWT payload: %v", err)
		return ""
	}
	return CovertToString(datum["cognito:username"])
}

func GetGameById(gameId string) model.Game {
	log.Println("getGameById:", gameId)
	params := &dynamodb.QueryInput{
		TableName:              aws.String(tables.GAMES),
		KeyConditionExpression: aws.String("id = :gameId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":gameId": {S: aws.String(gameId)},
		},
	}
	result, err := db.Query(params)
	if err != nil {
		log.Printf("GetGameById query failed: %v", err)
		return model.Game{}
	}

	var game model.Game
	for _, i := range result.Items {
		if err = dynamodbattribute.UnmarshalMap(i, &game); err != nil {
			log.Printf("GetGameById unmarshal error: %v", err)
		} else {
			break
		}
	}
	return game
}

func GetUserMap() (map[string]string, error) {
	params := &dynamodb.ScanInput{TableName: aws.String(tables.USERS)}
	result, err := db.Scan(params)
	if err != nil {
		return nil, fmt.Errorf("GetUserMap scan failed: %w", err)
	}

	userMap := make(map[string]string, len(result.Items))
	for _, i := range result.Items {
		var user model.User
		if err = dynamodbattribute.UnmarshalMap(i, &user); err != nil {
			log.Printf("GetUserMap unmarshal error: %v", err)
			continue
		}
		userMap[user.Id] = user.Username
	}
	return userMap, nil
}

func ListUsers() ([]model.User, error) {
	log.Println("listing users...")
	params := &dynamodb.ScanInput{TableName: aws.String(tables.USERS)}
	result, err := db.Scan(params)
	if err != nil {
		return nil, fmt.Errorf("ListUsers scan failed: %w", err)
	}

	var users []model.User
	for _, i := range result.Items {
		var user model.User
		if err = dynamodbattribute.UnmarshalMap(i, &user); err != nil {
			log.Printf("ListUsers unmarshal error: %v", err)
			continue
		}
		user.SecretQuestion = ""
		user.SecretAnswer = ""
		users = append(users, user)
	}
	return users, nil
}

func GetUserById(userId string) model.User {
	log.Println("getUserById:", userId)
	params := &dynamodb.QueryInput{
		TableName:              aws.String(tables.USERS),
		KeyConditionExpression: aws.String("id = :userId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userId": {S: aws.String(userId)},
		},
	}
	result, err := db.Query(params)
	if err != nil {
		log.Printf("GetUserById query failed: %v", err)
		return model.User{}
	}

	var user model.User
	for _, i := range result.Items {
		if err = dynamodbattribute.UnmarshalMap(i, &user); err != nil {
			log.Printf("GetUserById unmarshal error: %v", err)
		} else {
			break
		}
	}
	return user
}

func GetUserByEmail(email string) model.User {
	log.Println("getUserByEmail:", email)
	filt := expression.Name("email").Equal(expression.Value(email))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		log.Printf("GetUserByEmail expression build failed: %v", err)
		return model.User{}
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tables.USERS),
	}
	result, err := db.Scan(params)
	if err != nil {
		log.Printf("GetUserByEmail scan failed: %v", err)
		return model.User{}
	}

	var user model.User
	for _, i := range result.Items {
		if err = dynamodbattribute.UnmarshalMap(i, &user); err != nil {
			log.Printf("GetUserByEmail unmarshal error: %v", err)
		} else {
			break
		}
	}
	return user
}

func UserExistsByUsername(username string) bool {
	log.Println("UserExistsByUsername:", username)
	filt := expression.Name("username").Equal(expression.Value(username))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		log.Printf("UserExistsByUsername expression build failed: %v", err)
		return false
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tables.USERS),
	}
	result, err := db.Scan(params)
	if err != nil {
		log.Printf("UserExistsByUsername scan failed: %v", err)
		return false
	}
	return len(result.Items) == 0
}

func GetGameUsersByGameId(gameId string) model.GameUser {
	log.Println("getGameUsersByGameId:", gameId)
	params := &dynamodb.QueryInput{
		TableName:              aws.String(tables.GAME_USERS),
		KeyConditionExpression: aws.String("gameId = :gameId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":gameId": {S: aws.String(gameId)},
		},
	}

	var gameUser model.GameUser
	gameUser.GameId = gameId

	result, err := db.Query(params)
	if err != nil {
		log.Printf("GetGameUsersByGameId query failed: %v", err)
		return gameUser
	}
	if len(result.Items) == 0 {
		return gameUser
	}

	for _, i := range result.Items {
		if err = dynamodbattribute.UnmarshalMap(i, &gameUser); err != nil {
			log.Printf("GetGameUsersByGameId unmarshal error: %v", err)
		} else {
			break
		}
	}
	return gameUser
}

func GetUserGamesByUserId(userId string) model.UserGame {
	log.Println("getUserGamesByUserId:", userId)
	params := &dynamodb.QueryInput{
		TableName:              aws.String(tables.USER_GAMES),
		KeyConditionExpression: aws.String("userId = :userId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userId": {S: aws.String(userId)},
		},
	}

	var userGame model.UserGame
	userGame.UserId = userId

	result, err := db.Query(params)
	if err != nil {
		log.Printf("GetUserGamesByUserId query failed: %v", err)
		return userGame
	}
	if len(result.Items) == 0 {
		return userGame
	}

	for _, i := range result.Items {
		if err = dynamodbattribute.UnmarshalMap(i, &userGame); err != nil {
			log.Printf("GetUserGamesByUserId unmarshal error: %v", err)
		} else {
			break
		}
	}
	return userGame
}

func ListGames(tableName string, expr expression.Expression) ([]model.Game, error) {
	games, err := ListGamesWithPassword(tableName, expr)
	if err != nil {
		return nil, err
	}
	userMap, err := GetUserMap()
	if err != nil {
		return nil, err
	}
	for i, g := range games {
		games[i].RoomPassword = ""
		games[i].Discord = ""
		games[i].Winner1stName = userMap[g.Winner1st]
		games[i].Winner2ndName = userMap[g.Winner2nd]
		games[i].Winner3rdName = userMap[g.Winner3rd]
	}
	return games, nil
}

func ListGamesWithPassword(tableName string, expr expression.Expression) ([]model.Game, error) {
	log.Println("listGames...")
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
	}
	result, err := db.Scan(params)
	if err != nil {
		return nil, fmt.Errorf("ListGamesWithPassword scan failed: %w", err)
	}

	var games []model.Game
	for _, i := range result.Items {
		var game model.Game
		if err = dynamodbattribute.UnmarshalMap(i, &game); err != nil {
			log.Printf("ListGamesWithPassword unmarshal error: %v", err)
			continue
		}
		game.Cancellable = IsGameCancellable(game)
		games = append(games, game)
	}
	sort.Slice(games, func(i, j int) bool {
		return games[i].GameDate < games[j].GameDate
	})
	return games, nil
}

func GetTransactionLogs(expr expression.Expression) ([]model.TransactionLog, error) {
	log.Println("listTransactionLogs...")
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tables.TRANSACTION_LOG),
	}
	result, err := db.Scan(params)
	if err != nil {
		return nil, fmt.Errorf("GetTransactionLogs scan failed: %w", err)
	}

	var txLogs []model.TransactionLog
	for _, i := range result.Items {
		var tx model.TransactionLog
		if err = dynamodbattribute.UnmarshalMap(i, &tx); err != nil {
			log.Printf("GetTransactionLogs unmarshal error: %v", err)
			continue
		}
		txLogs = append(txLogs, tx)
	}
	sort.Slice(txLogs, func(i, j int) bool {
		return txLogs[i].InsertedAt > txLogs[j].InsertedAt
	})
	return txLogs, nil
}

func GetUserDetails(expr expression.Expression) (model.User, error) {
	log.Println("gettingUserDetails...")
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tables.USERS),
	}
	result, err := db.Scan(params)
	if err != nil {
		return model.User{}, fmt.Errorf("GetUserDetails scan failed: %w", err)
	}
	if len(result.Items) == 0 {
		return model.User{}, errors.New("user not found")
	}

	var user model.User
	if err = dynamodbattribute.UnmarshalMap(result.Items[0], &user); err != nil {
		return model.User{}, fmt.Errorf("GetUserDetails unmarshal failed: %w", err)
	}
	user.SecretQuestion = ""
	user.SecretAnswer = ""
	return user, nil
}

func ListMessages(expr expression.Expression) ([]model.Message, error) {
	log.Println("listing messages...")
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tables.MESSAGES),
	}
	result, err := db.Scan(params)
	if err != nil {
		return nil, fmt.Errorf("ListMessages scan failed: %w", err)
	}

	var messages []model.Message
	for _, i := range result.Items {
		var message model.Message
		if err = dynamodbattribute.UnmarshalMap(i, &message); err != nil {
			log.Printf("ListMessages unmarshal error: %v", err)
			continue
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func Contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
