package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/inventory-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// InventoryRepository GORM 实现的库存仓储
type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

// FindBalance 查询库存余额
func (r *InventoryRepository) FindBalance(ctx context.Context, warehouseID, skuID string) (*domain.InventoryBalance, error) {
	var model InventoryBalanceModel
	err := r.db.WithContext(ctx).
		Where("warehouse_id = ? AND sku_id = ?", warehouseID, skuID).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(errors.CodeSKUNotFound, "库存记录不存在")
		}
		return nil, err
	}
	return r.modelToBalance(&model), nil
}

// ListBalances 查询库存余额列表
func (r *InventoryRepository) ListBalances(ctx context.Context, tenantID string, offset, limit int) ([]*domain.InventoryBalance, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&InventoryBalanceModel{})
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []*InventoryBalanceModel
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	balances := make([]*domain.InventoryBalance, len(models))
	for i, m := range models {
		balances[i] = r.modelToBalance(m)
	}
	return balances, total, nil
}

// LockStock 锁定库存（事务内：更新余额 + 创建锁定记录 + 写入流水）
func (r *InventoryRepository) LockStock(ctx context.Context, lock *domain.InventoryLock, journal *domain.InventoryJournal) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 幂等检查
		var existingLock InventoryLockModel
		if err := tx.Where("lock_key = ?", lock.LockKey).First(&existingLock).Error; err == nil {
			return errors.NewBusinessError(errors.CodeStockIdempotencyConflict, "重复锁定请求: "+lock.LockKey)
		}

		// 行锁查询余额
		var balance InventoryBalanceModel
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("warehouse_id = ? AND sku_id = ?", lock.WarehouseID, lock.SKUID).
			First(&balance).Error
		if err != nil {
			return fmt.Errorf("查询库存失败: %w", err)
		}

		if balance.Available() < lock.Quantity {
			return errors.NewBusinessError(errors.CodeInsufficientStock, "库存不足")
		}

		beforeTotal := balance.TotalQuantity
		beforeLocked := balance.LockedQuantity
		beforeAvail := balance.Available()

		balance.LockedQuantity += lock.Quantity
		balance.Version++
		if err := tx.Save(&balance).Error; err != nil {
			return fmt.Errorf("更新库存余额失败: %w", err)
		}

		lockModel := &InventoryLockModel{
			ID:              lock.ID,
			TenantID:        lock.TenantID,
			OrderID:         lock.OrderID,
			SKUID:           lock.SKUID,
			WarehouseID:     lock.WarehouseID,
			Quantity:        lock.Quantity,
			ReleasedQty:     0,
			Status:          "locked",
			LockKey:         lock.LockKey,
			CreatedAt:       lock.CreatedAt,
		}
		if err := tx.Create(lockModel).Error; err != nil {
			return fmt.Errorf("创建锁定记录失败: %w", err)
		}

		journalModel := &InventoryJournalModel{
			ID:             journal.ID,
			TenantID:       journal.TenantID,
			WarehouseID:    journal.WarehouseID,
			SKUID:          journal.SKUID,
			OrderID:        journal.OrderID,
			ChangeType:     journal.ChangeType,
			ChangeQty:      journal.ChangeQty,
			BeforeTotal:    beforeTotal,
			AfterTotal:     balance.TotalQuantity,
			BeforeLocked:   beforeLocked,
			AfterLocked:    balance.LockedQuantity,
			BeforeAvail:    beforeAvail,
			AfterAvail:     balance.Available(),
			IdempotencyKey: journal.IdempotencyKey,
			Operator:       journal.Operator,
			CreatedAt:      journal.CreatedAt,
		}
		if err := tx.Create(journalModel).Error; err != nil {
			return fmt.Errorf("写入库存流水失败: %w", err)
		}

		return nil
	})
}

// ReleaseStock 释放库存
func (r *InventoryRepository) ReleaseStock(ctx context.Context, lockKey string, releaseQty int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lockModel InventoryLockModel
		if err := tx.Where("lock_key = ?", lockKey).First(&lockModel).Error; err != nil {
			return fmt.Errorf("锁定记录不存在: %w", err)
		}

		if lockModel.Status == "released" || lockModel.Status == "deducted" {
			return nil // 幂等：已释放或已扣减
		}

		actualRelease := lockModel.Quantity - lockModel.ReleasedQty
		if releaseQty > actualRelease {
			releaseQty = actualRelease
		}
		if releaseQty <= 0 {
			return nil
		}

		var balance InventoryBalanceModel
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("warehouse_id = ? AND sku_id = ?", lockModel.WarehouseID, lockModel.SKUID).
			First(&balance).Error
		if err != nil {
			return fmt.Errorf("查询库存失败: %w", err)
		}

		balance.LockedQuantity -= releaseQty
		balance.Version++
		if err := tx.Save(&balance).Error; err != nil {
			return fmt.Errorf("更新库存余额失败: %w", err)
		}

		lockModel.ReleasedQty += releaseQty
		lockModel.Status = "released"
		lockModel.UpdatedAt = time.Now()
		if err := tx.Save(&lockModel).Error; err != nil {
			return fmt.Errorf("更新锁定记录失败: %w", err)
		}

		return nil
	})
}

