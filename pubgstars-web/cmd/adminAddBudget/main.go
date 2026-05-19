package main

import (
	"context"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
)

func Handler(ctx context.Context, event svc.RequestEvent) (svc.Response, error) {
	switch event.Context["http-method"] {
	case "PUT", "POST":
		operator := svc.GetUsernameFromJwtTokenForAdmin(event.Params["header"]["Authorization"])
		userId := svc.CovertToString(event.Body["userId"])
		balance := svc.CovertToString(event.Body["balance"])
		bonus := svc.CovertToString(event.Body["bonus"])

		balanceFloat, err := strconv.ParseFloat(balance, 64)
		if err != nil {
			return svc.Response{StatusCode: 400, ErrorMessage: "Yatirilmak istenen bakiye tutari [" + balance + "] geçersiz!"}, nil
		}
		bonusFloat, err := strconv.ParseFloat(bonus, 64)
		if err != nil {
			return svc.Response{StatusCode: 400, ErrorMessage: "Yatirilmak istenen bonus tutari [" + bonus + "] geçersiz!"}, nil
		}

		return svc.Response{StatusCode: 200, Body: addBalanceToUser(operator, userId, balanceFloat, bonusFloat)}, nil
	default:
		return svc.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + event.Context["http-method"]}, nil
	}
}

func addBalanceToUser(operator string, userId string, balance float64, bonus float64) model.User {
	user := svc.GetUserById(userId)
	user.Balance += balance
	user.Bonus += bonus

	tx := svc.AddBalance(operator, user.Id, balance, bonus)
	if err := svc.UpdateUserWithTx(user, tx); err != nil {
		log.Printf("addBalanceToUser error: %v", err)
	}
	return user
}

func main() {
	lambda.Start(Handler)
}
