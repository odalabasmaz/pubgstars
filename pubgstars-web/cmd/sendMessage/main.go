package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
)

func Handler(ctx context.Context, event svc.RequestEvent) (svc.Response, error) {
	from := svc.CovertToString(event.Body["from"])
	message := svc.CovertToString(event.Body["message"])

	switch event.Context["http-method"] {
	case "POST":
		return sendMessage(from, message), nil
	default:
		return svc.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + event.Context["http-method"]}, nil
	}
}

func sendMessage(from string, message string) svc.Response {
	user := svc.GetUserByEmail(from)
	isCustomer := user.Id != ""

	requestText := fmt.Sprintf("Customer message received:\n\tIs Customer: [%s]\n\tFrom: [%s]\n\tMessage: [%s]",
		strconv.FormatBool(isCustomer), from, message)
	svc.SendMessage(requestText)

	msgMap := map[string]interface{}{
		"id":         svc.GenerateKey(10),
		"dateTime":   svc.CurrentTimeMillis(),
		"status":     "waiting",
		"isCustomer": isCustomer,
		"from":       from,
		"message":    message,
	}
	if err := svc.SaveMessage(msgMap); err != nil {
		return svc.Response{StatusCode: 400, ErrorMessage: err.Error()}
	}
	return svc.Response{StatusCode: 200, Body: "ok"}
}

func main() {
	lambda.Start(Handler)
}
