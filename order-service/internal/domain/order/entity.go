package order

import (
	"strings"
	"time"
)

type Order struct {
	ID                string
	TrackNumber       string
	Entry             string
	Delivery          Delivery
	Payment           Payment
	Items             []Item
	Locale            string
	InternalSignature string
	CustomerID        string
	DeliveryService   string
	ShardKey          string
	SmID              int
	DateCreated       time.Time
	OofShard          string
}
type NewOrderParams struct {
	ID          string
	TrackNumber string
	Entry       string
	Delivery    Delivery
	Payment     Payment
	Items       []Item
	Locale      string
	CustomerID  string
	CreatedAt   time.Time
	OofShard    string
}
type Payment struct {
	Transaction  string
	RequestID    string
	Currency     string
	Provider     string
	Amount       int
	PaymentAt    time.Time
	Bank         string
	DeliveryCost int
	GoodsTotal   int
	CustomFee    int
}

type Item struct {
	ChrtID      int
	TrackNumber string
	Price       int
	RID         string
	Name        string
	Sale        int
	Size        string
	NmID        int
	Brand       string
	Status      int
}

type Delivery struct {
	Name    string
	Phone   string
	Zip     string
	City    string
	Address string
	Region  string
	Email   string
}

func NewOrder(o NewOrderParams) (*Order, error) {
	if !strings.Contains(o.Delivery.Email, "@") || !strings.Contains(o.Delivery.Email, ".") {
		return nil, ErrInvalidMail
	}
	if o.Payment.Amount <= 0 {
		return nil, ErrInvalidPayment
	}
	if o.ID == "" {
		return nil, ErrEmptyID
	}
	if len(o.Items) == 0 {
		return nil, ErrNoItems
	}

	return &Order{
		ID:          o.ID,
		TrackNumber: o.TrackNumber,
		Entry:       o.Entry,
		Delivery:    o.Delivery,
		Payment:     o.Payment,
		Items:       o.Items,
		Locale:      o.Locale,
		CustomerID:  o.CustomerID,
		DateCreated: o.CreatedAt,
		OofShard:    o.OofShard,
	}, nil
}
