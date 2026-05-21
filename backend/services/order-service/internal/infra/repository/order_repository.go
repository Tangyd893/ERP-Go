package repository

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/order-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.SalesOrder) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		orderModel := &SalesOrderModel{
			ID: order.ID, TenantID: order.TenantID, StoreID: order.StoreID,
			PlatformOrderNo: order.PlatformOrderNo, OrderType: string(order.OrderType),
			OrderSource: string(order.OrderSource), OrderStatus: string(order.Status),
			BuyerName: order.BuyerName, BuyerEmail: order.BuyerEmail, Currency: order.Currency,
			TotalAmount: order.TotalAmount, ShippingAmount: order.ShippingFee,
			DiscountAmount: order.TaxAmount, IdempotencyKey: order.IdempotencyKey,
			CreatedAt: order.CreatedAt, UpdatedAt: order.UpdatedAt,
		}
		if err := tx.Create(orderModel).Error; err != nil {
			return err
		}

		for _, item := range order.Items {
			if err := tx.Create(&OrderItemModel{
				ID: item.ID, OrderID: order.ID, SKUID: item.SKUID,
				SKUCode: item.SKUCode, SKUName: item.SKUName,
				Quantity: item.Quantity, UnitPrice: item.UnitPrice, TotalPrice: item.TotalPrice,
			}).Error; err != nil {
				return err
			}
		}

		if order.Address != nil {
			if err := tx.Create(&OrderAddressModel{
				ID: order.ID + "-addr", OrderID: order.ID,
				ContactName: order.Address.ContactName, Phone: order.Address.Phone,
				Email: order.Address.Email, Country: order.Address.Country,
				State: order.Address.State, City: order.Address.City, District: order.Address.District,
				AddressLine1: order.Address.StreetLine1, AddressLine2: order.Address.StreetLine2,
				PostalCode: order.Address.PostalCode,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *OrderRepository) FindByID(ctx context.Context, id string) (*domain.SalesOrder, error) {
	var m SalesOrderModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(errors.CodeOrderNotFound, "订单不存在")
		}
		return nil, err
	}
	return r.loadOrder(ctx, &m)
}

func (r *OrderRepository) FindByIdempotencyKey(ctx context.Context, key string) (*domain.SalesOrder, error) {
	var m SalesOrderModel
	err := r.db.WithContext(ctx).Where("idempotency_key = ?", key).First(&m).Error
	if err != nil {
		return nil, err
	}
	return r.loadOrder(ctx, &m)
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).Model(&SalesOrderModel{}).Where("id = ?", id).Update("order_status", status).Error
}

func (r *OrderRepository) List(ctx context.Context, tenantID string, offset, limit int) ([]*domain.SalesOrder, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&SalesOrderModel{}).Where("tenant_id = ?", tenantID)
	query.Count(&total)

	var models []*SalesOrderModel
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&models)

	orders := make([]*domain.SalesOrder, 0, len(models))
	for i := range models {
		order, _ := r.loadOrder(ctx, models[i])
		if order != nil {
			orders = append(orders, order)
		}
	}
	return orders, total, nil
}

func (r *OrderRepository) loadOrder(ctx context.Context, m *SalesOrderModel) (*domain.SalesOrder, error) {
	var items []*OrderItemModel
	r.db.WithContext(ctx).Where("order_id = ?", m.ID).Find(&items)

	var addr OrderAddressModel
	r.db.WithContext(ctx).Where("order_id = ?", m.ID).First(&addr)

	var logs []*OrderStatusLogModel
	r.db.WithContext(ctx).Where("order_id = ?", m.ID).Order("created_at ASC").Find(&logs)

	domainItems := make([]*domain.OrderItem, len(items))
	for i, item := range items {
		domainItems[i] = &domain.OrderItem{
			ID: item.ID, OrderID: m.ID, SKUID: item.SKUID, SKUCode: item.SKUCode, SKUName: item.SKUName,
			Quantity: item.Quantity, UnitPrice: item.UnitPrice, TotalPrice: item.TotalPrice,
		}
	}

	statusLogs := make([]*domain.StatusLog, len(logs))
	for i, l := range logs {
		statusLogs[i] = &domain.StatusLog{
			FromStatus: domain.OrderStatus(l.FromStatus), ToStatus: domain.OrderStatus(l.ToStatus),
			Operator: l.Operator, Remark: l.Remark, CreatedAt: l.CreatedAt,
		}
	}

	order := &domain.SalesOrder{
		ID: m.ID, TenantID: m.TenantID, StoreID: m.StoreID,
		PlatformOrderNo: m.PlatformOrderNo, OrderType: domain.OrderType(m.OrderType),
		OrderSource: domain.OrderSource(m.OrderSource), Status: domain.OrderStatus(m.OrderStatus),
		BuyerName: m.BuyerName, BuyerEmail: m.BuyerEmail, Currency: m.Currency,
		TotalAmount: m.TotalAmount, ShippingFee: m.ShippingAmount, TaxAmount: m.DiscountAmount,
		Items: domainItems, StatusHistory: statusLogs,
		IdempotencyKey: m.IdempotencyKey, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
	if addr.ID != "" {
		order.Address = &domain.Address{
			ContactName: addr.ContactName, Phone: addr.Phone, Email: addr.Email,
			Country: addr.Country, State: addr.State, City: addr.City, District: addr.District,
			StreetLine1: addr.AddressLine1, StreetLine2: addr.AddressLine2, PostalCode: addr.PostalCode,
		}
	}

	return order, nil
}
