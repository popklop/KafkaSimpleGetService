package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"wbtech/internal/domain/order"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Save(ctx context.Context, o *order.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()
	query := `INSERT INTO orders (id, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT (id) DO NOTHING`
	res, err := tx.ExecContext(ctx, query, o.ID, o.TrackNumber, o.Entry, o.Locale, o.InternalSignature, o.CustomerID, o.DeliveryService, o.ShardKey, o.SmID, o.DateCreated, o.OofShard)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrOrderAlreadyExists
	}
	query = `INSERT INTO delivery(order_id, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = tx.ExecContext(ctx, query, o.ID, o.Delivery.Name, o.Delivery.Phone, o.Delivery.Zip, o.Delivery.City, o.Delivery.Address, o.Delivery.Region, o.Delivery.Email)
	if err != nil {
		return err
	}
	for i := 0; i < len(o.Items); i++ {
		query = `INSERT INTO items(order_id, chrt_id, track_number, price, rid, name, sale, size, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
		_, err := tx.ExecContext(ctx, query, o.ID, o.Items[i].ChrtID, o.Items[i].TrackNumber, o.Items[i].Price, o.Items[i].RID, o.Items[i].Name, o.Items[i].Sale, o.Items[i].Size, o.Items[i].NmID, o.Items[i].Brand, o.Items[i].Status)
		if err != nil {
			return err
		}

	}
	query = `INSERT INTO payments(order_id, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err = tx.ExecContext(ctx, query, o.ID, o.Payment.Transaction, o.Payment.RequestID, o.Payment.Currency, o.Payment.Provider, o.Payment.Amount, o.Payment.PaymentAt, o.Payment.Bank, o.Payment.DeliveryCost, o.Payment.GoodsTotal, o.Payment.CustomFee)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *OrderRepository) GetByID(ctx context.Context, id string) (*order.Order, error) {
	var o order.Order
	query := `SELECT id, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&o.ID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature, &o.CustomerID, &o.DeliveryService, &o.ShardKey, &o.SmID, &o.DateCreated, &o.OofShard)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	query = `SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_id = $1`
	row = r.db.QueryRowContext(ctx, query, id)
	var delivery order.Delivery
	err = row.Scan(&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	o.Delivery = delivery
	query = `SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payments WHERE order_id = $1`
	row = r.db.QueryRowContext(ctx, query, id)
	var payment order.Payment
	err = row.Scan(
		&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider,
		&payment.Amount, &payment.PaymentAt, &payment.Bank, &payment.DeliveryCost,
		&payment.GoodsTotal, &payment.CustomFee,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	o.Payment = payment
	query = `SELECT chrt_id, track_number, price, rid, name, sale, size, nm_id, brand, status FROM items WHERE order_id = $1`
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []order.Item
	for rows.Next() {
		var item order.Item
		err = rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale, &item.Size, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	o.Items = items

	return &o, nil
}

func (r *OrderRepository) GetAll(ctx context.Context) ([]*order.Order, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	var orders []*order.Order
	for i := 0; i < len(ids); i++ {
		o, err := r.GetByID(ctx, ids[i])
		if err != nil {
			continue
		}
		orders = append(orders, o)
	}
	return orders, nil
}
