<script setup lang="ts">
import { ref } from "vue";

const balances = ref([
  { sku_id: "sku-001", sku_code: "TSHIRT-001", sku_name: "纯棉T恤-白色-M", warehouse: "美国仓A", total_quantity: 500, locked_quantity: 75, available: 425 },
  { sku_id: "sku-002", sku_code: "MUG-001", sku_name: "陶瓷马克杯-350ml", warehouse: "美国仓A", total_quantity: 300, locked_quantity: 20, available: 280 },
  { sku_id: "sku-003", sku_code: "TSHIRT-002", sku_name: "纯棉T恤-黑色-L", warehouse: "美国仓A", total_quantity: 200, locked_quantity: 0, available: 200 },
]);

const journals = ref([
  { id: "1", sku_code: "TSHIRT-001", change_type: "lock", change_qty: 10, before_total: 500, after_total: 500, before_avail: 500, after_avail: 490, created_at: "2026-05-21 10:01" },
  { id: "2", sku_code: "TSHIRT-001", change_type: "lock", change_qty: 15, before_total: 500, after_total: 500, before_avail: 490, after_avail: 475, created_at: "2026-05-21 10:05" },
  { id: "3", sku_code: "TSHIRT-001", change_type: "deduct", change_qty: 10, before_total: 500, after_total: 490, before_avail: 475, after_avail: 480, created_at: "2026-05-21 14:00" },
  { id: "4", sku_code: "MUG-001", change_type: "lock", change_qty: 20, before_total: 300, after_total: 300, before_avail: 300, after_avail: 280, created_at: "2026-05-21 10:03" },
]);

const changeLabels: Record<string, { label: string; color: string }> = {
  lock: { label: "锁定", color: "#E6A23C" },
  release: { label: "释放", color: "#67C23A" },
  deduct: { label: "扣减", color: "#F56C6C" },
  increase: { label: "入库", color: "#409EFF" },
};
</script>

<template>
  <div>
    <el-card style="margin-bottom: 16px">
      <template #header><span>库存余额</span></template>

      <el-table :data="balances" stripe>
        <el-table-column prop="sku_code" label="SKU编码" width="150" />
        <el-table-column prop="sku_name" label="SKU名称" min-width="180" />
        <el-table-column prop="warehouse" label="仓库" width="120" />
        <el-table-column prop="total_quantity" label="总库存" width="100" align="right">
          <template #default="{ row }">
            <span style="font-weight: bold">{{ row.total_quantity }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="locked_quantity" label="已锁定" width="100" align="right">
          <template #default="{ row }">
            <span :style="{ color: row.locked_quantity > 0 ? '#E6A23C' : '#909399' }">{{ row.locked_quantity }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="available" label="可用库存" width="100" align="right">
          <template #default="{ row }">
            <span :style="{ color: row.available < 20 ? '#F56C6C' : '#67C23A', fontWeight: 'bold' }">{{ row.available }}</span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>库存流水</span>
          <el-select placeholder="变动类型" clearable size="small" style="width: 140px">
            <el-option label="锁定" value="lock" />
            <el-option label="释放" value="release" />
            <el-option label="扣减" value="deduct" />
            <el-option label="入库" value="increase" />
          </el-select>
        </div>
      </template>

      <el-table :data="journals" stripe size="small">
        <el-table-column prop="sku_code" label="SKU" width="140" />
        <el-table-column prop="change_type" label="类型" width="80">
          <template #default="{ row }">
            <span :style="{ color: changeLabels[row.change_type]?.color, fontWeight: 'bold' }">
              {{ changeLabels[row.change_type]?.label || row.change_type }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="change_qty" label="数量" width="80" align="center" />
        <el-table-column prop="before_total" label="变动前总数" width="100" align="right" />
        <el-table-column prop="after_total" label="变动后总数" width="100" align="right" />
        <el-table-column prop="before_avail" label="变动前可用" width="110" align="right" />
        <el-table-column prop="after_avail" label="变动后可用" width="110" align="right" />
        <el-table-column prop="created_at" label="时间" width="170" />
      </el-table>

      <el-pagination style="margin-top: 16px; justify-content: flex-end" background layout="total, prev, pager, next" :total="4" :page-size="20" size="small" />
    </el-card>
  </div>
</template>
