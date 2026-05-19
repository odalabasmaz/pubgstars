package internal

import (
	Model "../model"
	"errors"
	"fmt"
	"strconv"
)

func BaseTransaction(userId string) Model.TransactionLog {
	var tx Model.TransactionLog
	tx.Id = GenerateKey(10)
	tx.UserId = userId
	tx.Operator = "system"
	tx.InsertedAt = CurrentTimeMillis()
	return tx
}

func BaseTransactionWithOperator(userId string, operator string) Model.TransactionLog {
	var tx Model.TransactionLog
	tx.Id = GenerateKey(10)
	tx.UserId = userId
	tx.Operator = operator
	tx.InsertedAt = CurrentTimeMillis()
	return tx
}

func BalanceDeposit(userId string, quantity int64) {
	tx := BaseTransaction(userId)
	tx.TransactionType = Model.BALANCE
	tx.SubTransactionType = Model.BALANCE_DEPOSIT
	tx.Detail = strconv.FormatInt(quantity, 10)
}

func RegisterGame(userId string, gameId string) Model.TransactionLog {
	tx := BaseTransaction(userId)
	tx.TransactionType = Model.GAME
	tx.SubTransactionType = Model.GAME_REGISTER
	tx.Detail = gameId
	return tx
}

func WithdrawMoney(userId string, iban string, amount string) Model.TransactionLog {
	tx := BaseTransaction(userId)
	tx.TransactionType = Model.BALANCE
	tx.SubTransactionType = Model.BALANCE_WITHDRAW
	tx.Detail = iban + " numaralı hesaba " + amount + " TL aktarım talebi alınmıştır"
	return tx
}

func DepositMoney(userId string, amount string, description string) Model.TransactionLog {
	tx := BaseTransaction(userId)
	tx.TransactionType = Model.BALANCE
	tx.SubTransactionType = Model.BALANCE_WITHDRAW
	tx.Detail = userId + " numaralı hesaba " + amount + " TL aktarım talebi alınmıştır"
	return tx
}

func AddBalance(operator string, userId string, balance float64, bonus float64) Model.TransactionLog {
	tx := BaseTransactionWithOperator(userId, operator)
	tx.TransactionType = Model.BALANCE
	tx.SubTransactionType = Model.BALANCE_WITHDRAW
	if balance > 0 && bonus > 0 {
		tx.Detail = "Hesabiniza " + fmt.Sprintf("%.2f", balance) + " TL bakiye ve " + fmt.Sprintf("%.2f", bonus) + " TL bonus yuklenmistir"
	} else if balance > 0 {
		tx.Detail = "Hesabiniza " + fmt.Sprintf("%.2f", balance) + " TL bakiye yuklenmistir"
	} else if bonus > 0 {
		tx.Detail = "Hesabiniza " + fmt.Sprintf("%.2f", bonus) + " TL bonus yuklenmistir"
	} else {
		panic(errors.New("balance and bonus are both invalid"))
	}

	return tx
}

func UnregisterGame(userId string, gameId string) Model.TransactionLog {
	tx := BaseTransaction(userId)
	tx.TransactionType = Model.GAME
	tx.SubTransactionType = Model.GAME_UNREGISTER
	tx.Detail = gameId
	return tx
}

func GetTransactionType(transactionType Model.TransactionType) string {
	return Model.TRANSACTION_MAP[transactionType]
}

func GetSubTransactionType(subTransactionType Model.SubTransactionType) string {
	return Model.SUB_TRANSACTION_MAP[subTransactionType]
}

func WinGame(game Model.Game, operator string, userId string, award float64) Model.TransactionLog {
	tx := BaseTransactionWithOperator(userId, operator)
	tx.TransactionType = Model.GAME
	tx.SubTransactionType = Model.GAME_WIN
	tx.Detail = "Oyunu kazandiniz!\n" + game.Detail() + "\n" + fmt.Sprintf("%.2f", award) + " TL odul bakiyenize yuklendi."
	return tx
}
