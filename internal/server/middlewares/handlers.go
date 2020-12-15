package middlewares

import "github.com/gin-gonic/gin"

// NoMethodHandler function handles unknown methods
func NoMethodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(405, gin.H{"message": "Method not allowed"})
	}
}

// NoRouteHandler function handles unknown routes
func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Invalid URL"})
	}
}
