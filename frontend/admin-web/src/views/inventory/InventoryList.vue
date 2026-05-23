<script setup lang="ts">
import { computed, onMounted } from "vue";
import { ProTable } from "@erp/shared";
import { useInventoryStore } from "@/stores/inventory";

const inventoryStore = useInventoryStore();

const displayBalances = computed(() => inventoryStore.balances);
const displayJournals = computed(() => inventoryStore.journals);

const balanceColumns = [
  { prop: "sku_code", label: "SKU编码", width: 150 },
  { prop: "sku_name", label: "SKU名称", minWidth: 180 },
  { prop: "warehouse_name", label: "仓库", width: 120 },
  { prop: "qty", label: "总库存", width: 100, align: "right" as const },
  { prop: "locked_qty", label: "已锁定", width: 100, align: "right" as const },
  { prop: "available_qty", label: "可用库存", width: 100, align: "right" as const },
];

const journalColumns = [
  { prop: "sku_code", label: "SKU", width: 140 },
  { prop: "change_type", label: "类型", width: 80 },
  { prop: "change_qty", label: "数量", width: 80, align: "center" as const },
  { prop: "before_total", label: "变动前总数", width: 100, align: "right" as const },
  { prop: "after_total", label: "变动后总数", width: 100, align: "right" as const },
  { prop: "before_avail", label: "变动前可用", width: 110, align: "right" as const },
  { prop: "after_avail", label: "变动后可用", width: 110, align: "right" as const },
  { prop: "created_at", label: "时间", width: 170 },
];

const changeLabels: Record<string, { label: string; color: string }> = {
  lock: { label: "锁定", color: "#E6A23C" },
  release: { label: "释放", color: "#67C23A" },
  deduct: { label: "扣减", color: "#F56C6C" },
  increase: { label: "入库", color: "#409EFF" },
};

onMounted(() => {
  inventoryStore.fetchBalances();
  inventoryStore.fetchJournals();
});
</script>

<template>
  <div>
    <el-card style="margin-bottom: 16px">
      <template #header><span>库存余额</span></template>

      <ProTable
        :columns="balanceColumns"
        :data="displayBalances"
        :loading="inventoryStore.loading"
        :total="displayBalances.length"
      >
        <template #qty="{ row }">
          <span style="font-weight: bold">{{ row.qty }}</span>
        </template>
        <template #locked_qty="{ row }">
          <span :style="{ color: row.locked_qty > 0 ? '#E6A23C' : '#909399' }">{{ row.locked_qty }}</span>
        </template>
        <template #available_qty="{ row }">
          <span :style="{ color: row.available_qty < 20 ? '#F56C6C' : '#67C23A', fontWeight: 'bold' }">{{ row.available_qty }}</span>
        </template>
      </ProTable>
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

      <ProTable
        :columns="journalColumns"
        :data="displayJournals"
        :loading="inventoryStore.loading"
        :total="displayJournals.length"
      >
        <template #change_type="{ row }">
          <span :style="{ color: changeLabels[row.change_type]?.color, fontWeight: 'bold' }">
            {{ changeLabels[row.change_type]?.label || row.change_type }}
          </span>
        </template>
      </ProTable>
    </el-card>
  </div>
</template>
