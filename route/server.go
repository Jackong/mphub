package route

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/chanxuehong/wechat/mp"
	"github.com/chanxuehong/wechat/mp/menu"
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
	log.WithError(err).WithFields(log.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
		"query":  r.URL.Query(),
		"body":   r.PostForm,
	}).Errorln("Failed to handler wechat request")
}

//ServeWechat serve request from wechat
func ServeWechat(c *gin.Context) {
	ms.ServeHTTP(c.Writer, c.Request)
}

//SetServer set a server
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
		log.WithError(err).WithFields(log.Fields{
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

//GetMenu get menu from wechat server
func GetMenu(c *gin.Context) {
	server := c.Param(serverKey)
	client := menu.NewClient(ts[server], nil)
	menus, err := client.GetMenu()
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"server": server,
		}).Errorln("Failed to get menus from wechat")
		fail(c, "Failed to get menus from wechat")
		return
	}
	ok(c, map[string]interface{}{"menus": menus})
}
