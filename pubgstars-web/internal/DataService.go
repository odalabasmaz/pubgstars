package internal

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model/tables"
)

var db = GetDynamoDbClient("eu-central-1")

func RegisterUserToGame(user model.User, game model.Game) error {
	log.Println("userGameRegistration...")
	gameUsers := GetGameUsersByGameId(game.Id)
	userGames := GetUserGamesByUserId(user.Id)

	for _, g := range userGames.Games {
		if game.Id == g {
			return errors.New("oyuna zaten kayitlisiniz")
		}
	}

	if user.GetAvailableBalance() < game.Price {
		return errors.New("oyun icin yeterli bakiyeniz bulunmamakta")
	}
	if !IsGameDateValid(game) {
		return errors.New("oyun tarihi gecerli degil")
	}

	gameUsers.Users = append(gameUsers.Users, user.Id)
	userGames.Games = append(userGames.Games, game.Id)
	game.RegisteredUserCount += game.TeamPlayerCount

	if game.RegisteredUserCount > 100 {
		return errors.New("oyunda boş yer bulunmamaktadır")
	}

	game.TotalIncome += game.Price
	if game.Price >= user.Bonus {
		leftPrice := game.Price - user.Bonus
		user.Bonus = 0
		user.Balance -= leftPrice
	} else {
		user.Bonus -= game.Price
	}

	return userGameRegistration(0, user, game, gameUsers, userGames)
}

func UnregisterUserToGame(user model.User, game model.Game) error {
	log.Println("userGameUnregistration...")
	if !IsGameCancellable(game) {
		return errors.New("bu oyunun tarihi gectigi icin iptal edemezsiniz")
	}

	gameUsers := GetGameUsersByGameId(game.Id)
	userGames := GetUserGamesByUserId(user.Id)

	if !Contains(userGames.Games, game.Id) {
		log.Printf("bu oyuna kayitli degilisiniz, userId: %s, gameId: %s, userGames: %s",
			user.Id, game.Id, strings.Join(userGames.Games, ","))
		return errors.New("bu oyuna kayitli degilisiniz")
	}

	game.RegisteredUserCount -= game.TeamPlayerCount
	game.TotalIncome -= game.Price
	user.Bonus += game.Price

	var newGameUsers model.GameUser
	newGameUsers.GameId = game.Id
	for _, id := range gameUsers.Users {
		if id != user.Id {
			newGameUsers.Users = append(newGameUsers.Users, id)
		}
	}

	var newUserGames model.UserGame
	newUserGames.UserId = user.Id
	for _, id := range userGames.Games {
		if id != game.Id {
			newUserGames.Games = append(newUserGames.Games, id)
		}
	}

	return userGameRegistration(1, user, game, newGameUsers, newUserGames)
}

func userGameRegistration(registrationType int, user model.User, game model.Game, gameUser model.GameUser, userGame model.UserGame) error {
	gameMap, err := dynamodbattribute.MarshalMap(game)
	if err != nil {
		return fmt.Errorf("marshal game: %w", err)
	}
	userMap, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("marshal user: %w", err)
	}
	gameUserMap, err := dynamodbattribute.MarshalMap(gameUser)
	if err != nil {
		return fmt.Errorf("marshal gameUser: %w", err)
	}
	userGameMap, err := dynamodbattribute.MarshalMap(userGame)
	if err != nil {
		return fmt.Errorf("marshal userGame: %w", err)
	}

	var tx model.TransactionLog
	if registrationType == 0 {
		tx = RegisterGame(user.Id, game.Detail())
	} else {
		tx = UnregisterGame(user.Id, game.Detail())
	}
	txMap, err := dynamodbattribute.MarshalMap(tx)
	if err != nil {
		return fmt.Errorf("marshal tx: %w", err)
	}

	_, err = db.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{Put: &dynamodb.Put{TableName: aws.String(tables.GAMES), Item: gameMap}},
			{Put: &dynamodb.Put{TableName: aws.String(tables.USERS), Item: userMap}},
			{Put: &dynamodb.Put{TableName: aws.String(tables.GAME_USERS), Item: gameUserMap}},
			{Put: &dynamodb.Put{TableName: aws.String(tables.USER_GAMES), Item: userGameMap}},
			{Put: &dynamodb.Put{TableName: aws.String(tables.TRANSACTION_LOG), Item: txMap}},
		},
	})
	if err != nil {
		return fmt.Errorf("transact write failed: %w", err)
	}
	return nil
}

