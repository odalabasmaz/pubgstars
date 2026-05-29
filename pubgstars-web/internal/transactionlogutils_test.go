package internal

import (
	"fmt"
	"strings"
	"testing"

	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
)

func TestRegisterGame(t *testing.T) {
	tx := RegisterGame("user1", "gold/erangel/solo/pc @202501011200")
	if tx.UserId != "user1" {
		t.Errorf("UserId = %q, want %q", tx.UserId, "user1")
	}
	if tx.TransactionType != model.GAME {
		t.Errorf("TransactionType = %v, want GAME", tx.TransactionType)
	}
	if tx.SubTransactionType != model.GAME_REGISTER {
		t.Errorf("SubTransactionType = %v, want GAME_REGISTER", tx.SubTransactionType)
	}
	if tx.Operator != "system" {
		t.Errorf("Operator = %q, want %q", tx.Operator, "system")
	}
	if len(tx.Id) != 10 {
		t.Errorf("Id length = %d, want 10", len(tx.Id))
	}
	if tx.InsertedAt <= 0 {
		t.Error("InsertedAt should be a positive timestamp")
	}
}

func TestUnregisterGame(t *testing.T) {
	tx := UnregisterGame("user1", "detail")
	if tx.SubTransactionType != model.GAME_UNREGISTER {
		t.Errorf("SubTransactionType = %v, want GAME_UNREGISTER", tx.SubTransactionType)
	}
}

func TestWithdrawMoney(t *testing.T) {
	tx := WithdrawMoney("user1", "TR123456", "500")
	if tx.TransactionType != model.BALANCE {
		t.Errorf("TransactionType = %v, want BALANCE", tx.TransactionType)
	}
	if tx.SubTransactionType != model.BALANCE_WITHDRAW {
		t.Errorf("SubTransactionType = %v, want BALANCE_WITHDRAW", tx.SubTransactionType)
	}
	if !strings.Contains(tx.Detail, "TR123456") {
		t.Errorf("Detail missing IBAN: %q", tx.Detail)
	}
	if !strings.Contains(tx.Detail, "500") {
		t.Errorf("Detail missing amount: %q", tx.Detail)
	}
}

func TestDepositMoney(t *testing.T) {
	tx := DepositMoney("user1", "250", "bank transfer")
	if tx.SubTransactionType != model.BALANCE_DEPOSIT {
		t.Errorf("SubTransactionType = %v, want BALANCE_DEPOSIT", tx.SubTransactionType)
	}
	if !strings.Contains(tx.Detail, "250") {
		t.Errorf("Detail missing amount: %q", tx.Detail)
	}
	if !strings.Contains(tx.Detail, "bank transfer") {
		t.Errorf("Detail missing description: %q", tx.Detail)
	}
}

func TestAddBalance_BalanceAndBonus(t *testing.T) {
	tx := AddBalance("admin", "user1", 100, 50)
	if !strings.Contains(tx.Detail, "100.00") {
		t.Errorf("Detail missing balance: %q", tx.Detail)
	}
	if !strings.Contains(tx.Detail, "50.00") {
		t.Errorf("Detail missing bonus: %q", tx.Detail)
	}
}

func TestAddBalance_BalanceOnly(t *testing.T) {
	tx := AddBalance("admin", "user1", 100, 0)
	if !strings.Contains(tx.Detail, "100.00") {
		t.Errorf("Detail missing balance: %q", tx.Detail)
	}
	if strings.Contains(tx.Detail, "bonus") {
		t.Errorf("Detail should not mention bonus when bonus=0: %q", tx.Detail)
	}
}

func TestAddBalance_BonusOnly(t *testing.T) {
	tx := AddBalance("admin", "user1", 0, 75)
	if !strings.Contains(tx.Detail, "75.00") {
		t.Errorf("Detail missing bonus: %q", tx.Detail)
	}
}

func TestWinGame(t *testing.T) {
	game := model.Game{
		League:   "gold",
		Map:      "erangel",
		Type:     "solo",
		Platform: "pc",
		GameDate: "202501011200",
	}
	tx := WinGame(game, "admin", "user1", 500)
	if tx.SubTransactionType != model.GAME_WIN {
		t.Errorf("SubTransactionType = %v, want GAME_WIN", tx.SubTransactionType)
	}
	if !strings.Contains(tx.Detail, fmt.Sprintf("%.2f", 500.0)) {
		t.Errorf("Detail missing award amount: %q", tx.Detail)
	}
	if !strings.Contains(tx.Detail, game.Detail()) {
		t.Errorf("Detail missing game detail: %q", tx.Detail)
	}
}

func TestBaseTransactionWithOperator(t *testing.T) {
	tx := BaseTransactionWithOperator("user1", "operator1")
	if tx.Operator != "operator1" {
		t.Errorf("Operator = %q, want %q", tx.Operator, "operator1")
	}
	if tx.UserId != "user1" {
		t.Errorf("UserId = %q, want %q", tx.UserId, "user1")
	}
}
