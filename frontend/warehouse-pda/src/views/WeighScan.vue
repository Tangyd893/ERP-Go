<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { useWarehouseStore } from "@/stores/warehouse";

const router = useRouter();
const store = useWarehouseStore();
const weighWeight = ref<Record<string, number>>({});
const scanning = ref(false);

onMounted(() => store.fetchOutbounds());

const weighedOutbounds = computed(() =>
  store.outbounds.filter((o) => o.status === "packed")
);

async function handleWeigh(outboundId: string) {
  if (scanning.value) return;
  const weight = weighWeight.value[outboundId];
  if (!weight || weight <= 0) {
    ElMessage.warning("请输入有效重量");
    return;
  }
  scanning.value = true;
  try {
    await store.weigh(outboundId, weight);
    ElMessage.success("称重完成");
    weighWeight.value[outboundId] = 0;
    await store.fetchOutbounds();
  } catch (e: unknown) {
    const msg = (e as { response?: { status?: number } })?.response?.status === 409
      ? "该出库单已称重，请刷新"
      : "称重失败，请重试";
    ElMessage.error(msg);
  } finally {
    scanning.value = false;
  }
}
</script>

<template>
  <div style="padding: 16px">
    <el-page-header @back="router.push('/')" content="称重" />
    <el-empty v-if="weighedOutbounds.length === 0" description="暂无待称重出库单" />
    <el-card
      v-for="ob in weighedOutbounds"
      :key="ob.id"
      style="margin-top: 12px"
    >
      <div style="font-weight: bold">{{ ob.order_no }}</div>
      <div style="color: #909399; margin: 8px 0; font-size: 13px">{{ ob.id }}</div>
      <el-input-number
        v-model="weighWeight[ob.id]"
        :min="1"
        :step="1"
        placeholder="重量(g)"
        style="margin-right: 8px"
      />
      <el-button type="warning" :loading="scanning" @click="handleWeigh(ob.id)">
        确认称重
      </el-button>
    </el-card>
  </div>
</template>
