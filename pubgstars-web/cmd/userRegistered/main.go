package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
)

func Handler(event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	log.Printf("PostConfirmation for userRegistered: %s with email: %s\n",
		event.UserName, event.Request.UserAttributes["email"])

	now := svc.CurrentTimeMillis()
	user := model.User{
		Id:             svc.GenerateKey(10),
		Status:         "active",
		InsertedAt:     now,
		InsertedBy:     "system",
		UpdatedAt:      now,
		UpdatedBy:      "system",
		Username:       event.Request.UserAttributes["name"],
		Email:          event.Request.UserAttributes["email"],
		SecretQuestion: event.Request.UserAttributes["custom:secretQuestion"],
		SecretAnswer:   event.Request.UserAttributes["custom:secretAnswer"],
	}

	err := svc.SaveUser(user)
	return event, err
}

func main() {
	lambda.Start(Handler)
}
