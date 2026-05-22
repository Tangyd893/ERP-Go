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
