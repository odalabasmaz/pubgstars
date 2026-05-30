package testutil

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
)

// MockStore lets each test field-override only the operations it cares about.
// Any unset Fn field returns its zero value without error.
type MockStore struct {
	GetGameByIdFn           func(id string) model.Game
	GetUserByIdFn           func(id string) model.User
	GetUserByEmailFn        func(email string) model.User
	GetGameUsersByGameIdFn  func(gameId string) model.GameUser
	GetUserGamesByUserIdFn  func(userId string) model.UserGame
	ListGamesFn             func(tableName string, expr expression.Expression) ([]model.Game, error)
	ListGamesWithPasswordFn func(tableName string, expr expression.Expression) ([]model.Game, error)
	ListUsersFn             func() ([]model.User, error)
	GetUserMapFn            func() (map[string]string, error)
	GetTransactionLogsFn    func(expr expression.Expression) ([]model.TransactionLog, error)
	GetUserDetailsFn        func(expr expression.Expression) (model.User, error)
	ListMessagesFn          func(expr expression.Expression) ([]model.Message, error)
	UserExistsByUsernameFn  func(username string) bool

	RegisterUserToGameFn  func(user model.User, game model.Game) error
	UnregisterUserToGameFn func(user model.User, game model.Game) error
	CompleteGameFn        func(operator string, game model.Game, first, second, third model.User) error
	SaveGameFn            func(game model.Game) error
	SaveUserFn            func(user model.User) error
	UpdateUserWithTxFn    func(user model.User, tx model.TransactionLog) error
	SaveTransactionLogFn  func(tx model.TransactionLog) error
	SaveMessageFn         func(messageMap map[string]interface{}) error
}

func (m *MockStore) GetGameById(id string) model.Game {
	if m.GetGameByIdFn != nil {
		return m.GetGameByIdFn(id)
	}
	return model.Game{}
}
func (m *MockStore) GetUserById(id string) model.User {
	if m.GetUserByIdFn != nil {
		return m.GetUserByIdFn(id)
	}
	return model.User{}
}
func (m *MockStore) GetUserByEmail(email string) model.User {
	if m.GetUserByEmailFn != nil {
		return m.GetUserByEmailFn(email)
	}
	return model.User{}
}
func (m *MockStore) GetGameUsersByGameId(gameId string) model.GameUser {
	if m.GetGameUsersByGameIdFn != nil {
		return m.GetGameUsersByGameIdFn(gameId)
	}
	return model.GameUser{}
}
func (m *MockStore) GetUserGamesByUserId(userId string) model.UserGame {
	if m.GetUserGamesByUserIdFn != nil {
		return m.GetUserGamesByUserIdFn(userId)
	}
	return model.UserGame{}
}
func (m *MockStore) ListGames(t string, e expression.Expression) ([]model.Game, error) {
	if m.ListGamesFn != nil {
		return m.ListGamesFn(t, e)
	}
	return nil, nil
}
func (m *MockStore) ListGamesWithPassword(t string, e expression.Expression) ([]model.Game, error) {
	if m.ListGamesWithPasswordFn != nil {
		return m.ListGamesWithPasswordFn(t, e)
	}
	return nil, nil
}
func (m *MockStore) ListUsers() ([]model.User, error) {
	if m.ListUsersFn != nil {
		return m.ListUsersFn()
	}
	return nil, nil
}
func (m *MockStore) GetUserMap() (map[string]string, error) {
	if m.GetUserMapFn != nil {
		return m.GetUserMapFn()
	}
	return map[string]string{}, nil
}
func (m *MockStore) GetTransactionLogs(e expression.Expression) ([]model.TransactionLog, error) {
	if m.GetTransactionLogsFn != nil {
		return m.GetTransactionLogsFn(e)
	}
	return nil, nil
}
func (m *MockStore) GetUserDetails(e expression.Expression) (model.User, error) {
	if m.GetUserDetailsFn != nil {
		return m.GetUserDetailsFn(e)
	}
	return model.User{}, nil
}
func (m *MockStore) ListMessages(e expression.Expression) ([]model.Message, error) {
	if m.ListMessagesFn != nil {
		return m.ListMessagesFn(e)
	}
	return nil, nil
}
func (m *MockStore) UserExistsByUsername(u string) bool {
	if m.UserExistsByUsernameFn != nil {
		return m.UserExistsByUsernameFn(u)
	}
	return false
}
func (m *MockStore) RegisterUserToGame(u model.User, g model.Game) error {
	if m.RegisterUserToGameFn != nil {
		return m.RegisterUserToGameFn(u, g)
	}
	return nil
}
func (m *MockStore) UnregisterUserToGame(u model.User, g model.Game) error {
	if m.UnregisterUserToGameFn != nil {
		return m.UnregisterUserToGameFn(u, g)
	}
	return nil
}
func (m *MockStore) CompleteGame(op string, g model.Game, a, b, c model.User) error {
	if m.CompleteGameFn != nil {
		return m.CompleteGameFn(op, g, a, b, c)
	}
	return nil
}
func (m *MockStore) SaveGame(g model.Game) error {
	if m.SaveGameFn != nil {
		return m.SaveGameFn(g)
	}
	return nil
}
func (m *MockStore) SaveUser(u model.User) error {
	if m.SaveUserFn != nil {
		return m.SaveUserFn(u)
	}
	return nil
}
func (m *MockStore) UpdateUserWithTx(u model.User, t model.TransactionLog) error {
	if m.UpdateUserWithTxFn != nil {
		return m.UpdateUserWithTxFn(u, t)
	}
	return nil
}
func (m *MockStore) SaveTransactionLog(t model.TransactionLog) error {
	if m.SaveTransactionLogFn != nil {
		return m.SaveTransactionLogFn(t)
	}
	return nil
}
func (m *MockStore) SaveMessage(mp map[string]interface{}) error {
	if m.SaveMessageFn != nil {
		return m.SaveMessageFn(mp)
	}
	return nil
}
