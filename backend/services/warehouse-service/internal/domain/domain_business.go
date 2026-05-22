package domain

import (
	"fmt"
	"time"
)

// StartPicking 开始拣货
func (o *OutboundOrder) StartPicking() error {
	if o.Status != OutboundCreated && o.Status != OutboundWaved {
		return fmt.Errorf("出库单状态 %s 不能开始拣货", o.Status)
	}
	o.Status = OutboundPicking
	o.UpdatedAt = time.Now()
	return nil
}

// CompletePicking 完成拣货
func (o *OutboundOrder) CompletePicking() error {
	if o.Status != OutboundPicking {
		return fmt.Errorf("出库单状态 %s 不能完成拣货", o.Status)
	}
	o.Status = OutboundPicked
	o.UpdatedAt = time.Now()
	return nil
}

// StartChecking 开始复核
func (o *OutboundOrder) StartChecking() error {
	if o.Status != OutboundPicked {
		return fmt.Errorf("出库单状态 %s 不能开始复核", o.Status)
	}
	o.Status = OutboundChecking
	o.UpdatedAt = time.Now()
	return nil
}

// CompleteChecking 完成复核
func (o *OutboundOrder) CompleteChecking() error {
	if o.Status != OutboundChecking {
		return fmt.Errorf("出库单状态 %s 不能完成复核", o.Status)
	}
	o.Status = OutboundChecked
	o.UpdatedAt = time.Now()
	return nil
}

// StartPacking 开始打包
func (o *OutboundOrder) StartPacking() error {
	if o.Status != OutboundChecked {
		return fmt.Errorf("出库单状态 %s 不能开始打包", o.Status)
	}
	o.Status = OutboundPacking
	o.UpdatedAt = time.Now()
	return nil
}

// CompletePacking 完成打包
func (o *OutboundOrder) CompletePacking() error {
	if o.Status != OutboundPacking {
		return fmt.Errorf("出库单状态 %s 不能完成打包", o.Status)
	}
	o.Status = OutboundPacked
	o.UpdatedAt = time.Now()
	return nil
}

// Weigh 称重
func (o *OutboundOrder) Weigh() error {
	if o.Status != OutboundPacked {
		return fmt.Errorf("出库单状态 %s 不能称重", o.Status)
	}
	o.Status = OutboundWeighed
	o.UpdatedAt = time.Now()
	return nil
}

// Ship 发货出库
func (o *OutboundOrder) Ship() error {
	if o.Status != OutboundWeighed {
		return fmt.Errorf("出库单状态 %s 不能发货", o.Status)
	}
	o.Status = OutboundShipped
	o.UpdatedAt = time.Now()
	return nil
}

// MarkAbnormal 标记异常
func (o *OutboundOrder) MarkAbnormal(reason string) error {
	if o.Status == OutboundShipped || o.Status == OutboundAbnormal {
		return fmt.Errorf("出库单状态 %s 不能标记异常", o.Status)
	}
	_ = reason
	o.Status = OutboundAbnormal
	o.UpdatedAt = time.Now()
	return nil
}
