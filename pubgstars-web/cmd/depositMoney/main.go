package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
)

func Handler(ctx context.Context, event svc.RequestEvent) (svc.Response, error) {
	email := svc.GetUsernameFromJwtToken(event.Params["header"]["Authorization"])
	amount := svc.CovertToString(event.Body["amount"])
	description := svc.CovertToString(event.Body["description"])

	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return svc.Response{StatusCode: 400, ErrorMessage: "Yuklenmek istenen tutar geçersiz!"}, nil
	}

	user := svc.GetUserByEmail(email)
	tx := svc.DepositMoney(user.Id, amount, description)

	requestText := fmt.Sprintf("Deposit money request:\n\tUser: [%s]\n\tDescription: [%s]\n\tAmount: [%.2f TL]",
		email, description, amountFloat)
	svc.SendMessage(requestText)

	if err := svc.UpdateUserWithTx(user, tx); err != nil {
		log.Printf("depositMoney transaction error: %v", err)
		svc.SendMessage("!!! " + requestText)
		return svc.Response{StatusCode: 500, ErrorMessage: "Beklenmeyen bir hata oluştu!"}, nil
	}
	return svc.Response{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
