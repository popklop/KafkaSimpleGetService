package order

import (
	"context"
	"errors"
	"testing"
	"wbtech/internal/domain/order"
	"wbtech/internal/infrastructure/cache"
)

type mockRepository struct {
	saveFunc    func(ctx context.Context, o *order.Order) error
	getByIDFunc func(ctx context.Context, id string) (*order.Order, error)
	getAllFunc  func(ctx context.Context) ([]*order.Order, error)
}

func (m *mockRepository) Save(ctx context.Context, o *order.Order) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, o)
	}
	return nil
}
func (m *mockRepository) GetByID(ctx context.Context, id string) (*order.Order, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}
func (m *mockRepository) GetAll(ctx context.Context) ([]*order.Order, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx)
	}
	return nil, errors.New("not implemented")
}

func TestOrderUseCase_GetOrder_CacheHit(t *testing.T) {
	c := cache.NewOrderCache(10)
	expectedOrder := &order.Order{ID: "123", TrackNumber: "track"}
	c.Set("123", expectedOrder)
	repo := &mockRepository{
		getByIDFunc: func(ctx context.Context, id string) (*order.Order, error) {
			t.Error("GetByID should not be called on cache hit")
			return nil, nil
		},
	}
	uc := NewOrderUseCase(repo, c)
	got, err := uc.GetOrder(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != expectedOrder {
		t.Errorf("expected %v, got %v", expectedOrder, got)
	}
}

func TestOrderUseCase_GetOrder_CacheMiss(t *testing.T) {
	c := cache.NewOrderCache(10)
	expectedOrder := &order.Order{ID: "123", TrackNumber: "track"}
	repo := &mockRepository{
		getByIDFunc: func(ctx context.Context, id string) (*order.Order, error) {
			if id != "123" {
				t.Errorf("expected id 123, got %s", id)
			}
			return expectedOrder, nil
		},
	}
	uc := NewOrderUseCase(repo, c)
	got, err := uc.GetOrder(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != expectedOrder {
		t.Errorf("expected %v, got %v", expectedOrder, got)
	}
	if cached := c.Get("123"); cached != expectedOrder {
		t.Error("order not cached after Get")
	}
}

func TestOrderUseCase_GetOrder_NotFound(t *testing.T) {
	c := cache.NewOrderCache(10)
	repo := &mockRepository{
		getByIDFunc: func(ctx context.Context, id string) (*order.Order, error) {
			return nil, errors.New("some db error")
		},
	}
	uc := NewOrderUseCase(repo, c)
	_, err := uc.GetOrder(context.Background(), "123")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestOrderUseCase_SaveOrder(t *testing.T) {
	c := cache.NewOrderCache(10)
	var savedOrder *order.Order
	repo := &mockRepository{
		saveFunc: func(ctx context.Context, o *order.Order) error {
			savedOrder = o
			return nil
		},
	}
	uc := NewOrderUseCase(repo, c)

	dto := &OrderDTO{
		OrderUID:    "123",
		TrackNumber: "track",
		Delivery:    DeliveryDTO{Email: "test@test.com"},
		Payment:     PaymentDTO{Amount: 100},
		Items:       []ItemDTO{{ChrtID: 1}},
		DateCreated: "2023-01-01T00:00:00Z",
	}
	err := uc.SaveOrder(context.Background(), dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if savedOrder == nil || savedOrder.ID != "123" {
		t.Errorf("expected order to be saved, got %v", savedOrder)
	}
	if cached := c.Get("123"); cached == nil {
		t.Error("order not cached after Save")
	}
}

func TestOrderUseCase_RestoreCache(t *testing.T) {
	c := cache.NewOrderCache(10)
	ordersFromDB := []*order.Order{
		{ID: "1", TrackNumber: "t1"},
		{ID: "2", TrackNumber: "t2"},
	}
	repo := &mockRepository{
		getAllFunc: func(ctx context.Context) ([]*order.Order, error) {
			return ordersFromDB, nil
		},
	}
	uc := NewOrderUseCase(repo, c)
	err := uc.RestoreCache(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Get("1") == nil {
		t.Error("order 1 not cached")
	}
	if c.Get("2") == nil {
		t.Error("order 2 not cached")
	}
}
