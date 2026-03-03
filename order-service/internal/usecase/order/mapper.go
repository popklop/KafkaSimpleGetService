package order

import (
	"time"
	domain "wbtech/internal/domain/order"
)

func MapOrder(dto *OrderDTO) (*domain.Order, error) {
	dateCreated, err := time.Parse(time.RFC3339, dto.DateCreated)
	if err != nil {
		return &domain.Order{}, err
	}

	delivery := domain.Delivery{
		Name:    dto.Delivery.Name,
		Phone:   dto.Delivery.Phone,
		Zip:     dto.Delivery.Zip,
		City:    dto.Delivery.City,
		Address: dto.Delivery.Address,
		Region:  dto.Delivery.Region,
		Email:   dto.Delivery.Email,
	}

	payment := domain.Payment{
		Transaction:  dto.Payment.Transaction,
		RequestID:    dto.Payment.RequestID,
		Currency:     dto.Payment.Currency,
		Provider:     dto.Payment.Provider,
		Amount:       dto.Payment.Amount,
		PaymentAt:    time.Unix(dto.Payment.PaymentDt, 0),
		Bank:         dto.Payment.Bank,
		DeliveryCost: dto.Payment.DeliveryCost,
		GoodsTotal:   dto.Payment.GoodsTotal,
		CustomFee:    dto.Payment.CustomFee,
	}

	var items []domain.Item
	for i := 0; i < len(dto.Items); i++ {
		item := domain.Item{
			ChrtID:      dto.Items[i].ChrtID,
			TrackNumber: dto.Items[i].TrackNumber,
			Price:       dto.Items[i].Price,
			RID:         dto.Items[i].RID,
			Name:        dto.Items[i].Name,
			Sale:        dto.Items[i].Sale,
			Size:        dto.Items[i].Size,
			NmID:        dto.Items[i].NmID,
			Brand:       dto.Items[i].Brand,
			Status:      dto.Items[i].Status,
		}
		items = append(items, item)
	}

	params := domain.NewOrderParams{
		ID:          dto.OrderUID,
		TrackNumber: dto.TrackNumber,
		Entry:       dto.Entry,
		Delivery:    delivery,
		Payment:     payment,
		Items:       items,
		Locale:      dto.Locale,
		CustomerID:  dto.CustomerID,
		CreatedAt:   dateCreated,
		OofShard:    dto.OofShard,
	}

	return domain.NewOrder(params)
}
