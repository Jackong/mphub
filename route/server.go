package route

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/chanxuehong/wechat/mp"
	"github.com/gin-gonic/gin"
)

var (
	msm *mp.MessageServeMux
	//Servers multiple servers
	ms *mp.MultiServerFrontend

	ts map[string]mp.AccessTokenServer
)

const serverKey = "server"

func init() {
	msm = mp.NewMessageServeMux()
	ms = mp.NewMultiServerFrontend(serverKey, mp.ErrorHandlerFunc(errHandler), nil)
	ts = map[string]mp.AccessTokenServer{}
}

func errHandler(w http.ResponseWriter, r *http.Request, err error) {
	logrus.WithError(err).WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
		"query":  r.URL.Query(),
		"body":   r.PostForm,
	}).Errorln("Failed to handler wechat request")
}

//ServeWechat serve request from wechat
/**
 * @api {post} /api/wechat?server=:server 微信服务
 * @apiName ServeWechat
 * @apiGroup wechat
 * @apiParam (Query) {String} server 平台服务名称
 */
func ServeWechat(c *gin.Context) {
	ms.ServeHTTP(c.Writer, c.Request)
}

//SetServer set a server
/**
 * @api {post} /api/servers/:server 添加公众号服务
 * @apiName SetServer
 * @apiGroup wechat
 * @apiParam (Path) {String} server 平台服务名称
 * @apiParam (Body) {String} oriID 公众号原始ID
 * @apiParam (Body) {String} appID 公众号AppID
 * @apiParam (Body) {String} appSecret 公众号secrect
 * @apiParam (Body) {String} token 公众号token
 * @apiParam (Body) {String} aesKey 公众号AESKey
 */
func SetServer(c *gin.Context) {
	oriID := c.DefaultPostForm("oriID", "")
	appID := c.DefaultPostForm("appID", "")
	appSecret := c.DefaultPostForm("appSecret", "")
	token := c.DefaultPostForm("token", "test")
	aesKey := []byte(c.DefaultPostForm("aesKey", ""))
	if len(aesKey) == 0 {
		aesKey = nil
	}
	server := c.Param(serverKey)
	err := ms.SetServer(server, mp.NewDefaultServer(oriID, token, appID, aesKey, msm))
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"body":   c.Request.PostForm,
			"server": server,
		}).Errorln("Failed to set server for wechat")
		fail(c, "Failed to set server for wechat")
		return
	}
	if _, ok := ts[server]; ok {
		delete(ts, server)
	}
	ts[server] = mp.NewDefaultAccessTokenServer(appID, appSecret, nil)
	ok(c, nil)
}
