<script setup lang="ts">
import { ref } from "vue";

const orders = ref([
  { id: "1", platform_order_no: "AMZ-20260521-001", store: "美国站店铺A", buyer_name: "John Doe", currency: "USD", total_amount: 46.97, status: "locked", items: 2, created_at: "2026-05-21 10:00" },
  { id: "2", platform_order_no: "AMZ-20260521-002", store: "美国站店铺A", buyer_name: "Jane Smith", currency: "USD", total_amount: 15.99, status: "approved", items: 1, created_at: "2026-05-21 10:15" },
  { id: "3", platform_order_no: "AMZ-20260521-003", store: "英国站店铺B", buyer_name: "Tom Brown", currency: "GBP", total_amount: 25.98, status: "pending", items: 2, created_at: "2026-05-21 11:00" },
  { id: "4", platform_order_no: "AMZ-20260520-004", store: "美国站店铺A", buyer_name: "Alice Wang", currency: "USD", total_amount: 12.99, status: "shipped", items: 1, created_at: "2026-05-20 16:00" },
]);

const statusLabels: Record<string, { type: string; label: string }> = {
  pending: { type: "info", label: "待审核" },
  approved: { type: "warning", label: "已审核" },
  locked: { type: "", label: "已锁定" },
  picking: { type: "", label: "拣货中" },
  packed: { type: "", label: "已打包" },
  shipped: { type: "success", label: "已发货" },
  abnormal: { type: "danger", label: "异常" },
  cancelled: { type: "info", label: "已取消" },
};
</script>

<template>
  <div>
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>订单列表</span>
          <el-button type="primary">导入订单</el-button>
        </div>
      </template>

      <el-form inline style="margin-bottom: 16px">
        <el-form-item label="平台订单号">
          <el-input placeholder="搜索订单号" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select placeholder="全部" clearable style="width: 140px">
            <el-option v-for="(v, k) in statusLabels" :key="k" :label="v.label" :value="k" />
          </el-select>
        </el-form-item>
        <el-form-item label="店铺">
          <el-select placeholder="全部" clearable style="width: 160px">
            <el-option label="美国站店铺A" value="store-1" />
            <el-option label="英国站店铺B" value="store-2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary">查询</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="orders" stripe>
        <el-table-column prop="platform_order_no" label="平台订单号" width="200" />
        <el-table-column prop="store" label="店铺" width="150" />
        <el-table-column prop="buyer_name" label="买家" width="130" />
        <el-table-column prop="total_amount" label="金额" width="100" align="right">
          <template #default="{ row }">{{ row.currency }} {{ row.total_amount.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="items" label="件数" width="70" align="center" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusLabels[row.status]?.type || 'info'" size="small">
              {{ statusLabels[row.status]?.label || row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170" />
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" :disabled="row.status !== 'pending'">审核</el-button>
            <el-button type="warning" size="small" :disabled="row.status === 'shipped' || row.status === 'cancelled'">异常</el-button>
            <el-button type="danger" size="small" :disabled="row.status === 'cancelled'">取消</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination style="margin-top: 16px; justify-content: flex-end" background layout="total, prev, pager, next" :total="4" :page-size="20" />
    </el-card>
  </div>
</template>
