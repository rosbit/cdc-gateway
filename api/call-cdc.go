package api

import (
	"github.com/rosbit/gnet"
	"cdc-gateway/conf"
	"cdc-gateway/req"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	resp_ok = "SUC0000"
)

type cdcRespHead struct {
	FunCode string `json:"funcode"` // 接口名称
	UserId  string `json:"userid"`  // 企业一网通用户ID
	ReqId   string `json:"reqid"`   // 请求唯一编号
	RespId  string `json:"rspid"`   // 响应ID
	ResultCode string `json:"resultcode"` // 错误代码, 正常响应为SUC0000，表明请求已经被银行成功接收，不代表业务处理成功，业务处理结果请调用各个业务查询接口获取。其他错误代码请参见resultmsg的错误描述。
	ResultMsg string `json:"resultmsg"`   // 错误描述
}

func callCDC(cdcConf *gwconf.CDCConfT, apiName string, params map[string]interface{}, res interface{}) (err error) {
	body, e := req.MakeReq(cdcConf, apiName, params)
	if e != nil {
		err = e
		return
	}

	status, content, _, e := gnet.Http(cdcConf.ServiceURL, gnet.M("POST"), gnet.Params(body))
	if e != nil {
		err = e
		return
	}
	if status != http.StatusOK {
		err = fmt.Errorf("status: %d, resp: %s", status, content)
		return
	}

	decryptedBody, _, e := cdcConf.ParseResponse(string(content))
	if e != nil {
		err = e
		return
	}

	err = json.Unmarshal(decryptedBody, res)
	return
}
