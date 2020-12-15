package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/shreyas-sriram/gh-contributions-plus-plus/internal/serverless/handler"
)

func main() {

	lambda.Start(handler.GetContributionsChartLambda)
}