// DeductStock 扣减库存
func (r *InventoryRepository) DeductStock(ctx context.Context, lockKey string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lockModel InventoryLockModel
		if err := tx.Where("lock_key = ?", lockKey).First(&lockModel).Error; err != nil {
			return fmt.Errorf("锁定记录不存在: %w", err)
		}

		if lockModel.Status == "deducted" {
			return nil // 幂等：已扣减
		}

		deductQty := lockModel.Quantity - lockModel.ReleasedQty
		if deductQty <= 0 {
			return nil
		}

		var balance InventoryBalanceModel
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("warehouse_id = ? AND sku_id = ?", lockModel.WarehouseID, lockModel.SKUID).
			First(&balance).Error
		if err != nil {
			return fmt.Errorf("查询库存失败: %w", err)
		}

		if balance.TotalQuantity < deductQty || balance.LockedQuantity < deductQty {
			return errors.NewBusinessError(errors.CodeInsufficientStock, "库存不足，无法扣减")
		}

		balance.TotalQuantity -= deductQty
		balance.LockedQuantity -= deductQty
		balance.Version++
		if err := tx.Save(&balance).Error; err != nil {
			return fmt.Errorf("更新库存余额失败: %w", err)
		}

		lockModel.Status = "deducted"
		lockModel.UpdatedAt = time.Now()
		if err := tx.Save(&lockModel).Error; err != nil {
			return fmt.Errorf("更新锁定记录失败: %w", err)
		}

		return nil
	})
}

// FindLocksByOrderID 按订单ID查询锁定记录
func (r *InventoryRepository) FindLocksByOrderID(ctx context.Context, orderID string) ([]*domain.InventoryLock, error) {
	var models []*InventoryLockModel
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&models).Error
	if err != nil {
		return nil, err
	}
	locks := make([]*domain.InventoryLock, len(models))
	for i, m := range models {
		locks[i] = r.modelToLock(m)
	}
	return locks, nil
}

// ListJournals 查询库存流水
func (r *InventoryRepository) ListJournals(ctx context.Context, tenantID, skuID string, offset, limit int) ([]*domain.InventoryJournal, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&InventoryJournalModel{})
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if skuID != "" {
		query = query.Where("sku_id = ?", skuID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []*InventoryJournalModel
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	journals := make([]*domain.InventoryJournal, len(models))
	for i, m := range models {
		journals[i] = r.modelToJournal(m)
	}
	return journals, total, nil
}

// IncreaseStock 增加库存（入库）
func (r *InventoryRepository) IncreaseStock(ctx context.Context, journal *domain.InventoryJournal) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var balance InventoryBalanceModel
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("warehouse_id = ? AND sku_id = ?", journal.WarehouseID, journal.SKUID).
			First(&balance).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				balance = InventoryBalanceModel{
					ID:            journal.SKUID + "-" + journal.WarehouseID,
					TenantID:      journal.TenantID,
					WarehouseID:   journal.WarehouseID,
					SKUID:         journal.SKUID,
					TotalQuantity:  0,
				}
				if err := tx.Create(&balance).Error; err != nil {
					return fmt.Errorf("创建库存记录失败: %w", err)
				}
			} else {
				return fmt.Errorf("查询库存失败: %w", err)
			}
		}

		beforeTotal := balance.TotalQuantity
		balance.TotalQuantity += journal.ChangeQty
		balance.Version++
		if err := tx.Save(&balance).Error; err != nil {
			return fmt.Errorf("更新库存余额失败: %w", err)
		}

		journalModel := &InventoryJournalModel{
			ID:             journal.ID,
			TenantID:       journal.TenantID,
			WarehouseID:    journal.WarehouseID,
			SKUID:          journal.SKUID,
			ChangeType:     journal.ChangeType,
			ChangeQty:      journal.ChangeQty,
			BeforeTotal:    beforeTotal,
			AfterTotal:     balance.TotalQuantity,
			BeforeLocked:   balance.LockedQuantity,
			AfterLocked:    balance.LockedQuantity,
			BeforeAvail:    beforeTotal - balance.LockedQuantity,
			AfterAvail:     balance.Available(),
			IdempotencyKey: journal.IdempotencyKey,
			Operator:       journal.Operator,
			CreatedAt:      journal.CreatedAt,
		}
		return tx.Create(journalModel).Error
	})
}

func (r *InventoryRepository) modelToBalance(m *InventoryBalanceModel) *domain.InventoryBalance {
	return &domain.InventoryBalance{
		ID:             m.ID,
		TenantID:       m.TenantID,
		WarehouseID:    m.WarehouseID,
		SKUID:          m.SKUID,
		SKUCode:        m.SKUCode,
		TotalQuantity:  m.TotalQuantity,
		LockedQuantity: m.LockedQuantity,
		Version:        m.Version,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func (r *InventoryRepository) modelToLock(m *InventoryLockModel) *domain.InventoryLock {
	return &domain.InventoryLock{
		ID:          m.ID,
		TenantID:    m.TenantID,
		OrderID:     m.OrderID,
		SKUID:       m.SKUID,
		WarehouseID: m.WarehouseID,
		Quantity:    m.Quantity,
		ReleasedQty: m.ReleasedQty,
		Status:      m.Status,
		LockKey:     m.LockKey,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *InventoryRepository) modelToJournal(m *InventoryJournalModel) *domain.InventoryJournal {
	return &domain.InventoryJournal{
		ID:             m.ID,
		TenantID:       m.TenantID,
		WarehouseID:    m.WarehouseID,
		SKUID:          m.SKUID,
		OrderID:        m.OrderID,
		ChangeType:     m.ChangeType,
		ChangeQty:      m.ChangeQty,
		BeforeTotal:    m.BeforeTotal,
		AfterTotal:     m.AfterTotal,
		BeforeLocked:   m.BeforeLocked,
		AfterLocked:    m.AfterLocked,
		BeforeAvail:    m.BeforeAvail,
		AfterAvail:     m.AfterAvail,
		IdempotencyKey: m.IdempotencyKey,
		Operator:       m.Operator,
		CreatedAt:      m.CreatedAt,
	}
}
