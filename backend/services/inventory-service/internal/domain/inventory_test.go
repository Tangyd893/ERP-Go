package domain

import (
	"sync"
	"testing"
)

func SetupBalance(initialQty int) *InventoryBalance {
	return &InventoryBalance{
		ID:            "bal-001",
		WarehouseID:   "wh-001",
		SKUID:         "sku-001",
		TotalQuantity: initialQty,
	}
}

func TestLockInventory(t *testing.T) {
	balance := SetupBalance(100)

	if err := balance.Lock(10); err != nil {
		t.Fatalf("锁定失败: %v", err)
	}

	if balance.LockedQuantity != 10 {
		t.Errorf("已锁定数量应为10，实际 %d", balance.LockedQuantity)
	}
	if balance.Available() != 90 {
		t.Errorf("可用数量应为90，实际 %d", balance.Available())
	}
}

func TestLockInsufficientStock(t *testing.T) {
	balance := SetupBalance(5)

	err := balance.Lock(10)
	if err == nil {
		t.Error("锁定超过库存应返回错误")
	}
}

func TestReleaseInventory(t *testing.T) {
	balance := SetupBalance(100)
	balance.Lock(20)

	if err := balance.Release(5); err != nil {
		t.Fatalf("释放失败: %v", err)
	}

	if balance.LockedQuantity != 15 {
		t.Errorf("已锁定数量应为15，实际 %d", balance.LockedQuantity)
	}
	if balance.Available() != 85 {
		t.Errorf("可用数量应为85，实际 %d", balance.Available())
	}
}

func TestReleaseExceedsLocked(t *testing.T) {
	balance := SetupBalance(100)
	balance.Lock(5)

	err := balance.Release(10)
	if err == nil {
		t.Error("释放超过已锁定应返回错误")
	}
}

func TestDeductInventory(t *testing.T) {
	balance := SetupBalance(100)
	balance.Lock(20)

	if err := balance.Deduct(15); err != nil {
		t.Fatalf("扣减失败: %v", err)
	}

	if balance.TotalQuantity != 85 {
		t.Errorf("总数量应为85，实际 %d", balance.TotalQuantity)
	}
	if balance.LockedQuantity != 5 {
		t.Errorf("已锁定应剩余5，实际 %d", balance.LockedQuantity)
	}
	if balance.Available() != 80 {
		t.Errorf("可用数量应为80，实际 %d", balance.Available()) // 85 total - 5 locked = 80
	}
}

func TestDeductExceedsStock(t *testing.T) {
	balance := SetupBalance(10)

	err := balance.Deduct(15)
	if err == nil {
		t.Error("扣减超过库存应返回错误")
	}
}

func TestIncreaseInventory(t *testing.T) {
	balance := SetupBalance(100)

	if err := balance.Increase(50); err != nil {
		t.Fatalf("增加库存失败: %v", err)
	}

	if balance.TotalQuantity != 150 {
		t.Errorf("总数量应为150，实际 %d", balance.TotalQuantity)
	}
}

func TestConcurrentLocking(t *testing.T) {
	balance := SetupBalance(100)
	svc := NewInventoryService()

	var wg sync.WaitGroup
	errChan := make(chan error, 5)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := balance.Lock(15); err != nil {
				errChan <- err
			}
		}(i)
	}
	wg.Wait()
	close(errChan)

	errors := 0
	for range errChan {
		errors++
	}

	// 100 / 15 = 6.67, 5 goroutines should result in only 6 successful locks (90)
	// 5 goroutines * 15 = 75, should be fine with 100 available
	if balance.LockedQuantity != 75 {
		t.Errorf("并发锁定后已锁定应为75，实际 %d", balance.LockedQuantity)
	}
	if errors > 0 {
		t.Logf("并发冲突次数: %d", errors)
	}

	_ = svc
}

func TestLockReleaseDeductFullCycle(t *testing.T) {
	balance := SetupBalance(50)

	// 锁定
	if err := balance.Lock(20); err != nil {
		t.Fatalf("锁定失败: %v", err)
	}
	if balance.Available() != 30 {
		t.Error("锁定后可用应为30")
	}

	// 部分释放
	if err := balance.Release(5); err != nil {
		t.Fatalf("释放失败: %v", err)
	}
	if balance.LockedQuantity != 15 {
		t.Error("释放后已锁定应为15")
	}

	// 扣减剩余锁定
	if err := balance.Deduct(15); err != nil {
		t.Fatalf("扣减失败: %v", err)
	}
	if balance.TotalQuantity != 35 {
		t.Errorf("扣减后总数应为35，实际 %d", balance.TotalQuantity)
	}
	if balance.LockedQuantity != 0 {
		t.Errorf("扣减后已锁定应为0，实际 %d", balance.LockedQuantity)
	}
	if balance.Available() != 35 {
		t.Error("扣减后可用应为35")
	}
}
