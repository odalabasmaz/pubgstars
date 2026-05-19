#!/bin/bash -e

MODULE=$1
FUNCTION_NAME=${MODULE}-test
echo "Deploying ${MODULE} ..."

pushd .
echo "building.."
cd ../cmd/${MODULE}/
GOOS=linux go build -o main main.go
zip main.zip main

aws lambda get-function --function-name asd --profile pg --output json >res.out 2>&1
cat res.out | grep "Function not found:" | wc -n
rm -f res.out

echo "creating & uploading..."
aws lambda create-function \
    --function-name arn:aws:lambda:eu-central-1:470936150750:function:${FUNCTION_NAME} \
    --runtime go1.x \
    --role arn:aws:iam::470936150750:role/lambda-service-role \
    --handler main \
    --timeout 10 \
    --zip-file fileb://./main.zip \
    --profile pg

echo "deploying latest function.."
aws lambda update-function-code \
    --function-name arn:aws:lambda:eu-central-1:470936150750:function:${FUNCTION_NAME} \
    --zip-file fileb://./main.zip \
    --profile pg

echo "cleaning.."
rm -f main mai  n.zip
popd

echo "done."


# https://read.acloud.guru/serverless-golang-api-with-aws-lambda-34e442385a6a
# https://aws.amazon.com/getting-started/projects/build-serverless-web-app-lambda-apigateway-s3-dynamodb-cognito/module-4/
# https://docs.aws.amazon.com/apigateway/latest/developerguide/welcome.html
# https://hackernoon.com/error-handling-with-api-gateway-and-go-lambda-functions-fe0e10808732

aws lambda get-function --function-name