package internal

import (
	Model "../model"
	Tables "../model/tables"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"log"
	"os"
	"sort"
	"strings"
	"time"
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

//TODO Delete
var isForTest = os.Getenv("isForTest") == "true"

func GetDynamoDbClient(region string) *dynamodb.DynamoDB {
	if isForTest {
		// for test purpose
		sess, _ := session.NewSessionWithOptions(session.Options{
			Config:  aws.Config{Region: aws.String("eu-central-1")},
			Profile: "pg",
		})
		return dynamodb.New(sess)
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		log.Fatalf("DynamoDB client could not be created because of the err: %v\n", err)
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
	jwtPayload := strings.Split(jwtToken, ".")[1]
	//bytes, err := base64.StdEncoding.DecodeString(jwtPayload+"==")
	bytes, err := base64.RawURLEncoding.DecodeString(jwtPayload)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	jwtPayloadDecoded := string(bytes)
	log.Println(jwtPayloadDecoded)
	datum := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &datum); err != nil {
		log.Fatalf("error %v", err)
	}
	return CovertToString(datum["email"])
}

func GetUsernameFromJwtTokenForAdmin(jwtToken string) string {
	jwtPayload := strings.Split(jwtToken, ".")[1]
	bytes, err := base64.RawURLEncoding.DecodeString(jwtPayload)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	jwtPayloadDecoded := string(bytes)
	log.Println(jwtPayloadDecoded)
	datum := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &datum); err != nil {
		log.Fatalf("error %v", err)
	}
	return CovertToString(datum["cognito:username"])
}

func GetGameById(gameId string) Model.Game {
	log.Println("getGameById: ", gameId)
	/*filt := expression.Name("id").Equal(expression.Value(gameId))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		fmt.Println("failed filter: ", err)
	}*/

	params := &dynamodb.QueryInput{
		TableName:              aws.String(Tables.GAMES),
		KeyConditionExpression: aws.String("id = :gameId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":gameId": {
				S: aws.String(gameId),
			},
		},
	}

	result, err := db.Query(params)
	if err != nil {
		fmt.Println("failed to make Query API call: ", err)
	}

	//todo: check result if game exists

	var game Model.Game
	for _, i := range result.Items {
		err = dynamodbattribute.UnmarshalMap(i, &game)
		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
		} else {
			break
		}
	}

	if err != nil {
		panic(err)
	}

	log.Println("returning: ")
	log.Println(game)

	return game
}
func GetUserMap() (map[string]string, error) {
	var userMap = make(map[string]string)

	params := &dynamodb.ScanInput{
		TableName: aws.String(Tables.USERS),
	}
	result, err := db.Scan(params)
	if err != nil {
		fmt.Println("failed to make Query API call: ", err)
		return nil, err
	}

	if len(result.Items) == 0 {
		return userMap, nil
	}

	for _, i := range result.Items {
		user := Model.User{}

		err = dynamodbattribute.UnmarshalMap(i, &user)
		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
		}

		userMap[user.Id] = user.Username
	}

	return userMap, nil
}

func ListUsers() ([]Model.User, error) {
	log.Println("listing users...")
	params := &dynamodb.ScanInput{
		TableName: aws.String(Tables.USERS),
	}
	result, err := db.Scan(params)
	if err != nil {
		fmt.Println("failed to make Query API call: ", err)
		panic(err)
	}

	var users []Model.User
	if len(result.Items) == 0 {
		return users, nil
	}

	for _, i := range result.Items {
		user := Model.User{}

		err = dynamodbattribute.UnmarshalMap(i, &user)
		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
		}

		user.SecretQuestion = ""
		user.SecretAnswer = ""
		users = append(users, user)
	}

	log.Println("returning users: ")
	log.Println(users)
	return users, err
}

func GetUserById(userId string) Model.User {
	log.Println("getUserById: ", userId)

	params := &dynamodb.QueryInput{
		TableName:              aws.String(Tables.USERS),
		KeyConditionExpression: aws.String("id = :userId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userId": {
				S: aws.String(userId),
			},
		},
	}
	result, err := db.Query(params)
	if err != nil {
		fmt.Println("failed to make Query API call: ", err)
	}

	var user Model.User
	for _, i := range result.Items {
		err = dynamodbattribute.UnmarshalMap(i, &user)
		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
		} else {
			break
		}
	}

	log.Println("returning: ")
	log.Println(user)
	if err != nil {
		fmt.Println("panic!!")
		panic(err)
	}

	return user
}

func GetUserByEmail(email string) Model.User {
	log.Println("getUserByName: ", email)

	filt := expression.Name("email").Equal(expression.Value(email))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		fmt.Println(err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(Tables.USERS),
	}

	result, err := db.Scan(params)
	if err != nil {
		fmt.Println("failed to make Query API call: ", err)
	}

	var user Model.User
	for _, i := range result.Items {
		err = dynamodbattribute.UnmarshalMap(i, &user)
		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
		} else {
			break
		}
	}

	log.Println("returning: ")
	log.Println(user)
	if err != nil {
		fmt.Println("panic!!")
		panic(err)
	}

	return user
}

func UserExistsByUsername(username string) bool {
	log.Println("UserExistsByUsernameOrEmail: ", username)

	filt := expression.Name("username").Equal(expression.Value(username))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		fmt.Println(err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(Tables.USERS),
	}

	result, err := db.Scan(params)
	if err != nil {
		fmt.Println("failed to make Query API call: ", err)
	}
	fmt.Println("result size: ", len(result.Items))

	if len(result.Items) == 0 {
		return true
	}
	return false
}

