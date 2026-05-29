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
			return svc.Response{StatusCode: 400, ErrorMessage: "The balance amount to be deposited [" + balance + "] is invalid!"}, nil
		}
		bonusFloat, err := strconv.ParseFloat(bonus, 64)
		if err != nil {
			return svc.Response{StatusCode: 400, ErrorMessage: "The bonus amount to be deposited [" + bonus + "] is invalid!"}, nil
		}

		user, err := addBalanceToUser(operator, userId, balanceFloat, bonusFloat)
		if err != nil {
			log.Printf("addBalanceToUser error: %v", err)
			return svc.Response{StatusCode: 500, ErrorMessage: "Failed to update user balance"}, nil
		}
		return svc.Response{StatusCode: 200, Body: user}, nil
	default:
		return svc.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + event.Context["http-method"]}, nil
	}
}

func addBalanceToUser(operator string, userId string, balance float64, bonus float64) (model.User, error) {
	user := svc.GetUserById(userId)
	user.Balance += balance
	user.Bonus += bonus

	tx := svc.AddBalance(operator, user.Id, balance, bonus)
	if err := svc.UpdateUserWithTx(user, tx); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func main() {
	lambda.Start(Handler)
}
