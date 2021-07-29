package ce

import (
	"github.com/rosbit/http-helper"
	"cdc-gateway/conf"
	"net/http"
	"io"
)

// POST ${commonEndpoints.ParseResponse}/:app
// body为访问银企直连接口返回的内容
func ParseResponse(c *helper.Context) {
	app := c.Param("app")
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	cdcConf := gwconf.GetCDCConf(app)
	if cdcConf == nil {
		c.Error(http.StatusBadRequest, "unknown app name")
		return
	}

	_, res, err := cdcConf.ParseResponse(string(body))
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"resp": res,
	})
}

