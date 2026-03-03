package order

import (
	"testing"
	"time"
)

func TestMapOrder_Success(t *testing.T) {
	dto := &OrderDTO{
		OrderUID:    "123",
		TrackNumber: "track",
		Entry:       "entry",
		Delivery: DeliveryDTO{
			Name:  "Ivan",
			Email: "IvanIvanov@test.com",
		},
		Payment: PaymentDTO{
			Amount:    100,
			PaymentDt: 1600000000,
		},
		Items: []ItemDTO{
			{ChrtID: 1, Name: "item1"},
		},
		Locale:      "en",
		CustomerID:  "cust1",
		DateCreated: "2023-01-01T00:00:00Z",
		OofShard:    "1",
	}
	ord, err := MapOrder(dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ord.ID != "123" {
		t.Errorf("expected ID 123, got %s", ord.ID)
	}
	if ord.Delivery.Email != "IvanIvanov@test.com" {
		t.Errorf("expected email IvanIvanov@test.com, got %s", ord.Delivery.Email)
	}
	if ord.Payment.Amount != 100 {
		t.Errorf("expected amount 100, got %d", ord.Payment.Amount)
	}
	expectedTime := time.Unix(1600000000, 0)
	if !ord.Payment.PaymentAt.Equal(expectedTime) {
		t.Errorf("expected payment time %v, got %v", expectedTime, ord.Payment.PaymentAt)
	}
	if len(ord.Items) != 1 || ord.Items[0].ChrtID != 1 {
		t.Error("items not mapped correctly")
	}
}

func TestMapOrder_InvalidDate(t *testing.T) {
	dto := &OrderDTO{
		DateCreated: "invalid-date",
	}
	_, err := MapOrder(dto)
	if err == nil {
		t.Error("expected error for invalid date, got nil")
	}
}
