<script setup lang="ts">
import { computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useWarehouseStore } from "@/stores/warehouse";

const router = useRouter();
const store = useWarehouseStore();

const pickingCount = computed(
  () => store.outbounds.filter((o) => o.status === "picking" || o.status === "created").length
);
const checkingCount = computed(
  () => store.outbounds.filter((o) => o.status === "picked" || o.status === "checking").length
);
const packingCount = computed(
  () => store.outbounds.filter((o) => ["checked", "packed", "weighed"].includes(o.status)).length
);

onMounted(() => {
  store.fetchOutbounds();
});
</script>

<template>
  <div style="padding: 16px">
    <h3 style="margin-bottom: 16px">WMS 仓库作业</h3>
    <el-space direction="vertical" style="width: 100%" :size="12">
      <el-card shadow="hover" style="cursor: pointer" @click="router.push('/pick')">
        <div style="text-align: center; font-size: 18px">拣货任务</div>
        <div style="text-align: center; color: #409EFF; font-size: 24px; font-weight: bold">{{ pickingCount }}</div>
      </el-card>
      <el-card shadow="hover" style="cursor: pointer" @click="router.push('/check')">
        <div style="text-align: center; font-size: 18px">复核任务</div>
        <div style="text-align: center; color: #E6A23C; font-size: 24px; font-weight: bold">{{ checkingCount }}</div>
      </el-card>
      <el-card shadow="hover" style="cursor: pointer" @click="router.push('/ship')">
        <div style="text-align: center; font-size: 18px">待出库</div>
        <div style="text-align: center; color: #67C23A; font-size: 24px; font-weight: bold">{{ packingCount }}</div>
      </el-card>
    </el-space>
  </div>
</template>
