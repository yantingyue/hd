package service

import (
	"encoding/json"

	"demo/hd/constant"
	"demo/hd/dto"
	"demo/hd/util"
)

type HdService struct {
}

var HdSrv = &HdService{}

func (*HdService) QueryProductInfoByName(r *dto.GradOrderReq) {
	header := map[string]string{
		"token": constant.Token,
	}

	req := dto.SearchProductReq{
		Page:     1,
		Keyword:  r.ProductName,
		PageSize: 1,
	}

	jsonBytes, _ := json.Marshal(req)
	body, _ := util.Post(constant.ApiHost+constant.SearchProductPath, header, jsonBytes)
	if len(body) == 0 {
		return
	}
	resp := dto.SearchProductResp{}
	json.Unmarshal(body, &resp)
	for _, v := range resp.Data {
		r.ProductId = v.ProductId
		r.NftProductSizeId = v.Id
	}
}
