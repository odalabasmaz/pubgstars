#!/bin/bash

# you need aws cli profile with name pg

echo "## building..."
npm run build

echo "## uploading to S3..."
aws s3 sync build s3://pubgstars.com/ --profile pg

echo "## cleaning artifacts..."
rm -rf ./build/

echo "## done."
