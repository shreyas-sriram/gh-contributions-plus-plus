package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shreyas-sriram/gh-contributions-plus-plus/pkg/data"
	"github.com/shreyas-sriram/gh-contributions-plus-plus/pkg/draw"
	http_err "github.com/shreyas-sriram/gh-contributions-plus-plus/pkg/http-err"
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
		http_err.NewError(c, http.StatusBadRequest, fmt.Errorf("usage: <IP>/aggregate?username=<username1>&username=<username2>&year=<year>&theme=<light/dark>"))
		return
	}

	if len(year) == 0 {
		year = append(year, strconv.Itoa(time.Now().Year())) // Set default as current year
	}

	if len(theme) == 0 {
		theme = append(theme, "light") // Set default as "light"
	} else if !(theme[0] == "light" || theme[0] == "dark") {
		http_err.NewError(c, http.StatusBadRequest, fmt.Errorf("theme must be \"light\" or \"dark\""))
		return
	}

	request := new(data.Request)
	request.Usernames = usernames
	request.Year = year[0]
	request.Theme = theme[0]

	for _, username := range request.Usernames {
		rawHTML, err := data.GetRawPage(username, request.Year)
		if err != nil {
			http_err.NewError(c, http.StatusNotFound, fmt.Errorf("user not found"))
			return
		}

		contributions := data.ParseContributionsData(rawHTML, request.Year)

		request.ContributionList = append(request.ContributionList, contributions)
	}

	img, err := draw.ConstructMap(*request)
	if err != nil {
		http_err.NewError(c, http.StatusInternalServerError, fmt.Errorf("error creating chart"))
		return
	}

	imgHTML := "<html><body><img src=\"data:image/png;base64," + img + "\" /></body></html>"

	c.Data(http.StatusOK, "text/html", []byte(imgHTML))
}
