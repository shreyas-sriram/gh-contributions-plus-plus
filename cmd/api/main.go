package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/shreyas-sriram/gh-contributions-aggregator/internal/api"
	"github.com/shreyas-sriram/gh-contributions-aggregator/internal/api/controllers"
)

func main() {

	if env := os.Getenv("DEPLOY"); env == "server" {
		api.Run("")
	} else {
		lambda.Start(controllers.GetContributionsChartLambda)
	}
}
