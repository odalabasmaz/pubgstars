package main

import (
	AwsUtils "../../internal"
	TransactionLogUtils "../../internal"
	Model "../../model"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"log"
)

func Handler(ctx context.Context, event AwsUtils.RequestEvent) (AwsUtils.Response, error) {
	log.Println("begin !!")

	email := AwsUtils.GetUsernameFromJwtToken(event.Params["header"]["Authorization"])
	log.Println("username found: " + email)

	httpMethod := event.Context["http-method"]
	switch httpMethod {
	case "GET":
		return AwsUtils.Response{StatusCode: 200, Body: getTransactionLogs(email)}, nil
	default:
		return AwsUtils.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + httpMethod}, nil
	}
}

func getTransactionLogs(email string) []Model.TransactionLogResponse {
	userId := AwsUtils.GetUserByEmail(email).Id
	filt := expression.Name("userId").Equal(expression.Value(userId))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	txLogs, err := AwsUtils.GetTransactionLogs(expr)
	var txLogsResponse []Model.TransactionLogResponse

	for _, tx := range txLogs {
		var txRes Model.TransactionLogResponse
		txRes.UserName = email //TODO: caution
		txRes.DateTime = AwsUtils.ConvertMillisToString(tx.InsertedAt)
		txRes.TransactionType = TransactionLogUtils.GetTransactionType(tx.TransactionType)
		txRes.SubTransactionType = TransactionLogUtils.GetSubTransactionType(tx.SubTransactionType)
		txRes.Detail = tx.Detail
		txLogsResponse = append(txLogsResponse, txRes)
	}

	return txLogsResponse
}

func main() {
	//getTransactionLogs("odalabasmaz+pg1@gmail.com")
	lambda.Start(Handler)
	//test()
}

//todo remove after local test
func test() (string, error) {
	username := "orhun"
	logs := getTransactionLogs(username)
	for _, l := range logs {
		println(l.UserName, l.TransactionType, l.SubTransactionType, l.Detail, l.DateTime)
	}
	return "ok", nil
}
