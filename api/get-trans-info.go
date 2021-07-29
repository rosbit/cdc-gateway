// 账务查询
// 参考文档: https://openbiz.cmbchina.com/developer/UI/Business/CloudDirectConnect/Public/DocumentCenter/DocDetailHasTitile.aspx?child=DCCT20201214145038074&firstLevel=1&fabizkey=1

package api

import (
	"cdc-gateway/conf"
	"os"
	"fmt"
)

const (
	debtor_flag   = "D"
	creditor_flag = "C"

	trans_info_api = "DCTRSINF"
)

// 账户信息
type AccountInfo struct {
	ContinueFlag string `json:"cotflg"`   // 未传完标记, 如为Y则表示查询交易日在本次查询完成后，还有数据未查，为N则表示已查完交易日的所有交易
	TransactionSeq uint64 `json:"trsseq,string"` // 末位记账序号, 本次查询所查的最大记账序号，如cotflg为Y，下一次查询输入接口记账序号字段填入此序号+1
	DebtorCount uint32 `json:"dbtnbr,string"`  // 借方笔数，本次查询查到的借记的笔数
	DebtorAmount float64 `json:"dbtamt,string"`  // 借方金额，本次查询查到的借记的金额汇总
	CreditorCount uint32 `json:"crtnbr,string"` // 贷方笔数，本次查询查到的贷记的笔数
	CreditorAmount float64 `json:"crtamt,string"` // 贷方金额，本次查询查到的贷记的金额汇总
}

type TransInfo struct {
	TransDate string `json:"etydat"`  // 交易日, 交易发生的日期 YYYYMMDD
	TransTime string `json:"etytim"`  // 交易时间, 交易发生的时间，只有小时有效 hhmmss
	ValueDate string `json:"vltdat"`  // 起息日, 开始计息的日期, YYYYMMDD
	TransCode string `json:"trscod"`  // 交易类型
	Abstract  string `json:"naryur"`  // 摘要, 若为企业银行客户端经办的交易，则该字段为用途信息（4.0版代发代扣业务除外）若为其它渠道经办的交易，则该字段为交易的简单说明和注解。该字段显示在电子回单摘要上。
	TransAmount float64 `json:"trsamt,string"` // 交易金额, 当为正数时：企业为借方 当为负数时：企业为贷方
	CDMark    string `json:"amtcdr"`  // 借贷标记，C:贷；D:借
	Balance   float64 `json:"trsblv,string"`  // 余额，帐户的联机余额
	RefNumber string `json:"refnbr"`  // 流水号, 银行会计系统交易流水号,可以和回单命名中的流水号关联
	ReqNumber string `json:"reqnbr"`  // 企业银行交易序号，唯一标示企业银行客户端发起的一笔交易
	BusinessName string `json:"busnam"` // 业务名称
	Usage     string `json:"nusage"`  // 用途
	BusinessRef string `json:"yurref"` // 业务参考号, 企业银行客户端录入的业务参考号。用企业银行做的交易会有业务参考号，没有票据号，在柜台或其它地方生成的交易有票据号或其它的唯一标识，都统一称为业务参考号
	BusinessAbstract string `json:"busnar"` // 业务摘要，对业务的简单说明或注解。企业银行客户端录入的摘要信息
	OtherAbstract string `json:"otrnar"` // 其它摘要，对业务的其它说明或注解（暂不使用）
	BranchBankNumber string `json:"rpybbk"`  // 收/付方开户地区分行号, 附录A.3，收/付方帐号开户行所在地区，如北京、上海、深圳等
	Name string `json:"rpynam"`    // 收/付方名称, 收/付方帐户名称
	Account string `json:"rpyacc"` // 收/付方帐号, 收/付方的转入或转出帐号
	BankNumber string `json:"rpybbn"` // 收/付方开户行行号, 收/付方帐号的开户行的行号
	BankName string `json:"rpybnk"`   // 收/付方开户行名, 收/付方帐号的开户行的行名
	BankAddress string `json:"rpyadr"`   // 收/付方开户行地址
	ParentOrSubBranchBankNumber string `json:"gsbbbk"` // 母/子公司所在地区分行, 见附录A.3，母/子公司帐号的开户行所在地区，如北京、上海、深圳等
	ParentOrSubAccount string `json:"gsbacc"` // 母/子公司帐号
	ParentOrSubName string `json:"gsbnam"` // 母/子公司名称
	InfoFlag string `json:"infflg"`  // 信息标志, 用于标识收/付方帐号和母/子公司的信息。为空表示付方帐号和子公司；为“1”表示收方帐号和子公司；为“2”表示收方帐号和母公司；为“3”表示原收方帐号和子公司
	AttachFlag string `json:"athflg"` // 有否附件信息标志, Y：是 N：否
	CheckNumber string `json:"chknbr"` // 票据号
	ReverseFlag string `json:"rsvflg"` // 冲帐标志, *为冲帐，X为补帐 （冲账交易与原交易借贷相反）
	AbstractExtended string `json:"narext"` // 扩展摘要, 有效位数为16
	TransAnalysisCode string `json:"trsanl"` // 交易分析码, 1-2位取值含义见附录A.4，3-6位取值含义见trscod字段说明。
	RefBusinessNumber string `json:"refsub"` // 商务支付订单号, 由商务支付订单产生
	FirmCode string `json:"frmcod"` // 企业识别码, 开通收方识别功能的账户可以通过此码识别付款方
	LastSeq uint64 `json:"-"` // 一个计数，借用结构体存放
}

