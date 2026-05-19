package main

import (
	AwsUtils "../../internal"
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
		return AwsUtils.Response{StatusCode: 200, Body: getUserDetails(email)}, nil
	default:
		return AwsUtils.Response{StatusCode: 405, ErrorMessage: "unsupported operation: " + httpMethod}, nil
	}
}

func getUserDetails(email string) Model.User {
	filt := expression.Name("email").Equal(expression.Value(email))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		fmt.Println(err)
	}
	user, err := AwsUtils.GetUserDetails(expr)
	if err != nil {
		fmt.Println("Got error unmarshalling:")
		fmt.Println(err.Error())
	}
	return user
}

func main() {
	//getUserDetails("odalabasmaz+pg1@gmail.com")
	lambda.Start(Handler)
}
