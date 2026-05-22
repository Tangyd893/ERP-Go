package domain

// CalculateTotalCost 计算总成本
func (p *ProfitReport) CalculateTotalCost() {
	p.TotalCost = p.PurchaseCost + p.ShippingCost + p.CommissionCost + p.OtherCost
}

// CalculateProfit 计算利润和利润率
func (p *ProfitReport) CalculateProfit() {
	p.GrossProfit = p.SaleAmount - p.TotalCost
	if p.SaleAmount > 0 {
		p.ProfitMargin = (p.GrossProfit / p.SaleAmount) * 100
	}
}

// Calculate 执行完整利润计算（成本 → 利润 → 利润率）
func (p *ProfitReport) Calculate() {
	p.CalculateTotalCost()
	p.CalculateProfit()
}
