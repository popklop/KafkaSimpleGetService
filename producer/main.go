package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/segmentio/kafka-go"
)

type Order struct {
	OrderUID    string   `json:"order_uid"`
	TrackNumber string   `json:"track_number"`
	Entry       string   `json:"entry"`
	Delivery    Delivery `json:"delivery"`
	Payment     Payment  `json:"payment"`
	Items       []Item   `json:"items"`
	Locale      string   `json:"locale"`
	InternalSig string   `json:"internal_signature"`
	CustomerID  string   `json:"customer_id"`
	DeliverySvc string   `json:"delivery_service"`
	ShardKey    string   `json:"shardkey"`
	SmID        int      `json:"sm_id"`
	DateCreated string   `json:"date_created"`
	OofShard    string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
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

type Item struct {
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

func main() {
	broker := flag.String("broker", "localhost:9092", "Kafka broker")
	topic := flag.String("topic", "orders", "Kafka topic")
	count := flag.Int("n", 1, "Number of messages to send")
	flag.Parse()
	gofakeit.Seed(0)

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{*broker},
		Topic:    *topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	for i := 0; i < *count; i++ {
		order, _ := generateFakeOrder(i + 1)
		data, _ := json.Marshal(order)
		msg := kafka.Message{
			Key:   []byte(order.OrderUID),
			Value: data,
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := writer.WriteMessages(ctx, msg)
		cancel()
		if err != nil {
			log.Fatalf("failed to write message %d: %v", i+1, err)
		}
		log.Printf("message %d sent: %s", i+1, order.OrderUID)
	}
}

func generateFakeOrder(seq int) (*Order, error) {
	now := time.Now().UTC()
	orderUID := gofakeit.UUID()
	trackNumber, err := gofakeit.Generate("TRACK-??????")
	if err != nil {
		return nil, err
	}

	delivery := Delivery{
		Name:    gofakeit.Name(),
		Phone:   gofakeit.Phone(),
		Zip:     gofakeit.Zip(),
		City:    gofakeit.City(),
		Address: gofakeit.Street() + ", " + gofakeit.StreetNumber(),
		Region:  gofakeit.State(),
		Email:   gofakeit.Email(),
	}
	amount := gofakeit.Number(1000, 5000)
	payment := Payment{
		Transaction:  gofakeit.UUID(),
		RequestID:    gofakeit.UUID(),
		Currency:     gofakeit.CurrencyShort(),
		Provider:     gofakeit.RandomString([]string{"visa", "wbpay", "mir"}),
		Amount:       amount,
		PaymentDt:    now.Unix(),
		Bank:         gofakeit.Company(),
		DeliveryCost: gofakeit.Number(100, 500),
		GoodsTotal:   amount - gofakeit.Number(100, 500),
		CustomFee:    0,
	}

	itemsCount := gofakeit.Number(1, 3)
	items := make([]Item, itemsCount)
	for i := 0; i < itemsCount; i++ {
		items[i] = Item{
			ChrtID:      gofakeit.Number(1000, 9999),
			TrackNumber: trackNumber,
			Price:       gofakeit.Number(100, 1000),
			RID:         gofakeit.UUID(),
			Name:        gofakeit.ProductName(),
			Sale:        gofakeit.Number(0, 30),
			Size:        gofakeit.RandomString([]string{"S", "M", "L", "XL"}),
			NmID:        gofakeit.Number(10000, 99999),
			Brand:       gofakeit.Company(),
			Status:      202,
		}
	}

	return &Order{
		OrderUID:    orderUID,
		TrackNumber: trackNumber,
		Entry:       "WBIL",
		Delivery:    delivery,
		Payment:     payment,
		Items:       items,
		Locale:      "en",
		InternalSig: "",
		CustomerID:  gofakeit.UUID(),
		DeliverySvc: gofakeit.RandomString([]string{"wbdelivery", "sdek", "pochta_rossii"}),
		ShardKey:    gofakeit.RandomString([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}),
		SmID:        gofakeit.Number(1, 100),
		DateCreated: now.Format(time.RFC3339),
		OofShard:    gofakeit.RandomString([]string{"0", "1", "2"}),
	}, nil
}
