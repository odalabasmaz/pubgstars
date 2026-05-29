package internal

import (
	"fmt"

	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
)

func BaseTransaction(userId string) model.TransactionLog {
	var tx model.TransactionLog
	tx.Id = GenerateKey(10)
	tx.UserId = userId
	tx.Operator = "system"
	tx.InsertedAt = CurrentTimeMillis()
	return tx
}

func BaseTransactionWithOperator(userId string, operator string) model.TransactionLog {
	var tx model.TransactionLog
	tx.Id = GenerateKey(10)
	tx.UserId = userId
	tx.Operator = operator
	tx.InsertedAt = CurrentTimeMillis()
	return tx
}

func RegisterGame(userId string, gameDetail string) model.TransactionLog {
	tx := BaseTransaction(userId)
	tx.TransactionType = model.GAME
	tx.SubTransactionType = model.GAME_REGISTER
	tx.Detail = gameDetail
	return tx
}

func UnregisterGame(userId string, gameDetail string) model.TransactionLog {
	tx := BaseTransaction(userId)
	tx.TransactionType = model.GAME
	tx.SubTransactionType = model.GAME_UNREGISTER
	tx.Detail = gameDetail
	return tx
}

func WithdrawMoney(userId string, iban string, amount string) model.TransactionLog {
	tx := BaseTransaction(userId)
	tx.TransactionType = model.BALANCE
	tx.SubTransactionType = model.BALANCE_WITHDRAW
	tx.Detail = "Withdrawal request of " + amount + " TL to IBAN " + iban + " has been received"
	return tx
}

func DepositMoney(userId string, amount string, description string) model.TransactionLog {
	tx := BaseTransaction(userId)
	tx.TransactionType = model.BALANCE
	tx.SubTransactionType = model.BALANCE_DEPOSIT
	tx.Detail = "Deposit request of " + amount + " TL for account " + userId + " has been received. Description: " + description
	return tx
}

func AddBalance(operator string, userId string, balance float64, bonus float64) model.TransactionLog {
	tx := BaseTransactionWithOperator(userId, operator)
	tx.TransactionType = model.BALANCE
	tx.SubTransactionType = model.BALANCE_LOAD
	if balance > 0 && bonus > 0 {
		tx.Detail = fmt.Sprintf("%.2f TL balance and %.2f TL bonus have been loaded to your account", balance, bonus)
	} else if balance > 0 {
		tx.Detail = fmt.Sprintf("%.2f TL balance has been loaded to your account", balance)
	} else if bonus > 0 {
		tx.Detail = fmt.Sprintf("%.2f TL bonus has been loaded to your account", bonus)
	}
	return tx
}

func WinGame(game model.Game, operator string, userId string, award float64) model.TransactionLog {
	tx := BaseTransactionWithOperator(userId, operator)
	tx.TransactionType = model.GAME
	tx.SubTransactionType = model.GAME_WIN
	tx.Detail = "You won the game!\n" + game.Detail() + "\n" + fmt.Sprintf("%.2f TL prize has been loaded to your balance.", award)
	return tx
}

func GetTransactionType(transactionType model.TransactionType) string {
	return model.TRANSACTION_MAP[transactionType]
}

func GetSubTransactionType(subTransactionType model.SubTransactionType) string {
	return model.SUB_TRANSACTION_MAP[subTransactionType]
}
