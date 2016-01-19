package main

import (
	"os"

	"github.com/Jackong/mphub/route"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
}

func main() {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.Any("/wechat", route.ServeWechat)
		api.POST("/servers/:server", route.SetServer)
		api.GET("/servers/:server/menus", route.GetMenu)
	}
	r.Run(os.Getenv("HTTP_ADDR"))
}
