package main

import (
	tables "github.com/odalabasmaz/pubgstars/pubgstars-web/model/tables"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"math/rand"
	"os"
	"time"
)

type Game struct {
	Id           string `json:"id"`
	GameDate     string `json:"gameDate"`
	League       string `json:"league"`
	Type         string `json:"type"`
	Map          string `json:"map"`
	Price        int32  `json:"price"`
	RoomPassword string `json:"roomPassword"`
	Status       string `json:"status"`
	InsertedAt   int64  `json:"insertedAt"`
	InsertedBy   string `json:"insertedBy"`
	UpdatedAt    int64  `json:"updatedAt"`
	UpdatedBy    string `json:"updatedBy"`
}

type User struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	SecretPhrase string `json:"secretPhrase"`
	Status       string `json:"status"`
	InsertedAt   int64  `json:"insertedAt"`
	InsertedBy   string `json:"insertedBy"`
	UpdatedAt    int64  `json:"updatedAt"`
	UpdatedBy    string `json:"updatedBy"`
}

func main() {
	sess, _ := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String("eu-central-1")},
		Profile: "pg",
	})

	db := dynamodb.New(sess)
	//createGamesTable(db)
	//addNewItem(db)
	//listGames(db)

	//createTableWithHashKeyAndRangeKey(db, tables.GAMES, "id", "gameDate")
	//createTableWithHashKey(db, tables.GAMES, "id")
	//createTableWithHashKey(db, "userRegistered", "id")
	//createTableWithHashKey(db, tables.GAME_USERS, "gameId")
	//createTableWithHashKey(db, tables.USER_GAMES, "userId")
	//createTableWithHashKeyAndRangeKey(db, tables.TRANSACTION_LOG, "id", "userId")
	createTableWithHashKeyAndRangeKey(db, tables.MESSAGES, "id", "from")

	//addNewItem(db)
}

// TODO: consider pagination...
func listGames(db *dynamodb.DynamoDB) {
	params := &dynamodb.ScanInput{
		TableName: aws.String(tables.GAMES),
	}
	result, err := db.Scan(params)
	if err != nil {
		fmt.Println("failed to make Query API call: ", err)
	}

	var games []Game

	for _, i := range result.Items {
		game := Game{}

		err = dynamodbattribute.UnmarshalMap(i, &game)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		games = append(games, game)
	}

	for g := range games {
		fmt.Println(games[g])
	}
}

func listGames2(db *dynamodb.DynamoDB) {
	filt := expression.Name("state").Equal(expression.Value("active"))
	proj := expression.NamesList(
		expression.Name("id"),
		expression.Name("gameDate"),
		expression.Name("league"),
		expression.Name("type"),
		expression.Name("map"),
		expression.Name("price"),
		expression.Name("roomPassword"),
		expression.Name("state"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		fmt.Println(err)
	}

	params := &dynamodb.ScanInput{
		FilterExpression:     expr.Filter(),
		ProjectionExpression: expr.Projection(),
		TableName:            aws.String(tables.GAMES),
	}

	// Make the DynamoDB Query API call
	result, err := db.Scan(params)

	for _, i := range result.Items {
		game := Game{}

		err = dynamodbattribute.UnmarshalMap(i, &game)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println(dynamodbattribute.MarshalMap(game))
	}
}

func addNewItem(client *dynamodb.DynamoDB) {
	var jsonString = []byte(`{"gameDate":"20190101","league":"A","type":"2","map":"M2","price":10,"state":"active"}`)
	var game Game
	if err := json.Unmarshal(jsonString, &game); err != nil {
		fmt.Printf("error %v", err)
	}
	game.Id = GenerateKey(10)
	game.RoomPassword = "pass"
	addItem(client, game)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateKey(n int) string {
	var source = rand.NewSource(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	key := make([]byte, n)
	for i := range key {
		key[i] = charset[source.Int63()%int64(len(charset))]
	}
	return string(key)
}

func addItem(db *dynamodb.DynamoDB, game Game) {
	av, err := dynamodbattribute.MarshalMap(game)

	if err != nil {
		fmt.Println("Got error marshalling map:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create item in table Movies
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tables.GAMES),
	}

	// todo put item unique
	_, err = db.PutItem(input)

	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func createTableWithHashKeyAndRangeKey(db *dynamodb.DynamoDB, tableName string, hashKey string, rangeKey string) {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(hashKey),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String(rangeKey),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(hashKey),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String(rangeKey),
				KeyType:       aws.String("RANGE"),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
		TableName:   aws.String(tableName),
	}

	var err error
	_, err = db.CreateTable(input)
	if err != nil {
		fmt.Println("Got error calling CreateTable:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Created the table: " + tableName)
}

func createTableWithHashKey(db *dynamodb.DynamoDB, tableName string, hashKey string) {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(hashKey),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(hashKey),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
		TableName:   aws.String(tableName),
	}

	var err error
	_, err = db.CreateTable(input)
	if err != nil {
		fmt.Println("Got error calling CreateTable:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Created the table: " + tableName)
}
