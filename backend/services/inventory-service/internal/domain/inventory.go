package domain

import (
	"fmt"
	"sync"
	"time"
)

const journalIDFormat = "jrnl-%s-%s"

// InventoryBalance 库存余额聚合根
type InventoryBalance struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenant_id"`
	WarehouseID    string    `json:"warehouse_id"`
	SKUID          string    `json:"sku_id"`
	SKUCode        string    `json:"sku_code"`
	TotalQuantity  int       `json:"total_quantity"`
	LockedQuantity int       `json:"locked_quantity"`
	AvailableQty   int       `json:"available_quantity"`
	Version        int       `json:"version"` // 乐观锁版本号
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// InventoryLock 库存锁定记录
type InventoryLock struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	OrderID     string    `json:"order_id"`
	SKUID       string    `json:"sku_id"`
	WarehouseID string    `json:"warehouse_id"`
	Quantity    int       `json:"quantity"`
	ReleasedQty int       `json:"released_quantity"`
	Status      string    `json:"status"` // locked, released, deducted
	LockKey     string    `json:"lock_key"` // 幂等键
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// InventoryJournal 库存流水
type InventoryJournal struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	WarehouseID   string    `json:"warehouse_id"`
	SKUID         string    `json:"sku_id"`
	OrderID       string    `json:"order_id"`
	ChangeType    string    `json:"change_type"` // lock, release, deduct, increase, adjust
	ChangeQty     int       `json:"change_qty"`
	BeforeTotal   int       `json:"before_total"`
	AfterTotal    int       `json:"after_total"`
	BeforeLocked  int       `json:"before_locked"`
	AfterLocked   int       `json:"after_locked"`
	BeforeAvail   int       `json:"before_avail"`
	AfterAvail    int       `json:"after_avail"`
	IdempotencyKey string   `json:"idempotency_key"`
	Operator      string    `json:"operator"`
	CreatedAt     time.Time `json:"created_at"`
}

// Lock 锁定库存
func (b *InventoryBalance) Lock(quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("锁定数量必须大于0")
	}
	beforeAvail := b.Available()
	if beforeAvail < quantity {
		return fmt.Errorf("库存不足: 需要%d, 可用%d", quantity, beforeAvail)
	}
	b.LockedQuantity += quantity
	b.UpdatedAt = time.Now()
	b.Version++
	return nil
}

// Release 释放库存
func (b *InventoryBalance) Release(quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("释放数量必须大于0")
	}
	if b.LockedQuantity < quantity {
		return fmt.Errorf("释放数量(%d)超过已锁定数量(%d)", quantity, b.LockedQuantity)
	}
	b.LockedQuantity -= quantity
	b.UpdatedAt = time.Now()
	b.Version++
	return nil
}

// Deduct 扣减库存（出库确认）
func (b *InventoryBalance) Deduct(quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("扣减数量必须大于0")
	}
	if b.TotalQuantity < quantity {
		return fmt.Errorf("扣减库存不足: 需要%d, 现有%d", quantity, b.TotalQuantity)
	}
	if b.LockedQuantity < quantity {
		return fmt.Errorf("已锁定库存不足: 需要扣减%d, 已锁定%d", quantity, b.LockedQuantity)
	}
	b.TotalQuantity -= quantity
	b.LockedQuantity -= quantity
	b.UpdatedAt = time.Now()
	b.Version++
	return nil
}

// Increase 增加库存
func (b *InventoryBalance) Increase(quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("增加数量必须大于0")
	}
	b.TotalQuantity += quantity
	b.UpdatedAt = time.Now()
	b.Version++
	return nil
}

// Available 计算可用库存
func (b *InventoryBalance) Available() int {
	return b.TotalQuantity - b.LockedQuantity
}

// InventoryService 库存领域服务（并发安全）
type InventoryService struct {
	mu sync.Mutex // 简化版并发控制，生产环境使用 Redis 分布式锁或数据库行锁
}

func NewInventoryService() *InventoryService {
	return &InventoryService{}
}

