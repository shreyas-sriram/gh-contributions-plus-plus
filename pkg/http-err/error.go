package error

import "github.com/gin-gonic/gin"

// NewError function creates an object of "HTTPError"
func NewError(c *gin.Context, status int, err error) {
	c.Data(status, "text/plain", []byte(err.Error()))
}
