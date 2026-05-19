package main

import (
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
)

func Handler(event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	username := event.Request.UserAttributes["name"]
	log.Println("canUserRegister check for username:", username)

	if svc.UserExistsByUsername(username) {
		return event, nil
	}
	return event, errors.New("Girilen Kullanıcı Adı başka bir kullanıcı tarafından kullanılmaktadır!")
}

func main() {
	lambda.Start(Handler)
}
