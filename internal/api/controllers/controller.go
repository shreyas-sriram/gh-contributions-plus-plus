package controllers

import (
	"fmt"
	"net/http"

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

	var contributionList []utils.Contributions

	for _, username := range usernames {
		rawHTML, err := utils.GetRawPage(username)
		if err != nil {
			http_err.NewError(c, http.StatusNotFound, err)
			return
		}

		contributions, _ := utils.ParseContributionsData(rawHTML)

		contributionList = append(contributionList, contributions)
	}

	err := utils.ConstructMap(contributionList)
	if err != nil {
		http_err.NewError(c, http.StatusNotFound, fmt.Errorf("Error creating chart"))
		return
	}

	c.JSON(http.StatusOK, "success")
}
