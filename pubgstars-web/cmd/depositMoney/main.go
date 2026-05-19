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
	description := AwsUtils.CovertToString(inputMap["description"])
	amountFloat, err := strconv.ParseFloat(amount, 64)

	if err != nil {
		return AwsUtils.Response{StatusCode: 400, ErrorMessage: "Yuklenmek istenen tutar geçersiz!"}, err
	}

	user := AwsUtils.GetUserByEmail(email)

	tx := TransactionLogUtils.DepositMoney(user.Id, amount, description)
	requestText := "Deposit money request: \n" +
		"\tUser: [" + email + "]\n" +
		"\tDescription: [" + description + "]\n" +
		"\tAmount: [" + fmt.Sprintf("%.2f", amountFloat) + " TL]"
	SlackService.SendMessage(requestText)
	e := DataService.UpdateUserWithTx(user, tx)
	if e != nil {
		fmt.Println("Got error in transaction")
		fmt.Println(e.Error())
		SlackService.SendMessage("!!! " + requestText)
		return AwsUtils.Response{StatusCode: 500, ErrorMessage: "Beklenmeyen bir hata oluştu!"}, err
	}

	return AwsUtils.Response{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
