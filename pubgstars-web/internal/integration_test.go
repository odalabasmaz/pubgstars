//go:build integration

package internal

import (
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model/tables"
)

const localEndpoint = "http://localhost:8001"

// localDB creates a DynamoDB client pointed at the local instance.
func localDB(t *testing.T) *dynamodb.DynamoDB {
	t.Helper()
	resp, err := http.Get(localEndpoint)
	if err != nil || resp.StatusCode != 400 {
		// DynamoDB Local returns 400 on a bare GET — anything else means it's not up
		t.Skip("DynamoDB Local not available — run: docker-compose -f docker-compose.test.yml up -d")
	}

	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("eu-central-1"),
		Endpoint:    aws.String(localEndpoint),
		Credentials: credentials.NewStaticCredentials("dummy", "dummy", ""),
	})
	return dynamodb.New(sess)
}

func createTable(t *testing.T, client *dynamodb.DynamoDB, name, hashKey, hashType string) {
	t.Helper()
	_, err := client.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(name),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String(hashKey), AttributeType: aws.String(hashType)},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String(hashKey), KeyType: aws.String("HASH")},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
	})
	if err != nil {
		t.Logf("CreateTable %s: %v (may already exist)", name, err)
	}
}

// setupTables creates all required tables and overrides the package-level db.
func setupTables(t *testing.T) *dynamodb.DynamoDB {
	t.Helper()
	client := localDB(t)

	createTable(t, client, tables.GAMES, "id", "S")
	createTable(t, client, tables.USERS, "id", "S")
	createTable(t, client, tables.GAME_USERS, "gameId", "S")
	createTable(t, client, tables.USER_GAMES, "userId", "S")
	createTable(t, client, tables.TRANSACTION_LOG, "id", "S")
	createTable(t, client, tables.MESSAGES, "id", "S")

	// replace the package-level DynamoDB client used by all service functions
	db = client
	return client
}

// ─── Tests ────────────────────────────────────────────────────────────────────

func TestIntegration_SaveAndGetGame(t *testing.T) {
	setupTables(t)

	game := model.Game{
		Id:          "integ-game-1",
		GameDate:    time.Now().Add(2 * time.Hour).In(location).Format("200601021504"),
		League:      "gold",
		Type:        "solo",
		Map:         "erangel",
		Platform:    "pc",
		Status:      "active",
		Price:       50,
		InsertedAt:  CurrentTimeMillis(),
		InsertedBy:  "test",
		UpdatedAt:   CurrentTimeMillis(),
		UpdatedBy:   "test",
	}

	if err := SaveGame(game); err != nil {
		t.Fatalf("SaveGame: %v", err)
	}

	got := GetGameById("integ-game-1")
	if got.Id != game.Id {
		t.Errorf("GetGameById: got %q, want %q", got.Id, game.Id)
	}
	if got.Price != game.Price {
		t.Errorf("Price: got %.2f, want %.2f", got.Price, game.Price)
	}
}

func TestIntegration_SaveAndGetUser(t *testing.T) {
	setupTables(t)

	user := model.User{
		Id:         "integ-user-1",
		Username:   "testuser",
		Email:      "test@example.com",
		Status:     "active",
		Balance:    200,
		Bonus:      50,
		InsertedAt: CurrentTimeMillis(),
		InsertedBy: "test",
		UpdatedAt:  CurrentTimeMillis(),
		UpdatedBy:  "test",
	}

	if err := SaveUser(user); err != nil {
		t.Fatalf("SaveUser: %v", err)
	}

	byId := GetUserById("integ-user-1")
	if byId.Id != user.Id {
		t.Errorf("GetUserById: got %q, want %q", byId.Id, user.Id)
	}

	byEmail := GetUserByEmail("test@example.com")
	if byEmail.Email != user.Email {
		t.Errorf("GetUserByEmail: got %q, want %q", byEmail.Email, user.Email)
	}
}

func TestIntegration_RegisterAndUnregisterUser(t *testing.T) {
	setupTables(t)

	user := model.User{
		Id:         "integ-user-2",
		Username:   "player",
		Email:      "player@example.com",
		Status:     "active",
		Balance:    300,
		InsertedAt: CurrentTimeMillis(),
		InsertedBy: "test",
		UpdatedAt:  CurrentTimeMillis(),
		UpdatedBy:  "test",
	}
	game := model.Game{
		Id:              "integ-game-2",
		GameDate:        time.Now().Add(3 * time.Hour).In(location).Format("200601021504"),
		Price:           100,
		TeamPlayerCount: 1,
		Status:          "active",
		InsertedAt:      CurrentTimeMillis(),
		InsertedBy:      "test",
		UpdatedAt:       CurrentTimeMillis(),
		UpdatedBy:       "test",
	}

	if err := SaveUser(user); err != nil {
		t.Fatalf("SaveUser: %v", err)
	}
	if err := SaveGame(game); err != nil {
		t.Fatalf("SaveGame: %v", err)
	}

	if err := RegisterUserToGame(user, game); err != nil {
		t.Fatalf("RegisterUserToGame: %v", err)
	}

	userGame := GetUserGamesByUserId(user.Id)
	if !Contains(userGame.Games, game.Id) {
		t.Errorf("user not registered to game after RegisterUserToGame")
	}

	// re-fetch updated user for unregister (balance was deducted)
	updatedUser := GetUserById(user.Id)
	updatedGame := GetGameById(game.Id)

	if err := UnregisterUserToGame(updatedUser, updatedGame); err != nil {
		t.Fatalf("UnregisterUserToGame: %v", err)
	}

	userGameAfter := GetUserGamesByUserId(user.Id)
	if Contains(userGameAfter.Games, game.Id) {
		t.Errorf("user still registered after UnregisterUserToGame")
	}
}

func TestIntegration_UpdateUserWithTx(t *testing.T) {
	setupTables(t)

	user := model.User{
		Id:         "integ-user-3",
		Balance:    500,
		InsertedAt: CurrentTimeMillis(),
		InsertedBy: "test",
		UpdatedAt:  CurrentTimeMillis(),
		UpdatedBy:  "test",
	}
	if err := SaveUser(user); err != nil {
		t.Fatalf("SaveUser: %v", err)
	}

	user.Balance = 400
	tx := WithdrawMoney(user.Id, "TR000", "100")

	if err := UpdateUserWithTx(user, tx); err != nil {
		t.Fatalf("UpdateUserWithTx: %v", err)
	}

	updated := GetUserById(user.Id)
	if updated.Balance != 400 {
		t.Errorf("expected balance 400, got %.2f", updated.Balance)
	}
}
