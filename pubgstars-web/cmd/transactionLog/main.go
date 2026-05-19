package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
)

func Handler(ctx context.Context, event svc.RequestEvent) (svc.Response, error) {
	email := svc.GetUsernameFromJwtToken(event.Params["header"]["Authorization"])
	log.Println("username found:", email)

	switch event.Context["http-method"] {
	case "GET":
		return svc.Response{StatusCode: 200, Body: getTransactionLogs(email)}, nil
	default:
		return svc.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + event.Context["http-method"]}, nil
	}
}

func getTransactionLogs(email string) []model.TransactionLogResponse {
	userId := svc.GetUserByEmail(email).Id
	filt := expression.Name("userId").Equal(expression.Value(userId))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	txLogs, err := svc.GetTransactionLogs(expr)
	if err != nil {
		log.Println("getTransactionLogs error:", err)
		return nil
	}

	var txLogsResponse []model.TransactionLogResponse
	for _, tx := range txLogs {
		txLogsResponse = append(txLogsResponse, model.TransactionLogResponse{
			UserName:           email,
			DateTime:           svc.ConvertMillisToString(tx.InsertedAt),
			TransactionType:    svc.GetTransactionType(tx.TransactionType),
			SubTransactionType: svc.GetSubTransactionType(tx.SubTransactionType),
			Detail:             tx.Detail,
		})
	}
	return txLogsResponse
}

func main() {
	lambda.Start(Handler)
}
