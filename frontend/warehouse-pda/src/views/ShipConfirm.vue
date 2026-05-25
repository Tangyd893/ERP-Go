<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { useWarehouseStore } from "@/stores/warehouse";

const router = useRouter();
const store = useWarehouseStore();
const trackingNo = ref("");

onMounted(() => store.fetchOutbounds());

async function handleShip(outboundId: string) {
  try {
    await store.pack(outboundId);
    await store.weigh(outboundId, 1);
    await store.confirmShip(outboundId, trackingNo.value, "default");
    ElMessage.success("出库确认成功");
    await store.fetchOutbounds();
  } catch {
    ElMessage.error("出库确认失败");
  }
}
</script>

<template>
  <div style="padding: 16px">
    <el-page-header @back="router.push('/')" content="出库确认" />
    <el-input v-model="trackingNo" placeholder="物流单号（可选）" style="margin: 12px 0" />
    <el-card
      v-for="ob in store.outbounds.filter((o) => ['checked', 'packed', 'weighed', 'picked'].includes(o.status))"
      :key="ob.id"
      style="margin-top: 12px"
    >
      <div style="font-weight: bold">{{ ob.order_no }}</div>
      <el-button type="success" style="margin-top: 8px" @click="handleShip(ob.id)">确认出库</el-button>
    </el-card>
  </div>
</template>
