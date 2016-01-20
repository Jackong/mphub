package route

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

const (
	state = "test"
)

//GetAuthURL get auth url
/**
 * @api {get} /api/servers/:server/oauth/url?redirect=:redirect 获取授权登录url
 * @apiName GetAuthURL
 * @apiGroup oauth
 * @apiParam (Path) {String} server 平台服务名称
 * @apiParam (Query) {String} redirect 授权成功后跳转的web页面
 * @apiSuccess {String} url 授权登录url
 */
func GetAuthURL(c *gin.Context) {
	server := c.Param(serverKey)
	redirect := c.Query("redirect")
	if redirect == "" {
		fail(c, "Param redirect is required")
		return
	}
	app, exist := apps[server]
	if !exist || app == nil {
		logrus.WithFields(logrus.Fields{
			"server":   server,
			"app":      app,
			"redirect": redirect,
		}).Errorln("Server not found")
		fail(c, "Server not found")
		return
	}
	uri := c.Request.URL
	scheme := "http"
	if uri.Scheme != "" {
		scheme = uri.Scheme
	}
	host := c.Request.Host
	if uri.Host != "" {
		host = uri.Host
	}
	app.RedirectURI = fmt.Sprintf("%s://%s%s", scheme, host, strings.Replace(uri.Path, "/url", "/callback", 1))
	app.Scopes = []string{"snsapi_userinfo"}
	params := url.Values{}
	params.Set("redirect", redirect)
	ok(c, map[string]interface{}{"url": app.AuthCodeURL(state, params)})
}
