package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/valyala/fasthttp"
)

type ReplaceResp struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		OnSaleStatus       uint32 `json:"on_sale_status"`
		CurrentMilliTime   int64  `json:"current_milli_time"`
		StartTimeTimestamp int64  `json:"start_time_timestamp"`
		EndTimestamp       int64  `json:"end_time_timestamp"`
	} `json:"data"`
}
type ReplaceTimeResp struct {
	Code             int32  `json:"code"`
	Msg              string `json:"msg"`
	CurrentMilliTime int64  `json:"current_milli_time"`
}

type ResponseData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		OrderID    int    `json:"order_id"`
		Title      string `json:"title"`
		Picture    string `json:"picture"`
		UpdateTime string `json:"update_time"`
		SubTitle   string `json:"sub_title"`
		Type       string `json:"type"`
		Weight     int    `json:"weight"`
	} `json:"data"`
}

func main() {

	go Fj()
	select {}
}

const (
	b                 = 2                                  //1是分解 2是置换
	actId             = 535                                //活动id
	thread            = 2                                  //并发数
	tokenCommon       = "3662106dd9b749d3995348c5f2884a5b" //勿删
	tokenYanTingYue   = "3662106dd9b749d3995348c5f2884a5b" //颜庭跃
	tokenYanTingYueDa = "24715fa709414f6eb364ffb6f8c13485" //颜庭跃

)

func Fj() {
	go func() {
		for {
			//分解
			if b == 1 {
				for i := 0; i < thread; i++ {
					go func() {
						if FjDetail(actId, tokenCommon) {
							//颜庭跃
							go func() {
								if len(tokenYanTingYue) > 0 {
									//查看订单详情
									orderId := GetOrderId(actId, tokenYanTingYue)
									if orderId > 0 {
										Replace(actId, orderId, tokenYanTingYue)
									}
								}
							}()
						}
					}()
					time.Sleep(time.Millisecond * 20)
				}
			}
			//置换
			if b == 2 {
				for i := 0; i < thread; i++ {
					go func() {
						if ReplaceDetail(actId, tokenCommon) {
							Replace(actId, 206110329, tokenYanTingYue)
						}
					}()
				}
				time.Sleep(time.Millisecond * 20)

			}
		}
	}()

}

func GetOrderId(id uint64, token string) uint64 {
	header := GenerateHeader1(token)
	body := map[string]interface{}{
		"replace_id": id,
		"pageNumber": 1,
		"pageSize":   20,
	}
	jsonBytes, _ := json.Marshal(body)
	resp, _ := Post("https://api.aichaoliuapp.cn/aiera/ai_match_trading/nft/combination/choice/material", header, jsonBytes)
	log.Println(string(resp))
	if len(resp) == 0 {
		return 0
	}
	res := ResponseData{}
	json.Unmarshal(resp, &res)
	if res.Code == 0 && res.Msg == "success" && len(res.Data) > 0 {
		return uint64(res.Data[0].OrderID)
	}
	return 0

}

func FjDetail(id uint64, token string) bool {
	header := GenerateHeader1(token)
	body := map[string]interface{}{
		"replace_id": id,
	}
	jsonBytes, _ := json.Marshal(body)
	resp, _ := Post("https://api.aichaoliuapp.cn/aiera/ai_match_trading/nft/replace/active/detail", header, jsonBytes)
	log.Println(string(resp))
	if len(resp) == 0 {
		return false
	}
	res := ReplaceResp{}
	json.Unmarshal(resp, &res)
	if res.Code == 0 && res.Msg == "success" && res.Data.OnSaleStatus == 1 {
		return true
	}
	return false
}
func ReplaceDetail(id uint64, token string) bool {
	header := GenerateHeader1(token)
	body := map[string]interface{}{
		"replace_id": id,
	}
	jsonBytes, _ := json.Marshal(body)
	var (
		wg        sync.WaitGroup
		resDetail = ReplaceResp{}
		resTime   = ReplaceTimeResp{}
		resp1     []byte
		resp2     = make(map[string]interface{})
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp1, _ = Post("https://api.aichaoliuapp.cn/aiera/ai_match_trading/nft/replace/active/detail", header, jsonBytes)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp2, _ = Get("https://api.aichaoliuapp.cn/aiera/current/milli/time")
	}()
	wg.Wait()
	if len(resp1) == 0 || len(resp2) == 0 {
		return false
	}
	json.Unmarshal(resp1, &resDetail)
	g, _ := json.Marshal(resp2)
	json.Unmarshal(g, &resTime)
	//resDetail.Data.StartTimeTimestamp = 1689605160000
	diffTime := resDetail.Data.StartTimeTimestamp - resTime.CurrentMilliTime
	log.Println(resDetail.Data.StartTimeTimestamp, resTime.CurrentMilliTime, diffTime)
	if diffTime < 150 && diffTime > 0 {
		//time.Sleep(time.Millisecond * time.Duration(diffTime))
		log.Println(diffTime, resTime.CurrentMilliTime, resDetail.Data.StartTimeTimestamp, time.Now().UnixMilli())
		return true
	}
	if resDetail.Code == 0 && resDetail.Msg == "success" && resTime.CurrentMilliTime >= resDetail.Data.StartTimeTimestamp && resTime.CurrentMilliTime <= resDetail.Data.EndTimestamp {
		return true
	}
	return false
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
