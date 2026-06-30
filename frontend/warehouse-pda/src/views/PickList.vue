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
  <div class="pick-list">
    <el-page-header @back="router.push('/')" content="拣货任务" />

    <el-empty
      v-if="store.outbounds.length === 0 && !store.loading"
      description="暂无拣货任务"
    />

    <div v-if="store.loading" class="pick-list__loading">
      <el-skeleton :rows="3" animated />
    </div>

    <el-card
      v-for="ob in pickableOutbounds"
      :key="ob.id"
      class="pick-list__card"
      shadow="hover"
      @click="openPick(ob.id)"
    >
      <div class="pick-list__order-no">{{ ob.order_no }}</div>
      <div class="pick-list__meta">{{ ob.id }} · {{ ob.status }}</div>
    </el-card>
  </div>
</template>

<style scoped>
.pick-list__loading {
  padding: 16px 0;
}

.pick-list__card {
  margin-top: 12px;
  cursor: pointer;
}

.pick-list__order-no {
  font-weight: 600;
  font-size: 15px;
}

.pick-list__meta {
  color: var(--pda-text-secondary, #909399);
  font-size: 13px;
  margin-top: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
