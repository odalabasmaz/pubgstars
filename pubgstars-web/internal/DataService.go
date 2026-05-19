package internal

import (
	Model "../model"
	Tables "../model/tables"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"strings"
)

var (
	db = GetDynamoDbClient("eu-central-1")
)

func RegisterUserToGame(user Model.User, game Model.Game) error {
	log.Println("userGameRegistration...")
	gameId := game.Id
	gamePrice := game.Price
	userId := user.Id

	gameUsers := GetGameUsersByGameId(gameId)
	userGames := GetUserGamesByUserId(userId)

	for _, g := range userGames.Games {
		if gameId == g {
			fmt.Println("already registered")
			return errors.New("oyuna zaten kayitlisiniz")
		}
	}

	if user.GetAvailableBalance() < gamePrice {
		return errors.New("oyun icin yeterli bakiyeniz bulunmamakta")
	} else if !IsGameDateValid(game) {
		return errors.New("oyun tarihi gecerli degil")
	}

	gameUsers.Users = append(gameUsers.Users, userId)
	userGames.Games = append(userGames.Games, gameId)
	game.RegisteredUserCount += game.TeamPlayerCount

	if game.RegisteredUserCount > 100 {
		fmt.Println("oyunda boş yer bulunmamaktadır")
		return errors.New("oyunda 	boş yer bulunmamaktadır")
	}

	game.TotalIncome += gamePrice
	// use bonus first, then balance
	if gamePrice >= user.Bonus {
		leftPrice := gamePrice - user.Bonus
		user.Bonus = 0
		user.Balance -= leftPrice
	} else {
		user.Bonus -= gamePrice
	}

	return userGameRegistration(0, user, game, gameUsers, userGames)
}

func UnregisterUserToGame(user Model.User, game Model.Game) error {
	log.Println("userGameUnregistration...")
	gameId := game.Id
	gamePrice := game.Price
	userId := user.Id
	game.RegisteredUserCount -= game.TeamPlayerCount
	game.TotalIncome -= gamePrice
	user.Bonus += gamePrice

	if !IsGameCancellable(game) {
		log.Println("bu oyunun tarihi gectigi icin iptal edemezsiniz")
		return errors.New("bu oyunun tarihi gectigi icin iptal edemezsiniz")
	}

	gameUsers := GetGameUsersByGameId(gameId)
	userGames := GetUserGamesByUserId(userId)

	if !Contains(userGames.Games, gameId) {
		log.Println("bu oyuna kayitli degilisiniz, userId: " + userId + ", gameId: " + gameId + ", userGames: " + strings.Join(userGames.Games[:], ","))
		return errors.New("bu oyuna kayitli degilisiniz")
	}

	log.Println("before: users > ", gameUsers.Users)
	log.Println("before: games > ", userGames.Games)
	var newGameUsers Model.GameUser
	newGameUsers.GameId = gameId
	for _, id := range gameUsers.Users {
		if id != userId {
			newGameUsers.Users = append(newGameUsers.Users, id)
		}
	}

	var newUserGames Model.UserGame
	newUserGames.UserId = userId
	for _, id := range userGames.Games {
		if id != gameId {
			newUserGames.Games = append(newUserGames.Games, id)
		}
	}
	log.Println("after: users > ", newGameUsers.Users)
	log.Println("after: games > ", newUserGames.Games)

	return userGameRegistration(1, user, game, newGameUsers, newUserGames)
}

