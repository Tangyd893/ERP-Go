package domain

import (
	"testing"
)

func setupOrder() *SalesOrder {
	order := &SalesOrder{
		ID:     "order-001",
		Status: OrderPending,
		Items: []*OrderItem{
			{SKUID: "sku-1", SKUCode: "TSHIRT-001", Quantity: 2, UnitPrice: 15.99},
			{SKUID: "sku-2", SKUCode: "MUG-001", Quantity: 1, UnitPrice: 12.99},
		},
	}
	order.CalculateTotal()
	return order
}

func TestValidTransitions(t *testing.T) {
	order := setupOrder()

	// 待审核 -> 已审核
	if err := order.Approve("admin"); err != nil {
		t.Fatalf("审核失败: %v", err)
	}
	if order.Status != OrderApproved {
		t.Errorf("状态应为 approved，实际为 %s", order.Status)
	}

	// 审核后标记锁定成功 (切换到 locked 状态)
	if err := order.Transition(OrderLocked, "system", "库存锁定成功"); err != nil {
		t.Fatalf("锁定失败: %v", err)
	}
	if order.Status != OrderLocked {
		t.Errorf("状态应为 locked，实际为 %s", order.Status)
	}

	// locked -> picking
	order.Transition(OrderPicking, "warehouse", "开始拣货")
	// picking -> packed
	order.Transition(OrderPacked, "warehouse", "打包完成")
	// packed -> shipped
	order.Transition(OrderShipped, "warehouse", "已发货")
	// shipped -> delivered
	order.Transition(OrderDelivered, "system", "物流签收")
	// delivered -> completed
	order.Transition(OrderCompleted, "system", "订单完成")

	if order.Status != OrderCompleted {
		t.Errorf("最终状态应为 completed，实际为 %s", order.Status)
	}

	if len(order.StatusHistory) != 7 {
		t.Errorf("应有7条状态记录，实际 %d", len(order.StatusHistory))
	}
}

func TestInvalidTransition(t *testing.T) {
	order := setupOrder()

	// 待审核不能直接到已发货
	err := order.Transition(OrderShipped, "user", "跳过审核")
	if err == nil {
		t.Error("非法状态流转应返回错误")
	}
}

func TestCancelFromApproved(t *testing.T) {
	order := setupOrder()
	order.Approve("admin")
	if err := order.Cancel("admin", "买家取消"); err != nil {
		t.Fatalf("取消失败: %v", err)
	}
	if order.Status != OrderCancelled {
		t.Errorf("状态应为 cancelled，实际为 %s", order.Status)
	}
}

func TestAbnormalAndRetry(t *testing.T) {
	order := setupOrder()
	order.Approve("admin")

	// 审核后标记异常
	if err := order.MarkAbnormal("admin", "地址信息缺失"); err != nil {
		t.Fatalf("标记异常失败: %v", err)
	}
	if order.Status != OrderAbnormal {
		t.Errorf("状态应为 abnormal，实际为 %s", order.Status)
	}

	// 异常修复后回到待审核
	if err := order.Transition(OrderPending, "admin", "地址已补全"); err != nil {
		t.Fatalf("回退到待审核失败: %v", err)
	}
	if order.Status != OrderPending {
		t.Errorf("状态应为 pending，实际为 %s", order.Status)
	}
}

func TestOrderTotalCalculation(t *testing.T) {
	order := setupOrder()
	expectedTotal := 2*15.99 + 12.99 // 46.97
	if order.TotalAmount != expectedTotal {
		t.Errorf("订单总额应为 %.2f，实际 %.2f", expectedTotal, order.TotalAmount)
	}
}

func TestGetSKUQuantities(t *testing.T) {
	order := setupOrder()
	quantities := order.GetSKUQuantities()
	if quantities["sku-1"] != 2 {
		t.Errorf("sku-1 数量应为2，实际 %d", quantities["sku-1"])
	}
	if quantities["sku-2"] != 1 {
		t.Errorf("sku-2 数量应为1，实际 %d", quantities["sku-2"])
	}
}
