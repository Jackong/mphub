package route

import (
	"github.com/Sirupsen/logrus"
	"github.com/chanxuehong/wechat/mp/menu"
	"github.com/gin-gonic/gin"
)

//GetMenu get menu from wechat server
/**
 * @api {get} /api/servers/:server/menus 获取公众号菜单栏
 * @apiName GetMenu
 * @apiGroup menu
 * @apiParam (Path) {String} server 平台服务名称
 * @apiSuccess {Object} menus 菜单栏
 */
func GetMenu(c *gin.Context) {
	server := c.Param(serverKey)
	client := menu.NewClient(ts[server], nil)
	menus, err := client.GetMenu()
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"server": server,
		}).Errorln("Failed to get menus from wechat")
		fail(c, "Failed to get menus from wechat")
		return
	}
	ok(c, map[string]interface{}{"menus": menus})
}
