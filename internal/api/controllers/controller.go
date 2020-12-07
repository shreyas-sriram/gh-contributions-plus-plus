package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	http_err "github.com/shreyas-sriram/gh-contributions-aggregator/pkg/http-err"
	"github.com/shreyas-sriram/gh-contributions-aggregator/pkg/utils"
)

// GetContributionsChart godoc
// @Summary Prepares contributions chart based on given usernames
// @Description get contributions graph by usernames
// @Produce png
// @Param username path array true "GitHub Usernames"
// @Success 200 {object} tasks.Task
// @Router /api/contributions [get]
func GetContributionsChart(c *gin.Context) {
	usernames := c.Request.URL.Query()["username"]
	year := c.Request.URL.Query()["year"]
	theme := c.Request.URL.Query()["theme"]

	if len(usernames) == 0 {
		http_err.NewError(c, http.StatusNotFound, fmt.Errorf("No usernames given"))
		return
	}

	if len(year) == 0 {
		year[0] = strconv.Itoa(time.Now().Year()) // Set default as current year
	}

	if len(theme) == 0 {
		theme[0] = "classic" // Set default as "classic"
	}

	var request utils.Request
	request.Usernames = usernames
	request.Year = year[0]
	request.Theme = theme[0]

	for _, username := range request.Usernames {
		rawHTML, err := utils.GetRawPage(username, request.Year)
		if err != nil {
			http_err.NewError(c, http.StatusNotFound, err)
			return
		}

		contributions, _ := utils.ParseContributionsData(rawHTML, request.Year)

		request.ContributionList = append(request.ContributionList, contributions)
	}

	err := utils.ConstructMap(request)
	if err != nil {
		http_err.NewError(c, http.StatusNotFound, fmt.Errorf("Error creating chart"))
		return
	}

	c.JSON(http.StatusOK, "success")
}
