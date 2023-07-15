package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"demo/hd/constant"
	"demo/hd/dto"
	"demo/hd/util"
)

type TaskJob struct {
	Ctx              context.Context `json:"ctx"`
	QuitCh           chan struct{}   `json:"quit_ch"`
	TaskId           uint64          `json:"task_id"`
	Num              uint32          `json:"num"`
	ProductId        uint64          `json:"product_id"`
	NftProductSizeId uint64          `json:"nft_product_size_id"`
	ProductName      string          `json:"product_name"`
	Price            float32         `json:"price"`
	AutoPay          uint32          `json:"auto_pay"`
}

var (
	TaskId             uint64
	TaskIdToTaskJobMap = sync.Map{}
	TaskJobCh          = make(chan TaskJob, 1024)
)

func InitTask() {
	for task := range TaskJobCh {
		go GradOrder(&task)
	}
}

func GradOrder(req *TaskJob) {
	TaskIdToTaskJobMap.Store(req.ProductId, req)
	timer := time.NewTimer(time.Second)
	for {
		select {
		case <-req.Ctx.Done():
			return
		case <-req.QuitCh:
			TaskIdToTaskJobMap.Delete(req.ProductId)
			return
		case <-timer.C:
			monitorSecondList(req)
		}
		timer.Reset(time.Second * 3)
	}
}

func monitorSecondList(req *TaskJob) {
	if req.Num == 0 {
		req.QuitCh <- struct{}{}
		return
	}
	header := util.GenerateHeader(constant.Token)
	payload := map[string]interface{}{
		"product_id":          req.ProductId,
		"pageSize":            1 + req.Num,
		"unlock":              1,
		"pageNumber":          1,
		"nft_product_size_id": req.NftProductSizeId,
		"order_by":            "price",
		"prop_pack":           0,
	}

	jsonBytes, _ := json.Marshal(payload)
	body, _ := util.Post(constant.ApiHost+constant.SecondSellListPath, header, jsonBytes)
	if len(body) == 0 {
		log.Println("没有在寄售...")
		return
	}
	monitorResp := dto.MonitorProductListResp{}
	_ = json.Unmarshal(body, &monitorResp)
	if monitorResp.Code != 0 {
		log.Println(monitorResp.Msg)
		return
	}
	if len(monitorResp.Data.Res) == 0 {
		log.Println("没有寄售中藏品。。。")
	}
	for _, v := range monitorResp.Data.Res {
		if v.Price > req.Price {
			log.Println("========", v.Price, req.Price)
			return
		}
		if v.Price <= req.Price {
			if req.Num == 0 {
				req.QuitCh <- struct{}{}
				return
			}
			// 创建订单
			if createOrder(CreateOrderReq{
				SecondId:    v.SecondId,
				ProductName: req.ProductName,
				AutoPay:     req.AutoPay,
			}) {
				req.Num--
				log.Println("创建订单成功 num = ", req.Num)
			}
		}
	}
}

type CreateOrderReq struct {
	SecondId    uint64 `json:"second_id"`
	ProductName string `json:"product_name"`
	AutoPay     uint32 `json:"auto_pay"`
}

func createOrder(req CreateOrderReq) bool {
	header := util.GenerateCreateOrderHeader(constant.Token)
	body := dto.CreateOrderReq{
		OperateType: "buy",
		SecondId:    req.SecondId,
	}
	jsonBytes, _ := json.Marshal(body)
	jsonBytes, _ = util.Post(constant.ApiHost+constant.CreateOrderPath, header, jsonBytes)
	if len(jsonBytes) == 0 {
		return false
	}
	resp := dto.CreateOrderResp{}
	_ = json.Unmarshal(jsonBytes, &resp)
	fmt.Println(resp)
	if resp.Code == 0 {
		if req.AutoPay == 1 {
			if pay(resp.Data.OrderId) {
				notice("text", fmt.Sprintf("{\"text\":\"藏品《%s》%d\"}", req.ProductName, resp.Data.OrderId))

				//notice(req.ProductName, fmt.Sprintf("<%s>: %d", req.ProductName, resp.Data.OrderId))
				return true
			}
		}
		notice("text", fmt.Sprintf("{\"text\":\"藏品《%s》%d\"}", req.ProductName, resp.Data.OrderId))
		return true
	}
	return false
}

const (
	noticeUrl = "https://open.feishu.cn/open-apis/bot/v2/hook/f8e1aa28-f3fd-4bdc-aae9-257c773a7cdd"
)

type FSMsg struct {
	MsgType string `json:"msg_type"`
	Content string `json:"content"`
}

func notice(title, msg string) {
	m := FSMsg{
		MsgType: title,
		Content: msg,
	}
	fmt.Println(m)
	jsonBytes, _ := json.Marshal(m)
	_, _ = util.Post(noticeUrl, nil, jsonBytes)
}

func pay(orderId uint64) bool {
	header := util.GenerateHeader(constant.Token)
	req := dto.WalletPayReq{
		PayPwd:  constant.Pwd,
		OrderId: orderId,
	}

	jsonBytes, _ := json.Marshal(req)
	body, _ := util.Post(constant.ApiHost+constant.WalletPayPath, header, jsonBytes)
	resp := dto.WalletPayResp{}
	_ = json.Unmarshal(body, &resp)
	if resp.Code == 0 {
		log.Println(orderId, "购买成功")
		return true
	}
	log.Println("购买失败")
	return false
}
