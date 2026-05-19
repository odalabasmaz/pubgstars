package main

import (
	DataService "../../internal"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

func Handler(event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	log.Println("")
	log.Printf("%v\n", event.Request.UserAttributes)
	log.Printf("PostConfirmation for userRegistered: %s with email: %s\n", event.UserName, event.Request.UserAttributes["email"])

	// added from request
	var username = event.Request.UserAttributes["name"]
	log.Println("userName:" + username)
	result := DataService.UserExistsByUsername(username)
	if result {
		log.Println("başarılı")
		return event, nil
	}
	log.Println("kullanıcı adı email kullanılıyor")
	return event, errors.New("Girilen Kullanıcı Adı / Email başka bir kullanıcı tarafından kullanılmaktadır!")
}

func main() {
	/*	result := DataService.UserExistsByUsernameOrEmail("yusuf", "email")
		if (result) {
			log.Println("başarılı");
			//return event, nil
		} else{
			log.Println("kullanıcı adı email kullanılıyor");

		}*/
	lambda.Start(Handler)
}
