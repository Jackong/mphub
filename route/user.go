package route

import (
	"github.com/Sirupsen/logrus"
	"github.com/chanxuehong/wechat/mp/user"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	aopsID2User   = map[string][]*User{}
	openID2aopsID = map[string]string{}
)

//GetUserInfo get user info
/**
 * @api {get} /api/servers/:server/users 获取用户信息
 * @apiName GetUserInfo
 * @apiGroup user
 * @apiParam (Path) {String} server 平台服务名称
 * @apiSuccess {User} user 用户信息
 * @apiSuccess {Bind} bind 已绑定的信息
 * @apiSuccess {String} appID app ID
 */
func GetUserInfo(c *gin.Context) {
	server := c.Param(serverKey)
	if server == "" {
		fail(c, "Param server is required")
		return
	}
	token := c.MustGet("token").(*jwt.Token)
	openID := token.Claims["openID"].(string)
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
	aopsID := openID2aopsID[openID]
	bind := aopsID2User[aopsID]
	ok(c, map[string]interface{}{"user": user, "bind": bind, "appID": token.Claims["appID"]})
}

//BindAOPS bind wechat to AOPS
/**
 * @api {post} /api/servers/:server/users/bind 绑定到AOPS
 * @apiName BindAOPS
 * @apiGroup user
 * @apiParam (Path) {String} server 平台服务名称
 * @apiParam (Body) {String} aopsID AOPS ID
 * @apiSuccess {[]Bind} binds 已绑定的信息
 * @apiSuccess {bool} isNew 是否第一次绑定
 * @apiSuccess {bool} isUnion 是否unionID绑定
 */
func BindAOPS(c *gin.Context) {
	aopsID := c.PostForm("aopsID")
	if aopsID == "" {
		fail(c, "Param aopsID is required")
		return
	}
	token := c.MustGet("token").(*jwt.Token)
	appID := token.Claims["appID"].(string)
	openID := token.Claims["openID"].(string)
	unionID := token.Claims["unionID"].(string)
	users, exist := aopsID2User[aopsID]
	if !exist {
		users = append(users, &User{AppID: appID, OpenID: openID, UnionID: unionID})
		aopsID2User[aopsID] = users
		openID2aopsID[openID] = aopsID
		ok(c, map[string]interface{}{"binds": users, "isNew": true, "isUnion": false})
		return
	}
	for _, user := range users {
		if unionID != "" && unionID == user.UnionID {
			ok(c, map[string]interface{}{"binds": users, "isNew": false, "isUnion": true})
			return
		}
		if user.AppID == appID && user.OpenID == openID {
			ok(c, map[string]interface{}{"binds": users, "isNew": false, "isUnion": false})
			return
		}
		if user.OpenID != openID {
			fail(c, "The aopsID has been binded to same appID with other openID")
			return
		}
	}
}

//User user for wechat
type User struct {
	AppID   string
	OpenID  string
	UnionID string
}
