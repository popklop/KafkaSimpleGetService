package order

import (
	"context"
	"errors"
	"log"
	"wbtech/internal/domain/order"
	"wbtech/internal/infrastructure/cache"
	"wbtech/internal/infrastructure/postgres"
	"wbtech/metrics"
)

type OrderUseCase struct {
	repo  order.Repository
	cache *cache.OrderCache
}

func NewOrderUseCase(Repo order.Repository, Cache *cache.OrderCache) *OrderUseCase {
	return &OrderUseCase{
		repo:  Repo,
		cache: Cache,
	}
}

func (uc *OrderUseCase) SaveOrder(ctx context.Context, dto *OrderDTO) error {
	domainData, err := MapOrder(dto)
	if err != nil {
		return err
	}
	if err := uc.repo.Save(ctx, domainData); err != nil {
		return err
	}
	uc.cache.Set(domainData.ID, domainData)
	return nil
}

func (uc *OrderUseCase) GetOrder(ctx context.Context, id string) (*order.Order, error) {
	if cached := uc.cache.Get(id); cached != nil {
		log.Print("Got data from cache")
		metrics.IncCacheHit()
		return cached, nil
	}
	order, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	uc.cache.Set(order.ID, order)
	log.Print("Got data from db, and cached")
	return order, nil
}

func (uc *OrderUseCase) RestoreCache(ctx context.Context) error {
	orders, err := uc.repo.GetAll(ctx)
	if err != nil {
		return err
	}
	i := 0
	for i < len(orders) {
		uc.cache.Set(orders[i].ID, orders[i])
		i++
	}
	log.Printf("Restored %d elements from cache", i)
	return nil
}
