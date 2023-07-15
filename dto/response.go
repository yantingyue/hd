package dto

type GradOrderResp struct {
}

type MonitorProductListResp struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Res []InSellProduct `json:"res"`
	} `json:"data"`
}

type InSellProduct struct {
	SecondId         uint64        `json:"second_id"`
	ProductId        uint64        `json:"product_id"`
	SizeId           uint64        `json:"size_id"`
	NftProductSizeId uint64        `json:"nft_product_size_id"`
	OrderId          uint64        `json:"order_id"`
	IsLock           int32         `json:"is_lock"`
	Price            float32       `json:"price"`
	ProductTitle     string        `json:"product_title"`
	ReceiverCity     string        `json:"receiver_city"`
	ReceiverProvince string        `json:"receiver_province"`
	ReceiverRegion   string        `json:"receiver_region"`
	LoginName        string        `json:"login_name"`
	UserPicUrl       string        `json:"user_pic_url"`
	ProductPicture   string        `json:"product_picture"`
	IsDiscount       uint64        `json:"is_discount"`
	SellContent      uint64        `json:"sell_content"`
	SellerUserId     uint64        `json:"seller_user_id"`
	Extend           []interface{} `json:"extend"`
	IsOwner          uint64        `json:"is_owner"`
}

type GrabOrderTaskInfo struct {
	ProductId        uint64  `json:"product_id"`
	NftProductSizeId uint64  `json:"nft_product_size_id"`
	ProductName      string  `json:"product_name"`
	Price            float32 `json:"price"`
	Num              uint32  `json:"num"`
}

type CreateOrderResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		OrderId      uint64 `json:"order_id"`
		Status       int    `json:"status"`
		PayTypeThird int    `json:"pay_type_third"`
	} `json:"data"`
}

type SearchProductResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Id          uint64 `json:"id"`
		ProductId   uint64 `json:"product_id"`
		ProductName string `json:"product_name"`
		Image       string `json:"image"`
		Price       string `json:"price"`
		Label       struct {
			Text      string   `json:"text"`
			TextColor string   `json:"text_color"`
			BgColor   []string `json:"bg_color"`
		} `json:"label,omitempty"`
	} `json:"data"`
}

type WalletPayResp struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data struct{} `json:"data"`
}