// LockInventory 锁定库存（事务内调用）
func (s *InventoryService) LockInventory(balance *InventoryBalance, orderID string, skuQuantities map[string]int, warehouseID string, lockKeyPrefix string) ([]*InventoryLock, []*InventoryJournal, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var locks []*InventoryLock
	var journals []*InventoryJournal

	for skuID, qty := range skuQuantities {
		beforeTotal := balance.TotalQuantity
		beforeLocked := balance.LockedQuantity
		beforeAvail := balance.Available()

		if err := balance.Lock(qty); err != nil {
			return nil, nil, err
		}

		lock := &InventoryLock{
			ID:          fmt.Sprintf("%s-%s", lockKeyPrefix, skuID),
			OrderID:     orderID,
			SKUID:       skuID,
			WarehouseID: warehouseID,
			Quantity:    qty,
			Status:      "locked",
			LockKey:     lockKeyPrefix, // 幂等键
			CreatedAt:   time.Now(),
		}
		locks = append(locks, lock)

		journal := &InventoryJournal{
			ID:            fmt.Sprintf(journalIDFormat, lockKeyPrefix, skuID),
			WarehouseID:   warehouseID,
			SKUID:         skuID,
			OrderID:       orderID,
			ChangeType:    "lock",
			ChangeQty:     qty,
			BeforeTotal:   beforeTotal,
			AfterTotal:    balance.TotalQuantity,
			BeforeLocked:  beforeLocked,
			AfterLocked:   balance.LockedQuantity,
			BeforeAvail:   beforeAvail,
			AfterAvail:    balance.Available(),
			IdempotencyKey: lockKeyPrefix,
			CreatedAt:     time.Now(),
		}
		journals = append(journals, journal)
	}

	return locks, journals, nil
}

// ReleaseInventory 释放库存
func (s *InventoryService) ReleaseInventory(balance *InventoryBalance, locks []*InventoryLock, releaseKeyPrefix string) ([]*InventoryJournal, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var journals []*InventoryJournal

	for _, lock := range locks {
		releaseQty := lock.Quantity - lock.ReleasedQty
		if releaseQty <= 0 {
			continue
		}

		beforeTotal := balance.TotalQuantity
		beforeLocked := balance.LockedQuantity
		beforeAvail := balance.Available()

		if err := balance.Release(releaseQty); err != nil {
			return nil, err
		}

		lock.ReleasedQty += releaseQty
		lock.Status = "released"
		lock.UpdatedAt = time.Now()

		journal := &InventoryJournal{
			ID:            fmt.Sprintf(journalIDFormat, releaseKeyPrefix, lock.SKUID),
			WarehouseID:   lock.WarehouseID,
			SKUID:         lock.SKUID,
			OrderID:       lock.OrderID,
			ChangeType:    "release",
			ChangeQty:     releaseQty,
			BeforeTotal:   beforeTotal,
			AfterTotal:    balance.TotalQuantity,
			BeforeLocked:  beforeLocked,
			AfterLocked:   balance.LockedQuantity,
			BeforeAvail:   beforeAvail,
			AfterAvail:    balance.Available(),
			IdempotencyKey: releaseKeyPrefix,
			CreatedAt:     time.Now(),
		}
		journals = append(journals, journal)
	}

	return journals, nil
}

// DeductInventory 扣减库存（出库确认后）
func (s *InventoryService) DeductInventory(balance *InventoryBalance, locks []*InventoryLock, deductKeyPrefix string) ([]*InventoryJournal, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var journals []*InventoryJournal

	for _, lock := range locks {
		deductQty := lock.Quantity - lock.ReleasedQty
		if deductQty <= 0 {
			continue
		}

		beforeTotal := balance.TotalQuantity
		beforeLocked := balance.LockedQuantity
		beforeAvail := balance.Available()

		if err := balance.Deduct(deductQty); err != nil {
			return nil, err
		}

		lock.Status = "deducted"
		lock.UpdatedAt = time.Now()

		journal := &InventoryJournal{
			ID:            fmt.Sprintf(journalIDFormat, deductKeyPrefix, lock.SKUID),
			WarehouseID:   lock.WarehouseID,
			SKUID:         lock.SKUID,
			OrderID:       lock.OrderID,
			ChangeType:    "deduct",
			ChangeQty:     deductQty,
			BeforeTotal:   beforeTotal,
			AfterTotal:    balance.TotalQuantity,
			BeforeLocked:  beforeLocked,
			AfterLocked:   balance.LockedQuantity,
			BeforeAvail:   beforeAvail,
			AfterAvail:    balance.Available(),
			IdempotencyKey: deductKeyPrefix,
			CreatedAt:     time.Now(),
		}
		journals = append(journals, journal)
	}

	return journals, nil
}
