#!/bin/bash

#set -e

MODULE=$1
FUNCTION_NAME=${MODULE}
echo "## Deploying ${MODULE} ..."

pushd .
echo "## building.."
cd ../cmd/${MODULE}/ || exit 1
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
zip bootstrap.zip bootstrap

# create lambda if necessary
#echo "## creating lambda if necessary..."
#aws lambda create-function \
#    --function-name arn:aws:lambda:eu-central-1:470936150750:function:${FUNCTION_NAME} \
#    --runtime Go1.x \
#    --role arn:aws:iam::470936150750:role/lambda-service-role \
#    --handler main \
#    --region eu-central-1

echo "## deploying lambda function.."
aws lambda update-function-code \
    --function-name arn:aws:lambda:eu-central-1:470936150750:function:${FUNCTION_NAME} \
    --zip-file fileb://./bootstrap.zip \
    --profile pg

echo "## cleaning artifacts.."
rm -f bootstrap bootstrap.zip
popd || exit 1

echo "## done."

#aws lambda create-function \
# --region eu-central-1 \
# --function-name DiscoverMovies \
# --zip-file fileb://./deployment.zip \
# --runtime go1.x \
# --role arn:aws:iam::<account-id>:role/<role> \
# --handler main

# https://read.acloud.guru/serverless-golang-api-with-aws-lambda-34e442385a6a
# https://aws.amazon.com/getting-started/projects/build-serverless-web-app-lambda-apigateway-s3-dynamodb-cognito/module-4/
# https://docs.aws.amazon.com/apigateway/latest/developerguide/welcome.html
# https://hackernoon.com/error-handling-with-api-gateway-and-go-lambda-functions-fe0e10808732
