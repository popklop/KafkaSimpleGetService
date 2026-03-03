package cache

import (
	"testing"
	"wbtech/internal/domain/order"
)

func TestOrderCache_SetAndGet(t *testing.T) {
	c := NewOrderCache(5)
	ord := &order.Order{ID: "1", TrackNumber: "track1"}
	c.Set("1", ord)

	got := c.Get("1")
	if got == nil {
		t.Fatal("expected order, got nil")
	}
	if got.ID != "1" {
		t.Errorf("expected ID 1, got %s", got.ID)
	}
}

func TestOrderCache_GetNonExistent(t *testing.T) {
	c := NewOrderCache(5)
	got := c.Get("missing")
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestOrderCache_Eviction(t *testing.T) {
	c := NewOrderCache(2)
	ord1 := &order.Order{ID: "1"}
	ord2 := &order.Order{ID: "2"}
	ord3 := &order.Order{ID: "3"}
	c.Set("1", ord1)
	c.Set("2", ord2)
	c.Set("3", ord3)
	if c.Get("1") != nil {
		t.Error("expected order 1 to be evicted, but it's still present")
	}
	if c.Get("2") == nil {
		t.Error("expected order 2 to be present")
	}
	if c.Get("3") == nil {
		t.Error("expected order 3 to be present")
	}
}

func TestOrderCache_UpdateExisting(t *testing.T) {
	c := NewOrderCache(2)
	ord1 := &order.Order{ID: "1", TrackNumber: "old"}
	c.Set("1", ord1)
	ord1new := &order.Order{ID: "1", TrackNumber: "new"}
	c.Set("1", ord1new)

	got := c.Get("1")
	if got.TrackNumber != "new" {
		t.Errorf("expected updated track number 'new', got %s", got.TrackNumber)
	}
}
