package dto

type GradOrderReq struct {
	ProductId        uint64  `json:"product_id"`          // 藏品ID
	NftProductSizeId uint64  `json:"nft_product_size_id"` // 源文件ID
	ProductName      string  `json:"product_name"`        // 藏品名称
	Price            float32 `json:"price"`               // 藏品价格
	Num              uint32  `json:"num"`                 // 数量
	AutoPay          uint32  `json:"auto_pay"`            // 自动支付
}

type StopTaskReq struct {
	ProductId uint64 `json:"product_id"` // 藏品ID
}

type CreateOrderReq struct {
	OperateType  string `json:"operate_type"`
	SecondId     uint64 `json:"second_id"`
	UserCouponId int    `json:"user_coupon_id"`
}

type SearchProductReq struct {
	Page     int32  `json:"page"`
	Keyword  string `json:"keyword"`
	PageSize int32  `json:"page_size"`
}

type WalletPayReq struct {
	PayPwd  string `json:"pay_pwd"`
	OrderId uint64 `json:"order_id"`
}
