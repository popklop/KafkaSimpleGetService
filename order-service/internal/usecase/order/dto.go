package order

type OrderDTO struct {
	OrderUID    string      `json:"order_uid"`
	TrackNumber string      `json:"track_number"`
	Entry       string      `json:"entry"`
	Delivery    DeliveryDTO `json:"delivery"`
	Payment     PaymentDTO  `json:"payment"`
	Items       []ItemDTO   `json:"items"`
	Locale      string      `json:"locale"`
	InternalSig string      `json:"internal_signature"`
	CustomerID  string      `json:"customer_id"`
	DeliverySvc string      `json:"delivery_service"`
	ShardKey    string      `json:"shardkey"`
	SmID        int         `json:"sm_id"`
	DateCreated string      `json:"date_created"`
	OofShard    string      `json:"oof_shard"`
}

type DeliveryDTO struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type PaymentDTO struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type ItemDTO struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}