type FnGetTransInfo func(app string, date string, seq uint64) (count uint32, it <-chan *TransInfo, err error)

// 借方（收入)
func GetDebtorTransInfo(app string, date string, seq uint64) (count uint32, it <-chan *TransInfo, err error) {
	return getTransInfo(app, date, debtor_flag, seq)
}

// 贷方（支出)
func GetCreditorTransInfo(app string, date string, seq uint64) (count uint32, it <-chan *TransInfo, err error) {
	return getTransInfo(app, date, creditor_flag, seq)
}

func getTransInfo(app string, date string, cdMark string, seq uint64) (count uint32, it <-chan *TransInfo, err error) {
	cdcConf := gwconf.GetCDCConf(app)
	if cdcConf == nil {
		err = fmt.Errorf("app name %s not found", app)
		return
	}

	// get first page
	header, trans, e := get1PageTransInfo(cdcConf, date, seq)
	if e != nil {
		err = e
		return
	}
	if len(trans) == 0 {
		return
	}
	setSeq(trans, header.TransactionSeq)

	switch cdMark {
	case debtor_flag:
		if header.DebtorCount == 0 {
			return
		}
		count = header.DebtorCount
	case creditor_flag:
		if header.CreditorCount == 0 {
			return
		}
		count = header.CreditorCount
	}

	transPage := make(chan []TransInfo)
	it = makeChanTrans(transPage, cdMark)

	// fetch other pages
	go func() {
		for {
			transPage <- trans
			if header.ContinueFlag == "N" {
				break
			}
			seq := header.TransactionSeq + 1
			if header, trans, e = get1PageTransInfo(cdcConf, date, seq); e != nil {
				// ignore error
				fmt.Fprintf(os.Stderr, "err occurs when calling get1PageTransInfo(seq: %d): %v\n", seq, e)
				break
			}
			if len(trans) == 0 {
				break
			}
			setSeq(trans, header.TransactionSeq)
		}
		close(transPage)
	}()

	return
}

func setSeq(trans []TransInfo, lastSeq uint64) {
	for i, _ := range trans {
		t := &trans[i]
		t.LastSeq = lastSeq
	}
}

func get1PageTransInfo(cdcConf *gwconf.CDCConfT, date string, transSeq uint64) (header *AccountInfo, trans []TransInfo, err error) {
	reqParams := map[string]interface{}{
		"bbknbr": cdcConf.Account.BankNo,
		"accnbr": cdcConf.Account.No,
		"trsdat": date,
		"trsseq": fmt.Sprintf("%d", transSeq),
	}
	var res struct {
		Response struct {
			Head cdcRespHead `json:"head"`
			Body struct {
				InfoHead  []AccountInfo `json:"ntrbptrsz1"`
				InfoDetails []TransInfo `json:"ntqactrsz2"`
			} `json:"body"`
		} `json:"response"`
	}
	if err = callCDC(cdcConf, trans_info_api, reqParams, &res); err != nil {
		return
	}

	head := &res.Response.Head
	if head.ResultCode != resp_ok {
		err = fmt.Errorf("respId: %s, resCode: %s, resMsg: %s", head.RespId, head.ResultCode, head.ResultMsg)
		return
	}
	body := &res.Response.Body
	if len(body.InfoHead) == 0 {
		err = fmt.Errorf("ntrbptrsz1 expected in response")
		return
	}
	header = &body.InfoHead[0]
	trans = body.InfoDetails

	return
}

func makeChanTrans(transPage <-chan []TransInfo, cdMark string) (<-chan *TransInfo) {
	it := make(chan *TransInfo)
	go func() {
		for trans := range transPage {
			for i, _ := range trans {
				t := &trans[i]
				if t.CDMark == cdMark {
					it <- t
				}
			}
		}

		close(it)
	}()
	return it
}