func userGameRegistration(registrationType int, user Model.User, game Model.Game, gameUser Model.GameUser, userGame Model.UserGame) error {
	log.Println("game will be saved: ", game)
	log.Println("user will be saved: ", user)
	log.Println("gameUser will be saved: ", gameUser)
	log.Println("userGame will be saved: ", userGame)

	gameMap, err := dynamodbattribute.MarshalMap(game)
	userMap, err := dynamodbattribute.MarshalMap(user)
	gameUserMap, err := dynamodbattribute.MarshalMap(gameUser)
	userGameMap, err := dynamodbattribute.MarshalMap(userGame)

	var registerGame Model.TransactionLog
	if registrationType == 0 {
		registerGame = RegisterGame(user.Id, game.Detail())
	} else if registrationType == 1 {
		registerGame = UnregisterGame(user.Id, game.Detail())
	}
	txMap, err := dynamodbattribute.MarshalMap(registerGame)

	// TODO: just checking the last err
	if err != nil {
		fmt.Println("Got error marshalling map: ", err.Error())
		return err
	}
	_, err2 := db.TransactWriteItems(
		&dynamodb.TransactWriteItemsInput{
			TransactItems: []*dynamodb.TransactWriteItem{
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.GAMES),
						Item:      gameMap,
					},
				},
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.USERS),
						Item:      userMap,
					},
				},
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.GAME_USERS),
						Item:      gameUserMap,
					},
				},
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.USER_GAMES),
						Item:      userGameMap,
					},
				},
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.TRANSACTION_LOG),
						Item:      txMap,
					},
				},
			},
		},
	)

	if err2 != nil {
		fmt.Println("Got error calling transactional write: ", err2.Error())
		return errors.New("beklenmeyen bir hata olustu")
	}
	return nil
}

func CompleteGame(operator string, game Model.Game, firstWinner Model.User, secondWinner Model.User, thirdWinner Model.User) error {
	gameMap, err := dynamodbattribute.MarshalMap(game)
	firstMap, err := dynamodbattribute.MarshalMap(firstWinner)
	secondMap, err := dynamodbattribute.MarshalMap(secondWinner)
	thirdMap, err := dynamodbattribute.MarshalMap(thirdWinner)

	tx1, err := dynamodbattribute.MarshalMap(WinGame(game, operator, firstWinner.Id, game.Award1st))
	tx2, err := dynamodbattribute.MarshalMap(WinGame(game, operator, secondWinner.Id, game.Award2nd))
	tx3, err := dynamodbattribute.MarshalMap(WinGame(game, operator, thirdWinner.Id, game.Award3rd))

	// TODO: just checking the last err
	if err != nil {
		fmt.Println("Got error marshalling map: ", err.Error())
		return err
	}
	_, err2 := db.TransactWriteItems(
		&dynamodb.TransactWriteItemsInput{
			TransactItems: []*dynamodb.TransactWriteItem{
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.GAMES),
						Item:      gameMap,
					},
				},
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.USERS),
						Item:      firstMap,
					},
				},
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.USERS),
						Item:      secondMap,
					},
				},
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.USERS),
						Item:      thirdMap,
					},
				},
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.TRANSACTION_LOG),
						Item:      tx1,
					},
				},
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.TRANSACTION_LOG),
						Item:      tx2,
					},
				},
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.TRANSACTION_LOG),
						Item:      tx3,
					},
				},
			},
		},
	)
	// TODO: just checking the last err
	if err2 != nil {
		fmt.Println("Got error marshalling map: ", err2.Error())
		return errors.New("beklenmeyen bir hata olustu")
	}
	return nil
}

func SaveGame(game Model.Game) error {
	var gameUser Model.GameUser
	gameUser.GameId = game.Id

	log.Println("game will be saved: ", game)
	log.Println("gameUser will be saved: ", gameUser)

	gameMap, err := dynamodbattribute.MarshalMap(game)
	gameUserMap, err := dynamodbattribute.MarshalMap(gameUser)

	// TODO: just checking the last err
	if err != nil {
		fmt.Println("Got error marshalling map: ", err.Error())
		return err
	}
	_, err2 := db.TransactWriteItems(
		&dynamodb.TransactWriteItemsInput{
			TransactItems: []*dynamodb.TransactWriteItem{
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.GAMES),
						Item:      gameMap,
					},
				},
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.GAME_USERS),
						Item:      gameUserMap,
					},
				},
			},
		},
	)

	if err2 != nil {
		fmt.Println("Got error calling transactional write: ", err2.Error())
		return errors.New("beklenmeyen bir hata olustu")
	}
	return nil
}

