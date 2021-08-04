package ce

import (
	"github.com/rosbit/http-helper"
	"cdc-gateway/cache"
	"cdc-gateway/conf"
	"cdc-gateway/api"
	"net/http"
	"strings"
	"time"
	"fmt"
	"io"
	"bytes"
	"encoding/json"
)

// POST ${commonEndpoints.GetTransInfo}/:app/:cdMark
// {
//    "date": "YYYY-MM-DD",
//    "seq": 0
// }
func GetTransInfo(c *helper.Context) {
	app := c.Param("app")
	cdMark := strings.ToUpper(c.Param("cdMark"))

	var getTransInfo api.FnGetTransInfo
	switch cdMark {
	case "C":
		getTransInfo = api.GetCreditorTransInfo
	case "D":
		getTransInfo = api.GetDebtorTransInfo
	case "A", "":
		getTransInfo = api.GetTransInfo
	default:
		c.Error(http.StatusBadRequest, "unknown cdMark")
		return
	}

	var params struct {
		Date string `json:"date"`
		Seq  uint64 `json:"seq"`
	}
	if code, err := c.ReadJSON(&params); err != nil {
		c.Error(code, err.Error())
		return
	}
	if len(params.Date) == 0 {
		c.Error(http.StatusBadRequest, "date expected")
		return
	}

	d, err := time.ParseInLocation("2006-01-02", params.Date, gwconf.Loc)
	if err != nil {
		c.Error(http.StatusBadRequest, err.Error())
		return
	}
	date := d.Format("20060102")

	cachedContent := cache.Get(app, date, cdMark, params.Seq)
	fmt.Printf("get content from cache (%s, %s, %s, %d): %p\n", app, date, cdMark, params.Seq, cachedContent)
	if cachedContent != nil {
		c.Blob(http.StatusOK, "application/json", cachedContent)
		return
	}

	count, it, err := getTransInfo(app, date, params.Seq)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if count == 0 {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": http.StatusOK,
			"msg": "OK",
			"count": 0,
			"trans": nil,
		})
		return
	}

	w := &bytes.Buffer{}
	makeContent(w, count, it)
	cachedContent = w.Bytes()
	go func() {
		cache.Set(app, date, cdMark, params.Seq, cachedContent, gwconf.ServiceConf.CacheLiveTime)
	}()

	c.Blob(http.StatusOK, "application/json", cachedContent)
}

func makeContent(w io.Writer, count uint32, it <-chan *api.TransInfo) {
	j := json.NewEncoder(w)
	j.SetEscapeHTML(false)

	fmt.Fprintf(w, `{"code":%d,"msg":"OK","count":%d,"trans":[`, http.StatusOK, count)
	first := true
	lastSeq := uint64(0)
	for trans := range it {
		if first {
			first = false
		} else {
			fmt.Fprintf(w, ",")
		}
		j.Encode(trans)
		lastSeq = trans.LastSeq
	}
	fmt.Fprintf(w, `],"last-seq":%d}`, lastSeq)
}

