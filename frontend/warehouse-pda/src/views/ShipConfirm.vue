<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { useWarehouseStore, getErrorMessage } from "@/stores/warehouse";

const router = useRouter();
const store = useWarehouseStore();
const trackingNo = ref("");

onMounted(() => store.fetchOutbounds());

const shippableOutbounds = computed(() =>
  store.outbounds.filter((o) =>
    ["checked", "packed", "weighed", "picked"].includes(o.status)
  )
);

async function handleShip(outboundId: string) {
  try {
    await store.pack(outboundId);
    await store.weigh(outboundId, 1);
    await store.confirmShip(outboundId, trackingNo.value, "default");
    ElMessage.success("出库确认成功");
    await store.fetchOutbounds();
  } catch (e: unknown) {
    ElMessage.error(getErrorMessage(e, "出库确认失败"));
  }
}
</script>

<template>
  <div class="ship-confirm">
    <el-page-header @back="router.push('/')" content="出库确认" />

    <el-input
      v-model="trackingNo"
      placeholder="物流单号（可选）"
      class="ship-confirm__tracking"
    />

    <el-empty
      v-if="shippableOutbounds.length === 0 && !store.loading"
      description="暂无待出库单"
    />

    <div v-if="store.loading" class="ship-confirm__loading">
      <el-skeleton :rows="3" animated />
    </div>

    <el-card
      v-for="ob in shippableOutbounds"
      :key="ob.id"
      class="ship-confirm__card"
    >
      <div class="ship-confirm__order-no">{{ ob.order_no }}</div>
      <div class="ship-confirm__meta">{{ ob.id }} · {{ ob.status }}</div>
      <el-button
        type="success"
        class="ship-confirm__btn"
        @click="handleShip(ob.id)"
      >
        确认出库
      </el-button>
    </el-card>
  </div>
</template>

<style scoped>
.ship-confirm__tracking {
  margin: 12px 0;
}

.ship-confirm__loading {
  padding: 16px 0;
}

.ship-confirm__card {
  margin-top: 12px;
}

.ship-confirm__order-no {
  font-weight: 600;
  font-size: 15px;
}

.ship-confirm__meta {
  color: var(--pda-text-secondary, #909399);
  font-size: 13px;
  margin: 4px 0 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ship-confirm__btn {
  margin-top: 8px;
  min-height: var(--pda-touch-min, 44px);
}
</style>
