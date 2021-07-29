package req

import (
	"cdc-gateway/conf"
	"time"
	"fmt"
	"net/url"
)

func MakeReq(cdcConf *gwconf.CDCConfT, api string, reqBodyJSON map[string]interface{}) (reqBodyStr string, err error) {
	now := time.Now()
	reqId := fmt.Sprintf("%s%s", now.Format("20060102150405.000"), api)
	reqId = reqId[:14] + reqId[15:]

	uID := cdcConf.ChooseUID()
	body := map[string]interface{} {
		"request": map[string]interface{}{
			"head": map[string]interface{} {
				"funcode": api,
				"userid": uID,
				"reqid": reqId,
			},
			"body": reqBodyJSON,
		},
	}

	res, e := cdcConf.MakeSignature(body)
	if e != nil {
		err = e
		return
	}
	reqBodyStr = fmt.Sprintf("UID=%s&DATA=%s", uID, url.QueryEscape(res))
	return
}

func MakeRequest(app, api string, reqBodyJSON map[string]interface{}) (reqBodyStr string, err error) {
	cdcConf := gwconf.GetCDCConf(app)
	if cdcConf == nil {
		err = fmt.Errorf("unknown app name %s", app)
		return
	}
	return MakeReq(cdcConf, api, reqBodyJSON)
}
