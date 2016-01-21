package route

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/chanxuehong/wechat/oauth2"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	state  = "test"
	secret = "jackong"
)

var (
	tokens = map[string]*oauth2.Token{}
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

//CallbackAuth callback for oauth
func CallbackAuth(c *gin.Context) {
	server := c.Param(serverKey)
	redirect := c.Query("redirect")
	code := c.Query("code")
	if code == "" {
		fail(c, "Param code is required")
		return
	}
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
	cli := &oauth2.Client{Config: app}
	accessToken, err := cli.Exchange(code)
	if err != nil || accessToken == nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"server":   server,
			"app":      app,
			"redirect": redirect,
			"token":    accessToken,
		}).Errorln("Failed to exchange code and token")
		fail(c, "Failed to exchange code and token")
		return
	}
	expires := time.Now().Add(time.Hour * 24)
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["exp"] = expires.Unix()
	token.Claims["openID"] = accessToken.OpenId
	token.Claims["unionID"] = accessToken.UnionId
	token.Claims["appID"] = app.AppId
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"server":      server,
			"app":         app,
			"redirect":    redirect,
			"accessToken": accessToken,
			"token":       token,
		}).Errorln("Failed to sign a token")
		fail(c, "Failed to sign a token")
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{Name: "token", Value: tokenString, Path: "/", Expires: expires, HttpOnly: true})
	c.Redirect(http.StatusTemporaryRedirect, redirect)
}

//ValidJWT valid token for jwt
func ValidJWT(c *gin.Context) {
	cookie, err := c.Request.Cookie("token")
	if err != nil {
		c.Abort()
		fail(c, "Invalid token")
		return
	}
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	server := c.Param(serverKey)

	if err != nil || !token.Valid {
		logrus.WithError(err).WithFields(logrus.Fields{
			"server": server,
			"token":  cookie,
		}).Errorln("Invalid token")
		c.Abort()
		fail(c, "Invalid token")
		return
	}
	c.Set("token", token)
	c.Next()
}
