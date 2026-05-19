package main

import (
	AwsUtils "../../internal"
	DataService "../../internal"
	TransactionLogUtils "../../internal"
	Model "../../model"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"strconv"
)

var (
	db = AwsUtils.GetDynamoDbClient("eu-central-1")
)

func Handler(ctx context.Context, event AwsUtils.RequestEvent) (AwsUtils.Response, error) {
	log.Println("begin !!")

	httpMethod := event.Context["http-method"]
	switch httpMethod {
	case "PUT", "POST":
		operator := AwsUtils.GetUsernameFromJwtTokenForAdmin(event.Params["header"]["Authorization"])
		inputMap := event.Body
		userId := AwsUtils.CovertToString(inputMap["userId"])
		balance := AwsUtils.CovertToString(inputMap["balance"])
		bonus := AwsUtils.CovertToString(inputMap["bonus"])

		balanceFloat, err := strconv.ParseFloat(balance, 64)
		if err != nil {
			return AwsUtils.Response{StatusCode: 400, ErrorMessage: "Yatirilmak istenen bakiye tutari [" + balance + "] geçersiz! "}, err
		}
		bonusFloat, err := strconv.ParseFloat(bonus, 64)
		if err != nil {
			return AwsUtils.Response{StatusCode: 400, ErrorMessage: "Yatirilmak istenen bonus tutari [" + bonus + "] geçersiz! " + bonus}, err
		}

		return AwsUtils.Response{StatusCode: 200, Body: addBalanceToUser(operator, userId, balanceFloat, bonusFloat)}, nil
	default:
		return AwsUtils.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + httpMethod}, nil
	}
}

func addBalanceToUser(operator string, userId string, balance float64, bonus float64) Model.User {
	user := AwsUtils.GetUserById(userId)

	user.Balance += balance
	user.Bonus += bonus

	tx := TransactionLogUtils.AddBalance(operator, user.Id, balance, bonus)
	err := DataService.UpdateUserWithTx(user, tx)
	if err != nil {
		log.Println("Error occurred.")
		log.Println(err)
	}

	return user
}

func main() {
	lambda.Start(Handler)
}
