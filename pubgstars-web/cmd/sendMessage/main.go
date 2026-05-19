package main

import (
	AwsUtils "../../internal"
	DataService "../../internal"
	ModelUtils "../../internal"
	SlackService "../../internal"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"strconv"
)

func Handler(ctx context.Context, event AwsUtils.RequestEvent) (AwsUtils.Response, error) {
	log.Println("begin !!")
	//email := AwsUtils.GetUsernameFromJwtToken(event.Params["header"]["Authorization"])
	httpMethod := event.Context["http-method"]
	inputMap := event.Body
	from := AwsUtils.CovertToString(inputMap["from"])
	message := AwsUtils.CovertToString(inputMap["message"])
	switch httpMethod {
	case "POST":
		return sendMessage(from, message), nil
	default:
		return AwsUtils.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + httpMethod}, nil
	}
}

func sendMessage(from string, message string) AwsUtils.Response {
	isCustomer := false
	user := AwsUtils.GetUserByEmail(from)
	if user.Id != "" {
		isCustomer = true
	}
	requestText := "Customer message received: \n" +
		"\tIs Customer: [" + strconv.FormatBool(isCustomer) + "]\n" +
		"\tFrom: [" + from + "]\n" +
		"\tMessage: [" + message + "]"

	// send message to slack
	SlackService.SendMessage(requestText)

	// save to db
	msgMap := map[string]interface{}{
		"id":         ModelUtils.GenerateKey(10),
		"dateTime":   AwsUtils.CurrentTimeMillis(),
		"status":     "waiting",
		"isCustomer": isCustomer,
		"from":       from,
		"message":    message,
	}
	if err := DataService.SaveMessage(msgMap); err != nil {
		return AwsUtils.Response{StatusCode: 400, ErrorMessage: err.Error()}
	}
	return AwsUtils.Response{StatusCode: 200, Body: "ok"}
}

func main() {
	//sendMessage("odalabasmaz+pg1@gmail.com", "msg")
	lambda.Start(Handler)
}