func GetGameUsersByGameId(gameId string) Model.GameUser {
	log.Println("getGameUsersByGameId: ", gameId)
	params := &dynamodb.QueryInput{
		TableName:              aws.String(Tables.GAME_USERS),
		KeyConditionExpression: aws.String("gameId = :gameId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":gameId": {
				S: aws.String(gameId),
			},
		},
	}

	var gameUser Model.GameUser
	result, err := db.Query(params)
	if err != nil {
		fmt.Println("failed to make Query API call: ", err)
		return gameUser
	} else if len(result.Items) == 0 {
		fmt.Print("no result found for gameUsers")
		gameUser.GameId = gameId
		return gameUser
	}

	for _, i := range result.Items {
		err = dynamodbattribute.UnmarshalMap(i, &gameUser)
		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
		} else {
			break
		}
	}

	log.Println("returning: ")
	log.Println(gameUser)
	if err != nil {
		panic(err)
	}

	return gameUser
}

func GetUserGamesByUserId(userId string) Model.UserGame {
	log.Println("getUserGamesByGameId: ", userId)

	params := &dynamodb.QueryInput{
		TableName:              aws.String(Tables.USER_GAMES),
		KeyConditionExpression: aws.String("userId = :userId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userId": {
				S: aws.String(userId),
			},
		},
	}

	var userGame Model.UserGame
	result, err := db.Query(params)
	if err != nil {
		fmt.Println("failed to make Query API call: ", err)
		return userGame
	} else if len(result.Items) == 0 {
		fmt.Print("no result found for userGames")
		userGame.UserId = userId
		return userGame
	}

	for _, i := range result.Items {
		err = dynamodbattribute.UnmarshalMap(i, &userGame)
		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
		} else {
			break
		}
	}

	if userGame.Games == nil {
		log.Println("")
	}

	log.Println("returning: ")
	log.Println(userGame)
	if err != nil {
		panic(err)
	}

	return userGame
}

// TODO: consider pagination...
func ListGames(tableName string, expr expression.Expression) ([]Model.Game, error) {
	games, err1 := ListGamesWithPassword(tableName, expr)
	userMap, err2 := GetUserMap()

	if err1 == nil && err2 == nil {
		for i, g := range games {
			games[i].RoomPassword = ""
			games[i].Discord = ""
			games[i].Winner1stName = userMap[g.Winner1st]
			games[i].Winner2ndName = userMap[g.Winner2nd]
			games[i].Winner3rdName = userMap[g.Winner3rd]
		}
	}
	return games, errors.New("unexpected error during listing games")
}

func ListGamesWithPassword(tableName string, expr expression.Expression) ([]Model.Game, error) {
	log.Println("listGames...")

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
	}
	result, err := db.Scan(params)
	if err != nil {
		fmt.Println("failed to make Query API call: ", err)
		panic(err)
	}

	var games []Model.Game
	if len(result.Items) == 0 {
		return games, nil
	}

	for _, i := range result.Items {
		game := Model.Game{}

		err = dynamodbattribute.UnmarshalMap(i, &game)
		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
		}

		game.Cancellable = IsGameCancellable(game)
		games = append(games, game)
	}

	sort.Slice(games, func(i, j int) bool {
		return games[i].GameDate < games[j].GameDate
	})

	log.Println("returning games: ")
	log.Println(games)
	return games, err
}

func GetTransactionLogs(expr expression.Expression) ([]Model.TransactionLog, error) {
	log.Println("listTransactionLogs...")
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(Tables.TRANSACTION_LOG),
	}
	result, err := db.Scan(params)
	if err != nil {
		fmt.Println("failed to make Query API call: ", err)
		panic(err)
	}

	var txLogs []Model.TransactionLog
	for _, i := range result.Items {
		tx := Model.TransactionLog{}

		err = dynamodbattribute.UnmarshalMap(i, &tx)
		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
		}

		txLogs = append(txLogs, tx)
	}

	// sort
	sort.Slice(txLogs, func(i, j int) bool {
		return txLogs[i].InsertedAt > txLogs[j].InsertedAt
	})

	log.Println("returning: ")
	log.Println(txLogs)
	return txLogs, err
}

func GetUserDetails(expr expression.Expression) (Model.User, error) {
	log.Println("gettingUserDetails...")
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(Tables.USERS),
	}
	var users []Model.User
	result, err := db.Scan(params)
	if err != nil {
		fmt.Println("failed to make Query API call: ", err)
		return users[0], errors.New("failed to make Query API call")
	}

	for _, i := range result.Items {
		user := Model.User{}

		err = dynamodbattribute.UnmarshalMap(i, &user)
		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
		}

		user.SecretQuestion = ""
		user.SecretAnswer = ""
		users = append(users, user)
	}

	return users[0], err
}

func ListMessages(expr expression.Expression) ([]Model.Message, error) {
	log.Println("listing messages...")
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String("messages"),
	}
	result, err := db.Scan(params)
	if err != nil {
		fmt.Println("failed to make Query API call: ", err)
		return nil, errors.New("failed to make Query API call")
	}

	var messages []Model.Message
	if len(result.Items) == 0 {
		return messages, nil
	}

	for _, i := range result.Items {
		message := Model.Message{}

		err = dynamodbattribute.UnmarshalMap(i, &message)
		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
		}

		messages = append(messages, message)
	}
	return messages, err
}

func Contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
