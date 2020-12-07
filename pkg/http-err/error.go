package error

import "github.com/gin-gonic/gin"

// HTTPError encapsulates a custom http error
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// NewError function creates an object of "HTTPError"
func NewError(c *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	c.JSON(status, er)
}
