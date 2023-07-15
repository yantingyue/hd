package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/valyala/fasthttp"
)

type ReplaceResp struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

func main() {
	go Token1()
	go Token2()
	//go Token3()
	select {}
}

var TT = time.NewTicker(time.Millisecond * 100)

//var TT = time.NewTicker(time.Millisecond * 500)

func Token1() {
	go func() {
		defer panicRecover()
		defer TT.Stop()
		for {
			<-TT.C
			go func() {
				defer panicRecover()
				fmt.Println(1111)
				if Replace(412, 40878538, "d76fa8fd19d84d769bce85e06390ac40") {
					time.Sleep(time.Second * 2)
				}
				fmt.Println(44444)
				if Replace(412, 40878970, "d76fa8fd19d84d769bce85e06390ac40") {
					time.Sleep(time.Second * 2)
				}
			}()
		}
	}()
}
func Token2() {
	go func() {
		defer panicRecover()
		defer TT.Stop()

		for {
			<-TT.C
			go func() {
				defer panicRecover()
				fmt.Println(2222)
				_ = Replace(412, 40878144, "7b7bebc482b34c429026d6e7b09c08f8")
			}()
		}
	}()
}
func Token3() {
	go func() {
		defer panicRecover()
		defer TT.Stop()
		for {
			<-TT.C
			go func() {
				defer panicRecover()
				fmt.Println(33333)
				_ = Replace(412, 40878614, "035c5639f15c44f9b838d38e242399ed")
			}()
		}
	}()
}

var panicRecover = func() {
	if err := recover(); err != nil {
		log.Println("panic err", err)
	}
}

func Replace(id, orderId uint64, token string) bool {
	header := GenerateHeader1(token)
	body := map[string]interface{}{
		"order_id":   orderId,
		"replace_id": id,
	}
	jsonBytes, _ := json.Marshal(body)
	resp, _ := Post("https://api.aichaoliuapp.cn/aiera/ai_match_trading/nft/replace/active/exchange", header, jsonBytes)
	log.Println(orderId, string(resp))
	//fmt.Println(orderId, string(resp))
	if len(resp) == 0 {
		return false
	}
	res := ReplaceResp{}
	json.Unmarshal(resp, &res)
	if res.Code == 0 && res.Msg == "success" {
		return true
	}
	return false
}

const (
	timeOut  = 60 * time.Second
	version  = "31850"
	channel  = "010100"
	platform = "ios"
	appname  = "aiera.sneaker.snkrs.shoe"
	salt     = "5c33494d1b277902d1b78f98093f6fd4"
)

func GenerateHeader1(token string) map[string]string {
	timestamp := cast.ToString(time.Now().UnixMilli())
	return map[string]string{
		"token":     token,
		"version":   version,
		"channel":   channel,
		"platform":  platform,
		"appname":   appname,
		"timestamp": timestamp,
		"sign":      MD5(timestamp + salt),
	}
}

func Get(host string) (m map[string]interface{}, err error) {
	statusCode, body, err := fasthttp.GetTimeout(nil, host, timeOut)
	if err != nil {
		return
	}
	if statusCode != fasthttp.StatusOK {
		err = errors.New(fmt.Sprintf("request failed statusCode[%d]", statusCode))
		return
	}
	if body == nil {
		err = errors.New("response body is nil")
		return
	}
	m = make(map[string]interface{})
	if err = json.Unmarshal(body, &m); err != nil {
		return
	}
	return
}

func Post(host string, header map[string]string, payload []byte) (body []byte, err error) {
	req := &fasthttp.Request{}
	req.SetRequestURI(host)

	req.Header.SetContentType("application/json")
	req.Header.SetMethod("POST")

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	req.SetBody(payload)

	resp := &fasthttp.Response{}
	client := &fasthttp.Client{}
	if err = client.DoTimeout(req, resp, timeOut); err != nil {
		return
	}
	body = resp.Body()
	return
}

// 生成32位MD5
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}
