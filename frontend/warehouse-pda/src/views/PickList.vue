<script setup lang="ts">
import { computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useWarehouseStore } from "@/stores/warehouse";

const router = useRouter();
const store = useWarehouseStore();

onMounted(() => {
  store.fetchOutbounds();
});

const pickableOutbounds = computed(() =>
  store.outbounds.filter((o) => o.status === "picking" || o.status === "created")
);

function openPick(outboundId: string) {
  router.push({ path: "/pick/scan", query: { outbound_id: outboundId } });
}
</script>

<template>
  <div style="padding: 16px">
    <el-page-header @back="router.push('/')" content="拣货任务" />
    <el-card
      v-for="ob in pickableOutbounds"
      :key="ob.id"
      style="margin-top: 12px; cursor: pointer"
      @click="openPick(ob.id)"
    >
      <div style="font-weight: bold">{{ ob.order_no }}</div>
      <div style="color: #909399; font-size: 13px">{{ ob.id }} · {{ ob.status }}</div>
    </el-card>
    <el-empty v-if="store.outbounds.length === 0" description="暂无拣货任务" />
  </div>
</template>
