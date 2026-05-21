package repository

import "time"

// InventoryBalanceModel 库存余额 GORM 模型
type InventoryBalanceModel struct {
	ID             string    `gorm:"column:id;primaryKey"`
	TenantID       string    `gorm:"column:tenant_id;index"`
	WarehouseID    string    `gorm:"column:warehouse_id"`
	SKUID          string    `gorm:"column:sku_id"`
	SKUCode        string    `gorm:"column:sku_code"`
	TotalQuantity  int       `gorm:"column:total_quantity"`
	LockedQuantity int       `gorm:"column:locked_quantity"`
	Version        int       `gorm:"column:version;default:1"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (InventoryBalanceModel) TableName() string { return "inventory_balances" }

// Available 计算可用库存
func (m *InventoryBalanceModel) Available() int {
	return m.TotalQuantity - m.LockedQuantity
}

// InventoryLockModel 库存锁定 GORM 模型
type InventoryLockModel struct {
	ID              string    `gorm:"column:id;primaryKey"`
	TenantID        string    `gorm:"column:tenant_id"`
	OrderID         string    `gorm:"column:order_id;index"`
	SKUID           string    `gorm:"column:sku_id"`
	WarehouseID     string    `gorm:"column:warehouse_id"`
	Quantity        int       `gorm:"column:quantity"`
	ReleasedQty     int       `gorm:"column:released_quantity"`
	Status          string    `gorm:"column:status"`
	LockKey         string    `gorm:"column:lock_key;uniqueIndex"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (InventoryLockModel) TableName() string { return "inventory_locks" }

// InventoryJournalModel 库存流水 GORM 模型
type InventoryJournalModel struct {
	ID              string    `gorm:"column:id;primaryKey"`
	TenantID        string    `gorm:"column:tenant_id"`
	WarehouseID     string    `gorm:"column:warehouse_id"`
	SKUID           string    `gorm:"column:sku_id;index"`
	OrderID         string    `gorm:"column:order_id"`
	ChangeType      string    `gorm:"column:change_type"`
	ChangeQty       int       `gorm:"column:change_qty"`
	BeforeTotal     int       `gorm:"column:before_total"`
	AfterTotal      int       `gorm:"column:after_total"`
	BeforeLocked    int       `gorm:"column:before_locked"`
	AfterLocked     int       `gorm:"column:after_locked"`
	BeforeAvail     int       `gorm:"column:before_avail"`
	AfterAvail      int       `gorm:"column:after_avail"`
	IdempotencyKey  string    `gorm:"column:idempotency_key;index"`
	Operator        string    `gorm:"column:operator"`
	CreatedAt       time.Time `gorm:"column:created_at;index"`
}

func (InventoryJournalModel) TableName() string { return "inventory_journals" }
