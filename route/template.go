package route

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/chanxuehong/wechat/mp/message/template"
	"github.com/gin-gonic/gin"
)

var (
	tpls = map[string]string{
		"test": "Gu4eTujTcSL9O_HwZxyFJ7UMuQ7Ge9nkTCl2miA3Vdw",
	}
)

//PushTemplate push template message
/**
 * @api {POST} /api/messages/template 推送模板消息
 * @apiName PushTemplate
 * @apiGroup message
 * @apiParam (Body) {String} aopsID AOPS ID
 */
func PushTemplate(c *gin.Context) {
	aopsID := c.PostForm("aopsID")
	if aopsID == "" {
		fail(c, "Param aopsID is required")
		return
	}
	users := aopsID2User[aopsID]
	if len(users) == 0 {
		fail(c, "Can not find any openID")
		return
	}
	appID2Server := map[string]string{}
	for server, config := range apps {
		appID2Server[config.AppId] = server
	}
	for _, user := range users {
		server := appID2Server[user.AppID]
		_, err := pushTpl(server, user.OpenID, json.RawMessage(`
      {
        "first": {
          "value": "推送测试",
          "color":"#173177"
        },
        "coupon": {
          "value": "1W元",
          "color":"#173177"
        },
        "expDate": {
          "value": "2016.01.01",
          "color":"#173177"
        },
        "remark": {
          "value": "非常感谢",
          "color":"#173177"
        }
      }
      `))
		if err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"server": server,
				"aopsID": aopsID,
				"openID": user.OpenID,
			}).Errorln("Failed to push template message")
			fail(c, "Failed to push template message")
			return
		}
	}
	ok(c, nil)
}

func pushTpl(server, openID string, data json.RawMessage) (int64, error) {
	client := template.NewClient(ts[server], nil)
	return client.Send(&template.TemplateMessage{
		ToUser:      openID,
		TemplateId:  tpls[server],
		URL:         "http://hcz.pingan.com",
		RawJSONData: data,
	})
}
