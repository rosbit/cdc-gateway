package ce

import (
	"github.com/rosbit/http-helper"
	"cdc-gateway/req"
	"net/http"
)

// POST ${commonEndpoints.MakeRequest}/:app/:api
//  - :app 配置文件中的应用名称
//  - :api 银企直连的接口名称，8个字母
//  - POST body内容为接口请求报文中的"body"对应的值内容。
// 返回: 是一个可直接用于访问银企直连接口的请求体
func MakeRequest(c *helper.Context) {
	app := c.Param("app")
	api := c.Param("api")
	if len(api) != 8 {
		c.Error(http.StatusBadRequest, "bad api size")
		return
	}

	var body map[string]interface{}
	if code, err := c.ReadJSON(&body); err != nil {
		c.Error(code, err.Error())
		return
	}
	if body == nil || len(body) == 0 {
		c.Error(http.StatusBadRequest, "body expected")
		return
	}

	res, err := req.MakeRequest(app, api, body)
	if err != nil {
		c.Error(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"req": res,
	})
}

