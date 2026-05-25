<script setup lang="ts">
import { computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { useWarehouseStore, getErrorMessage } from "@/stores/warehouse";

const router = useRouter();
const store = useWarehouseStore();

onMounted(() => store.fetchOutbounds());

const checkingOutbounds = computed(() =>
  store.outbounds.filter((o) => o.status === "picked" || o.status === "picking")
);

async function handleCheck(ob: { id: string }) {
  try {
    await store.checkScan(ob.id, "", 1);
    ElMessage.success("复核完成");
    await store.fetchOutbounds();
  } catch (e: unknown) {
    ElMessage.error(getErrorMessage(e, "复核失败"));
  }
}
</script>

<template>
  <div style="padding: 16px">
    <el-page-header @back="router.push('/')" content="复核任务" />
    <el-card
      v-for="ob in checkingOutbounds"
      :key="ob.id"
      style="margin-top: 12px"
    >
      <div style="font-weight: bold">{{ ob.order_no }}</div>
      <el-button type="warning" style="margin-top: 8px" @click="handleCheck(ob)">复核确认</el-button>
    </el-card>
  </div>
</template>
