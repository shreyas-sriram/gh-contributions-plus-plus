package router

import (
	"fmt"
	"io"
	"os"

	"github.com/shreyas-sriram/gh-contributions-aggregator/internal/api/controllers"
	"github.com/shreyas-sriram/gh-contributions-aggregator/internal/api/middlewares"

	"github.com/gin-gonic/gin"
)

// Setup function sets-up the router
func Setup() *gin.Engine {
	app := gin.New()

	// Logging to a file.
	f, _ := os.Create("log/api.log")
	gin.DisableConsoleColor()
	gin.DefaultWriter = io.MultiWriter(f)

	// Middlewares
	app.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - - [%s] \"%s %s %s %d %s \" \" %s\" \" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	app.Use(gin.Recovery())
	app.NoRoute(middlewares.NoRouteHandler())

	// Routes
	app.GET("/api/contributions", controllers.GetContributionsChart)

	return app
}
