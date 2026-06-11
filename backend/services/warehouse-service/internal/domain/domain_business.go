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

// ── Wave 波次业务规则 ──────────────────────────────────

// WaveStatus 波次状态
type WaveStatus string

const (
	WaveCreated  WaveStatus = "created"
	WavePicking  WaveStatus = "picking"
	WaveCompleted WaveStatus = "completed"
)

// NewWave 创建波次
func NewWave(id, warehouseID, name string) *Wave {
	return &Wave{
		ID:          id,
		WarehouseID: warehouseID,
		Name:        name,
		Status:      string(WaveCreated),
		CreatedAt:   time.Now(),
	}
}

// AddOutbound 向波次添加出库单
func (w *Wave) AddOutbound(outboundID string) error {
	if w.Status != string(WaveCreated) {
		return fmt.Errorf("波次状态 %s 不可添加出库单", w.Status)
	}
	for _, id := range w.OutboundIDs {
		if id == outboundID {
			return fmt.Errorf("出库单 %s 已在波次中", outboundID)
		}
	}
	w.OutboundIDs = append(w.OutboundIDs, outboundID)
	return nil
}

// StartPicking 开始波次拣货
func (w *Wave) StartPicking() error {
	if w.Status != string(WaveCreated) {
		return fmt.Errorf("波次状态 %s 不能开始拣货", w.Status)
	}
	if len(w.OutboundIDs) == 0 {
		return fmt.Errorf("波次没有出库单")
	}
	w.Status = string(WavePicking)
	return nil
}

// Complete 完成波次
func (w *Wave) Complete() error {
	if w.Status != string(WavePicking) {
		return fmt.Errorf("波次状态 %s 不能完成", w.Status)
	}
	w.Status = string(WaveCompleted)
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
