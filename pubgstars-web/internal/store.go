package internal

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
)

// Store abstracts all persistence operations so handlers can be tested with a mock.
type Store interface {
	// reads
	GetGameById(id string) model.Game
	GetUserById(id string) model.User
	GetUserByEmail(email string) model.User
	GetGameUsersByGameId(gameId string) model.GameUser
	GetUserGamesByUserId(userId string) model.UserGame
	ListGames(tableName string, expr expression.Expression) ([]model.Game, error)
	ListGamesWithPassword(tableName string, expr expression.Expression) ([]model.Game, error)
	ListUsers() ([]model.User, error)
	GetUserMap() (map[string]string, error)
	GetTransactionLogs(expr expression.Expression) ([]model.TransactionLog, error)
	GetUserDetails(expr expression.Expression) (model.User, error)
	ListMessages(expr expression.Expression) ([]model.Message, error)
	UserExistsByUsername(username string) bool

	// writes
	RegisterUserToGame(user model.User, game model.Game) error
	UnregisterUserToGame(user model.User, game model.Game) error
	CompleteGame(operator string, game model.Game, first, second, third model.User) error
	SaveGame(game model.Game) error
	SaveUser(user model.User) error
	UpdateUserWithTx(user model.User, tx model.TransactionLog) error
	SaveTransactionLog(tx model.TransactionLog) error
	SaveMessage(messageMap map[string]interface{}) error
}