func CompleteGame(operator string, game model.Game, first, second, third model.User) error {
	gameMap, err := dynamodbattribute.MarshalMap(game)
	if err != nil {
		return fmt.Errorf("marshal game: %w", err)
	}
	firstMap, err := dynamodbattribute.MarshalMap(first)
	if err != nil {
		return fmt.Errorf("marshal first winner: %w", err)
	}
	secondMap, err := dynamodbattribute.MarshalMap(second)
	if err != nil {
		return fmt.Errorf("marshal second winner: %w", err)
	}
	thirdMap, err := dynamodbattribute.MarshalMap(third)
	if err != nil {
		return fmt.Errorf("marshal third winner: %w", err)
	}
	tx1, err := dynamodbattribute.MarshalMap(WinGame(game, operator, first.Id, game.Award1st))
	if err != nil {
		return fmt.Errorf("marshal tx1: %w", err)
	}
	tx2, err := dynamodbattribute.MarshalMap(WinGame(game, operator, second.Id, game.Award2nd))
	if err != nil {
		return fmt.Errorf("marshal tx2: %w", err)
	}
	tx3, err := dynamodbattribute.MarshalMap(WinGame(game, operator, third.Id, game.Award3rd))
	if err != nil {
		return fmt.Errorf("marshal tx3: %w", err)
	}

	_, err = db.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{Put: &dynamodb.Put{TableName: aws.String(tables.GAMES), Item: gameMap}},
			{Put: &dynamodb.Put{TableName: aws.String(tables.USERS), Item: firstMap}},
			{Put: &dynamodb.Put{TableName: aws.String(tables.USERS), Item: secondMap}},
			{Put: &dynamodb.Put{TableName: aws.String(tables.USERS), Item: thirdMap}},
			{Put: &dynamodb.Put{TableName: aws.String(tables.TRANSACTION_LOG), Item: tx1}},
			{Put: &dynamodb.Put{TableName: aws.String(tables.TRANSACTION_LOG), Item: tx2}},
			{Put: &dynamodb.Put{TableName: aws.String(tables.TRANSACTION_LOG), Item: tx3}},
		},
	})
	if err != nil {
		return fmt.Errorf("transact write failed: %w", err)
	}
	return nil
}

func SaveGame(game model.Game) error {
	var gameUser model.GameUser
	gameUser.GameId = game.Id

	gameMap, err := dynamodbattribute.MarshalMap(game)
	if err != nil {
		return fmt.Errorf("marshal game: %w", err)
	}
	gameUserMap, err := dynamodbattribute.MarshalMap(gameUser)
	if err != nil {
		return fmt.Errorf("marshal gameUser: %w", err)
	}

	_, err = db.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{Put: &dynamodb.Put{TableName: aws.String(tables.GAMES), Item: gameMap}},
			{Put: &dynamodb.Put{TableName: aws.String(tables.GAME_USERS), Item: gameUserMap}},
		},
	})
	if err != nil {
		return fmt.Errorf("transact write failed: %w", err)
	}
	return nil
}

func SaveUser(user model.User) error {
	var userGame model.UserGame
	userGame.UserId = user.Id

	userMap, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("marshal user: %w", err)
	}
	userGameMap, err := dynamodbattribute.MarshalMap(userGame)
	if err != nil {
		return fmt.Errorf("marshal userGame: %w", err)
	}

	_, err = db.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{Put: &dynamodb.Put{TableName: aws.String(tables.USERS), Item: userMap}},
			{Put: &dynamodb.Put{TableName: aws.String(tables.USER_GAMES), Item: userGameMap}},
		},
	})
	if err != nil {
		return fmt.Errorf("transact write failed: %w", err)
	}
	SendMessageToChannel("customer-registered", "new user registered: "+user.Email)
	return nil
}

func UpdateUserWithTx(user model.User, tx model.TransactionLog) error {
	userMap, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("marshal user: %w", err)
	}
	txMap, err := dynamodbattribute.MarshalMap(tx)
	if err != nil {
		return fmt.Errorf("marshal tx: %w", err)
	}

	_, err = db.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{Put: &dynamodb.Put{TableName: aws.String(tables.USERS), Item: userMap}},
			{Put: &dynamodb.Put{TableName: aws.String(tables.TRANSACTION_LOG), Item: txMap}},
		},
	})
	if err != nil {
		return fmt.Errorf("transact write failed: %w", err)
	}
	return nil
}

func SaveTransactionLog(tx model.TransactionLog) error {
	txMap, err := dynamodbattribute.MarshalMap(tx)
	if err != nil {
		return fmt.Errorf("marshal tx: %w", err)
	}
	_, err = db.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{Put: &dynamodb.Put{TableName: aws.String(tables.TRANSACTION_LOG), Item: txMap}},
		},
	})
	if err != nil {
		return fmt.Errorf("transact write failed: %w", err)
	}
	return nil
}

func SaveMessage(messageMap map[string]interface{}) error {
	msgMap, err := dynamodbattribute.MarshalMap(messageMap)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}
	_, err = db.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{Put: &dynamodb.Put{TableName: aws.String(tables.MESSAGES), Item: msgMap}},
		},
	})
	if err != nil {
		return fmt.Errorf("transact write failed: %w", err)
	}
	return nil
}
