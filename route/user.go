package route

import (
	"github.com/Sirupsen/logrus"
	"github.com/chanxuehong/wechat/mp/user"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//GetUserInfo get user info
/**
 * @api {get} /api/servers/:server/users 获取用户信息
 * @apiName GetUserInfo
 * @apiGroup user
 * @apiParam (Path) {String} server 平台服务名称
 * @apiSuccess {User} user 用户信息
 */
func GetUserInfo(c *gin.Context) {
	server := c.Param(serverKey)
	if server == "" {
		fail(c, "Param server is required")
		return
	}
	cookie, err := c.Request.Cookie("token")
	if err != nil {
		fail(c, "Invalid token")
		return
	}
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		logrus.WithError(err).WithFields(logrus.Fields{
			"server": server,
			"token":  cookie,
		}).Errorln("Invalid token")
		fail(c, "Invalid token")
		return
	}

	client := user.NewClient(ts[server], nil)
	user, err := client.UserInfo(token.Claims["openID"].(string), "")
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"server": server,
			"token":  token,
		}).Errorln("Failed to get user info from wechat")
		fail(c, "Failed to get user info from wechat")
		return
	}
	ok(c, map[string]interface{}{"user": user})
}
