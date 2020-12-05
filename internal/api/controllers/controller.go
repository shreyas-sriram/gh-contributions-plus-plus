package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shreyas-sriram/gh-contributions-aggregator/pkg/helpers"
	http_err "github.com/shreyas-sriram/gh-contributions-aggregator/pkg/http-err"
)

// GetContributionsChart godoc
// @Summary Prepares contributions chart based on given usernames
// @Description get contributions graph by usernames
// @Produce png
// @Param username path array true "GitHub Usernames"
// @Success 200 {object} tasks.Task
// @Router /api/contributions [get]
func GetContributionsChart(c *gin.Context) {
	// if task, err := s.Get(id); err != nil {
	// 	http_err.NewError(c, http.StatusNotFound, errors.New("task not found"))
	// 	log.Println(err)
	// } else {
	// 	c.JSON(http.StatusOK, task)
	// }
	queryParams := c.Request.URL.Query()
	usernames := queryParams["username"]

	var contributionList []helpers.Contributions

	for _, username := range usernames {
		rawHTML, err := helpers.GetRawPage(username)
		if err != nil {
			http_err.NewError(c, http.StatusNotFound, err)
			return
		}

		contributions, _ := helpers.ParseContributionsData(rawHTML)

		contributionList = append(contributionList, contributions)
	}

	aggregateContributions, _ := helpers.AggregateContributions(contributionList)
	b, err := json.Marshal(aggregateContributions)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(string(b))

	c.JSON(http.StatusOK, string(b))
}
