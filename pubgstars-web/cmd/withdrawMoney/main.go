package main

import (
	AwsUtils "../../internal"
	DataService "../../internal"
	SlackService "../../internal"
	TransactionLogUtils "../../internal"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"strconv"
)

func Handler(ctx context.Context, event AwsUtils.RequestEvent) (AwsUtils.Response, error) {
	log.Println("begin !!")

	email := AwsUtils.GetUsernameFromJwtToken(event.Params["header"]["Authorization"])
	inputMap := event.Body
	amount := AwsUtils.CovertToString(inputMap["amount"])
	iban := AwsUtils.CovertToString(inputMap["iban"])
	nameSurname := AwsUtils.CovertToString(inputMap["nameSurname"])
	secretQuestion := AwsUtils.CovertToString(inputMap["secretQuestion"])
	secretAnswer := AwsUtils.CovertToString(inputMap["secretAnswer"])
	amountFloat, err := strconv.ParseFloat(amount, 64)

	if err != nil {
		return AwsUtils.Response{StatusCode: 400, ErrorMessage: "Çekilmek istenen tutar geçersiz!"}, err
	}

	if amountFloat < 100 {
		return AwsUtils.Response{StatusCode: 400, ErrorMessage: "Çekilmek istenen tutar en az 100₺ olmalidir!"}, err
	}

	user := AwsUtils.GetUserByEmail(email)
	user.Balance -= amountFloat
	if user.SecretQuestion != secretQuestion || user.SecretAnswer != secretAnswer {
		return AwsUtils.Response{StatusCode: 400, ErrorMessage: "Gizli soru veya cevabı yanlış."}, err
	}
	if user.Balance < 0 {
		return AwsUtils.Response{StatusCode: 400, ErrorMessage: "Yetersiz Bakiye!"}, err
	}

	tx := TransactionLogUtils.WithdrawMoney(user.Id, iban, amount)
	requestText := "Withdraw money request: \n" +
		"\tUser: [" + email + "]\n" +
		"\tTo: [" + nameSurname + "]\n" +
		"\tIBAN: [" + iban + "]\n" +
		"\tAmount: [" + fmt.Sprintf("%.2f", amountFloat) + " TL]"
	SlackService.SendMessage(requestText)
	log.Println("after: update > ", user)
	e := DataService.UpdateUserWithTx(user, tx)
	if e != nil {
		fmt.Println("Got error in transaction")
		fmt.Println(e.Error())
		SlackService.SendMessage("!!! " + requestText)
		return AwsUtils.Response{StatusCode: 500, ErrorMessage: "Beklenmeyen Hata Oluştu!"}, err
	}

	return AwsUtils.Response{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