func SaveUser(user Model.User) error {
	var userGame Model.UserGame
	userGame.UserId = user.Id

	log.Println("user will be saved: ", user)
	log.Println("userGame will be saved: ", userGame)

	userMap, err := dynamodbattribute.MarshalMap(user)
	userGameMap, err := dynamodbattribute.MarshalMap(userGame)

	// TODO: just checking the last err
	if err != nil {
		fmt.Println("Got error marshalling map: ", err.Error())
		return err
	}
	_, err2 := db.TransactWriteItems(
		&dynamodb.TransactWriteItemsInput{
			TransactItems: []*dynamodb.TransactWriteItem{
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.USERS),
						Item:      userMap,
					},
				},
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.USER_GAMES),
						Item:      userGameMap,
					},
				},
			},
		},
	)

	if err2 != nil {
		fmt.Println("Got error calling transactional write: ", err2.Error())
		return errors.New("beklenmeyen bir hata olustu")
	}
	//send slack message
	SendMessageToChannel("customer-registered", "new user registered: "+user.Email)
	return nil
}

func UpdateUserWithTx(user Model.User, tx Model.TransactionLog) error {
	log.Println("user will be updated: ", user)

	userMap, err := dynamodbattribute.MarshalMap(user)
	// TODO: just checking the last err
	if err != nil {
		fmt.Println("Got error marshalling map: ", err.Error())
		return err
	}

	txMap, err := dynamodbattribute.MarshalMap(tx)
	// TODO: just checking the last err
	if err != nil {
		fmt.Println("Got error marshalling map: ", err.Error())
		return err
	}

	_, err2 := db.TransactWriteItems(
		&dynamodb.TransactWriteItemsInput{
			TransactItems: []*dynamodb.TransactWriteItem{
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.USERS),
						Item:      userMap,
					},
				}, {
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.TRANSACTION_LOG),
						Item:      txMap,
					},
				},
			},
		},
	)

	if err2 != nil {
		fmt.Println("Got error calling transactional write: ", err2.Error())
		return errors.New("beklenmeyen bir hata olustu")
	}
	return nil
}

func SaveTransactionLog(tx Model.TransactionLog) error {
	log.Println("tx will be saved: ", tx)
	txMap, err := dynamodbattribute.MarshalMap(tx)
	if err != nil {
		fmt.Println("Got error marshalling map: ", err.Error())
		return err
	}

	_, err2 := db.TransactWriteItems(
		&dynamodb.TransactWriteItemsInput{
			TransactItems: []*dynamodb.TransactWriteItem{
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.TRANSACTION_LOG),
						Item:      txMap,
					},
				},
			},
		},
	)
	if err2 != nil {
		fmt.Println("Got error calling transactional write: ", err2.Error())
		return errors.New("beklenmeyen bir hata olustu")
	}
	return nil
}

func SaveMessage(messageMap map[string]interface{}) error {
	msgMap, err := dynamodbattribute.MarshalMap(messageMap)
	if err != nil {
		fmt.Println("Got error marshalling map: ", err.Error())
		return errors.New("Got error marshalling map: " + err.Error())
	}
	_, err2 := db.TransactWriteItems(
		&dynamodb.TransactWriteItemsInput{
			TransactItems: []*dynamodb.TransactWriteItem{
				{
					Put: &dynamodb.Put{
						TableName: aws.String(Tables.MESSAGES),
						Item:      msgMap,
					},
				},
			},
		},
	)

	if err2 != nil {
		fmt.Println("Got error calling transactional write: ", err2.Error())
		return errors.New("beklenmeyen bir hata olustu")
	}
	return nil
}
