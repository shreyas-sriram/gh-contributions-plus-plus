package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/shreyas-sriram/gh-contributions-plus-plus/pkg/data"
	"github.com/shreyas-sriram/gh-contributions-plus-plus/pkg/draw"
)

// Response is a wrapper for events.APIGatewayProxyResponse
type Response events.APIGatewayProxyResponse

// GetContributionsChartLambda godoc
// @Summary Prepares contributions chart based on given usernames
// @Description get contributions graph by usernames
// @Produce png
// @Param username path array true "GitHub Usernames"
// @Success 200 {object} tasks.Task
// @Router /api/contributions [get]
func GetContributionsChartLambda(gatewayRequest events.APIGatewayProxyRequest) (Response, error) {
	usernames := gatewayRequest.MultiValueQueryStringParameters["username"]
	year := gatewayRequest.QueryStringParameters["year"]
	theme := gatewayRequest.QueryStringParameters["theme"]

	if len(usernames) == 0 {
		return Response{Body: "usage: <IP>/aggregate?username=<username1>&username=<username2>&year=<year>&theme=<light/dark>", StatusCode: http.StatusOK}, nil
	}

	if len(year) == 0 {
		year = strconv.Itoa(time.Now().Year()) // Set default as current year
	}

	if len(theme) == 0 {
		theme = "light" // Set default as "light"
	} else if !(theme == "light" || theme == "dark") {
		return Response{Body: "theme must be \"light\" or \"dark\"", StatusCode: http.StatusBadRequest}, nil
	}

	request := new(data.Request)
	request.Usernames = usernames
	request.Year = year
	request.Theme = theme

	for _, username := range request.Usernames {
		rawHTML, err := data.GetRawPage(username, request.Year)
		if err != nil {
			return Response{Body: "user not found", StatusCode: http.StatusNotFound}, nil
		}

		contributions := data.ParseContributionsData(rawHTML, request.Year)

		request.ContributionList = append(request.ContributionList, contributions)
	}

	img, err := draw.ConstructMap(*request)
	if err != nil {
		return Response{Body: "error creating chart", StatusCode: http.StatusInternalServerError}, nil
	}

	return Response{Body: img, Headers: map[string]string{"Content-Type": "image/png"}, StatusCode: http.StatusOK, IsBase64Encoded: true}, nil
}
