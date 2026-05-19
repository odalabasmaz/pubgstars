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
	tx.Detail = iban + " numaralı hesaba " + amount + " TL aktarım talebi alınmıştır"
	return tx
}

func DepositMoney(userId string, amount string, description string) model.TransactionLog {
	tx := BaseTransaction(userId)
	tx.TransactionType = model.BALANCE
	tx.SubTransactionType = model.BALANCE_DEPOSIT
	tx.Detail = userId + " numaralı hesaba " + amount + " TL yatırım talebi alınmıştır. Açıklama: " + description
	return tx
}

func AddBalance(operator string, userId string, balance float64, bonus float64) model.TransactionLog {
	tx := BaseTransactionWithOperator(userId, operator)
	tx.TransactionType = model.BALANCE
	tx.SubTransactionType = model.BALANCE_LOAD
	if balance > 0 && bonus > 0 {
		tx.Detail = "Hesabiniza " + fmt.Sprintf("%.2f", balance) + " TL bakiye ve " + fmt.Sprintf("%.2f", bonus) + " TL bonus yuklenmistir"
	} else if balance > 0 {
		tx.Detail = "Hesabiniza " + fmt.Sprintf("%.2f", balance) + " TL bakiye yuklenmistir"
	} else if bonus > 0 {
		tx.Detail = "Hesabiniza " + fmt.Sprintf("%.2f", bonus) + " TL bonus yuklenmistir"
	}
	return tx
}

func WinGame(game model.Game, operator string, userId string, award float64) model.TransactionLog {
	tx := BaseTransactionWithOperator(userId, operator)
	tx.TransactionType = model.GAME
	tx.SubTransactionType = model.GAME_WIN
	tx.Detail = "Oyunu kazandiniz!\n" + game.Detail() + "\n" + fmt.Sprintf("%.2f", award) + " TL odul bakiyenize yuklendi."
	return tx
}

func GetTransactionType(transactionType model.TransactionType) string {
	return model.TRANSACTION_MAP[transactionType]
}

func GetSubTransactionType(subTransactionType model.SubTransactionType) string {
	return model.SUB_TRANSACTION_MAP[subTransactionType]
}
