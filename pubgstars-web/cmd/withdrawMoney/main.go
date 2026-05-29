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
	iban := svc.CovertToString(event.Body["iban"])
	nameSurname := svc.CovertToString(event.Body["nameSurname"])
	secretQuestion := svc.CovertToString(event.Body["secretQuestion"])
	secretAnswer := svc.CovertToString(event.Body["secretAnswer"])

	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return svc.Response{StatusCode: 400, ErrorMessage: "The withdrawal amount is invalid!"}, nil
	}
	if amountFloat < 100 {
		return svc.Response{StatusCode: 400, ErrorMessage: "The withdrawal amount must be at least 100₺!"}, nil
	}

	user := svc.GetUserByEmail(email)
	if user.SecretQuestion != secretQuestion || user.SecretAnswer != secretAnswer {
		return svc.Response{StatusCode: 400, ErrorMessage: "Secret question or answer is incorrect."}, nil
	}
	if user.Balance < amountFloat {
		return svc.Response{StatusCode: 400, ErrorMessage: "Insufficient balance!"}, nil
	}

	user.Balance -= amountFloat
	tx := svc.WithdrawMoney(user.Id, iban, amount)

	requestText := fmt.Sprintf("Withdraw money request:\n\tUser: [%s]\n\tTo: [%s]\n\tIBAN: [%s]\n\tAmount: [%.2f TL]",
		email, nameSurname, iban, amountFloat)
	svc.SendMessage(requestText)

	if err := svc.UpdateUserWithTx(user, tx); err != nil {
		log.Printf("withdrawMoney transaction error: %v", err)
		svc.SendMessage("!!! " + requestText)
		return svc.Response{StatusCode: 500, ErrorMessage: "An unexpected error occurred!"}, nil
	}
	return svc.Response{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
