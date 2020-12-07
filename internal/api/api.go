package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shreyas-sriram/gh-contributions-aggregator/internal/api/router"
	"github.com/shreyas-sriram/gh-contributions-aggregator/internal/pkg/config"
)

func setConfiguration(configPath string) {
	config.Setup(configPath)
	gin.SetMode(config.GetConfig().Server.Mode)
}

func Run(configPath string) {
	if configPath == "" {
		configPath = "data/config.yml"
	}
	setConfiguration(configPath)
	conf := config.GetConfig()
	web := router.Setup()
	fmt.Println("GitHub Contributions Aggregator is running on port " + conf.Server.Port)
	fmt.Println("==================>")
	_ = web.Run(":" + conf.Server.Port)
}
