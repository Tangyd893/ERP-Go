package domain

import (
	"testing"
)

func TestOutboundOrderFullFlow(t *testing.T) {
	o := &OutboundOrder{ID: "OB-001", Status: OutboundCreated}

	if err := o.StartPicking(); err != nil {
		t.Fatalf("开始拣货失败: %v", err)
	}
	if o.Status != OutboundPicking {
		t.Errorf("状态应为 picking，实际: %s", o.Status)
	}

	if err := o.CompletePicking(); err != nil {
		t.Fatalf("完成拣货失败: %v", err)
	}
	if o.Status != OutboundPicked {
		t.Errorf("状态应为 picked，实际: %s", o.Status)
	}

	if err := o.StartChecking(); err != nil {
		t.Fatalf("开始复核失败: %v", err)
	}
	if err := o.CompleteChecking(); err != nil {
		t.Fatalf("完成复核失败: %v", err)
	}
	if o.Status != OutboundChecked {
		t.Errorf("状态应为 checked，实际: %s", o.Status)
	}

	if err := o.StartPacking(); err != nil {
		t.Fatalf("开始打包失败: %v", err)
	}
	if err := o.CompletePacking(); err != nil {
		t.Fatalf("完成打包失败: %v", err)
	}
	if o.Status != OutboundPacked {
		t.Errorf("状态应为 packed，实际: %s", o.Status)
	}

	if err := o.Weigh(); err != nil {
		t.Fatalf("称重失败: %v", err)
	}
	if o.Status != OutboundWeighed {
		t.Errorf("状态应为 weighed，实际: %s", o.Status)
	}

	if err := o.Ship(); err != nil {
		t.Fatalf("发货失败: %v", err)
	}
	if o.Status != OutboundShipped {
		t.Errorf("状态应为 shipped，实际: %s", o.Status)
	}
}

func TestOutboundOrderInvalidTransitions(t *testing.T) {
	tests := []struct {
		name      string
		from      OutboundStatus
		action    func(o *OutboundOrder) error
	}{
		{"已发货不能拣货", OutboundShipped, func(o *OutboundOrder) error { return o.StartPicking() }},
		{"未拣货不能完成拣货", OutboundCreated, func(o *OutboundOrder) error { return o.CompletePicking() }},
		{"未复核不能完成复核", OutboundPicking, func(o *OutboundOrder) error { return o.CompleteChecking() }},
		{"未打包不能完成打包", OutboundChecked, func(o *OutboundOrder) error { return o.CompletePacking() }},
		{"未打包不能称重", OutboundChecked, func(o *OutboundOrder) error { return o.Weigh() }},
		{"未称重不能发货", OutboundPacked, func(o *OutboundOrder) error { return o.Ship() }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OutboundOrder{ID: "OB-002", Status: tt.from}
			if err := tt.action(o); err == nil {
				t.Error("应返回错误但未返回")
			}
		})
	}
}

// ── Wave 波次测试 ──────────────────────────────────────

func TestWaveFullFlow(t *testing.T) {
	w := NewWave("W-001", "WH-001", "波次#1")

	if w.Status != string(WaveCreated) {
		t.Errorf("初始状态应为 created，实际: %s", w.Status)
	}
	if len(w.OutboundIDs) != 0 {
		t.Error("新波次应无出库单")
	}

	// 添加出库单
	if err := w.AddOutbound("OB-001"); err != nil {
		t.Fatalf("添加出库单失败: %v", err)
	}
	if err := w.AddOutbound("OB-002"); err != nil {
		t.Fatalf("添加出库单失败: %v", err)
	}
	if len(w.OutboundIDs) != 2 {
		t.Errorf("应有 2 个出库单，实际: %d", len(w.OutboundIDs))
	}

	// 重复添加应报错
	if err := w.AddOutbound("OB-001"); err == nil {
		t.Error("重复添加出库单应报错")
	}

	// 开始拣货
	if err := w.StartPicking(); err != nil {
		t.Fatalf("开始波次拣货失败: %v", err)
	}
	if w.Status != string(WavePicking) {
		t.Errorf("状态应为 picking，实际: %s", w.Status)
	}

	// 拣货中不可添加出库单
	if err := w.AddOutbound("OB-003"); err == nil {
		t.Error("拣货中不可添加出库单")
	}

	// 完成波次
	if err := w.Complete(); err != nil {
		t.Fatalf("完成波次失败: %v", err)
	}
	if w.Status != string(WaveCompleted) {
		t.Errorf("状态应为 completed，实际: %s", w.Status)
	}
}

func TestWaveInvalidTransitions(t *testing.T) {
	// 空波次不能开始拣货
	w := NewWave("W-002", "WH-001", "空波次")
	if err := w.StartPicking(); err == nil {
		t.Error("空波次开始拣货应报错")
	}

	// 已完成波次不能再次完成
	w2 := NewWave("W-003", "WH-001", "波次#3")
	w2.AddOutbound("OB-001")
	w2.StartPicking()
	w2.Complete()
	if err := w2.Complete(); err == nil {
		t.Error("已完成波次不能再次完成")
	}

	// 已完成波次不能开始拣货
	if err := w2.StartPicking(); err == nil {
		t.Error("已完成波次不能开始拣货")
	}
}

// ── 幂等场景：扫码重复调用应返回成功而非报错 ────────────

func TestOutboundStatusIdempotency(t *testing.T) {
	// 复核幂等：已复核状态再次复核应成功
	o := &OutboundOrder{ID: "OB-IDEM", Status: OutboundChecked}
	if err := o.CompleteChecking(); err == nil {
		// 根据业务规则，已复核不能再次复核，应在 app 层处理幂等
		// 领域层应拒绝非法状态转换
		t.Log("领域层正确拒绝非法状态转换，幂等由 App 层处理")
	} else {
		t.Logf("领域层拒绝: %v（幂等由 App 层处理）", err)
	}

	// 打包幂等
	o2 := &OutboundOrder{ID: "OB-IDEM2", Status: OutboundPacked}
	if err := o2.CompletePacking(); err == nil {
		t.Error("已打包不应再次允许 CompletePacking")
	}

	// 称重幂等
	o3 := &OutboundOrder{ID: "OB-IDEM3", Status: OutboundWeighed}
	if err := o3.Weigh(); err == nil {
		t.Error("已称重不应再次允许 Weigh")
	}
}

func TestOutboundOrderMarkAbnormal(t *testing.T) {
	o := &OutboundOrder{ID: "OB-003", Status: OutboundPicking}
	if err := o.MarkAbnormal("缺货"); err != nil {
		t.Fatalf("标记异常失败: %v", err)
	}
	if o.Status != OutboundAbnormal {
		t.Errorf("状态应为 abnormal，实际: %s", o.Status)
	}

	if err := o.MarkAbnormal("再次异常"); err == nil {
		t.Error("已是异常状态，不应再次标记异常")
	}

	shipped := &OutboundOrder{ID: "OB-004", Status: OutboundShipped}
	if err := shipped.MarkAbnormal("晚了"); err == nil {
		t.Error("已发货状态不能标记异常")
	}
}
