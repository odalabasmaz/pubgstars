package internal

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
)

// DynamoStore is the production Store implementation backed by DynamoDB.
// It delegates to the package-level functions so existing code is unchanged.
type DynamoStore struct{}

func NewDynamoStore() Store { return &DynamoStore{} }

func (s *DynamoStore) GetGameById(id string) model.Game        { return GetGameById(id) }
func (s *DynamoStore) GetUserById(id string) model.User        { return GetUserById(id) }
func (s *DynamoStore) GetUserByEmail(email string) model.User  { return GetUserByEmail(email) }
func (s *DynamoStore) UserExistsByUsername(u string) bool      { return UserExistsByUsername(u) }
func (s *DynamoStore) GetGameUsersByGameId(id string) model.GameUser {
	return GetGameUsersByGameId(id)
}
func (s *DynamoStore) GetUserGamesByUserId(id string) model.UserGame {
	return GetUserGamesByUserId(id)
}
func (s *DynamoStore) ListGames(t string, e expression.Expression) ([]model.Game, error) {
	return ListGames(t, e)
}
func (s *DynamoStore) ListGamesWithPassword(t string, e expression.Expression) ([]model.Game, error) {
	return ListGamesWithPassword(t, e)
}
func (s *DynamoStore) ListUsers() ([]model.User, error) { return ListUsers() }
func (s *DynamoStore) GetUserMap() (map[string]string, error)  { return GetUserMap() }
func (s *DynamoStore) GetTransactionLogs(e expression.Expression) ([]model.TransactionLog, error) {
	return GetTransactionLogs(e)
}
func (s *DynamoStore) GetUserDetails(e expression.Expression) (model.User, error) {
	return GetUserDetails(e)
}
func (s *DynamoStore) ListMessages(e expression.Expression) ([]model.Message, error) {
	return ListMessages(e)
}

func (s *DynamoStore) RegisterUserToGame(u model.User, g model.Game) error {
	return RegisterUserToGame(u, g)
}
func (s *DynamoStore) UnregisterUserToGame(u model.User, g model.Game) error {
	return UnregisterUserToGame(u, g)
}
func (s *DynamoStore) CompleteGame(op string, g model.Game, a, b, c model.User) error {
	return CompleteGame(op, g, a, b, c)
}
func (s *DynamoStore) SaveGame(g model.Game) error                        { return SaveGame(g) }
func (s *DynamoStore) SaveUser(u model.User) error                        { return SaveUser(u) }
func (s *DynamoStore) UpdateUserWithTx(u model.User, t model.TransactionLog) error {
	return UpdateUserWithTx(u, t)
}
func (s *DynamoStore) SaveTransactionLog(t model.TransactionLog) error { return SaveTransactionLog(t) }
func (s *DynamoStore) SaveMessage(m map[string]interface{}) error       { return SaveMessage(m) }
