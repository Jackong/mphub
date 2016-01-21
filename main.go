package main

import (
	"os"

	"github.com/Jackong/mphub/route"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.Any("/wechat", route.ServeWechat)
		server := api.Group("/servers/:server")
		server.Use(route.ValidServer)
		{
			server.POST("", route.SetServer)
			server.GET("/menus", route.GetMenu)
			oauth := server.Group("/oauth")
			{
				oauth.GET("/url", route.GetAuthURL)
				oauth.GET("/callback", route.CallbackAuth)
			}
			user := server.Group("/users")
			user.Use(route.ValidJWT)
			{
				user.GET("", route.GetUserInfo)
				user.POST("/bind", route.BindAOPS)
			}
		}
		msg := api.Group("/messages")
		{
			msg.POST("/template", route.PushTemplate)
		}
	}
	r.Run(os.Getenv("HTTP_ADDR"))
}
