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
		http_err.NewError(c, http.StatusNotFound, fmt.Errorf("no usernames given"))
		return
	}

	if len(year) == 0 {
		year[0] = strconv.Itoa(time.Now().Year()) // Set default as current year
	}

	if len(theme) == 0 {
		theme[0] = "light" // Set default as "light"
	} else if !(theme[0] == "light" || theme[0] == "dark") {
		http_err.NewError(c, http.StatusNotFound, fmt.Errorf("theme must be \"light\" or \"dark\""))
		return
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

		contributions, err := utils.ParseContributionsData(rawHTML, request.Year)
		if err != nil {
			http_err.NewError(c, http.StatusNotFound, fmt.Errorf("unable to find contributions, try again later"))
			return
		}

		request.ContributionList = append(request.ContributionList, contributions)
	}

	img, err := utils.ConstructMap(request)
	if err != nil {
		http_err.NewError(c, http.StatusNotFound, fmt.Errorf("error creating chart"))
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(img))
}
