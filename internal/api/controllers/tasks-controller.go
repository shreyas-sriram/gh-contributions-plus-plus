package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTaskById godoc
// @Summary Prepares contributions chart based on given usernames
// @Description get Task by ID
// @Produce json
// @Param username path array true "GitHub Usernames"
// @Success 200 {object} tasks.Task
// @Router /api/contributions [get]
func GetTaskById(c *gin.Context) {
	// id := c.Param("id")
	// if task, err := s.Get(id); err != nil {
	// 	http_err.NewError(c, http.StatusNotFound, errors.New("task not found"))
	// 	log.Println(err)
	// } else {
	// 	c.JSON(http.StatusOK, task)
	// }
	c.JSON(http.StatusOK, "hello world")
}
