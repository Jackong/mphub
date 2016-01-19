package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func send(c *gin.Context, code int64, data map[string]interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	data["code"] = code
	c.JSON(http.StatusOK, data)
}

func fail(c *gin.Context, err string) {
	send(c, codeFail, map[string]interface{}{"error": err})
}

func ok(c *gin.Context, data map[string]interface{}) {
	send(c, codeOK, data)
}
