<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { useWarehouseStore } from "@/stores/warehouse";

const router = useRouter();
const store = useWarehouseStore();
const packingWeight = ref<Record<string, number>>({});
const scanning = ref(false);

onMounted(() => store.fetchOutbounds());

const packedOutbounds = computed(() =>
  store.outbounds.filter((o) => o.status === "checked")
);

async function handlePack(outboundId: string) {
  if (scanning.value) return;
  scanning.value = true;
  try {
    const weight = packingWeight.value[outboundId] || 0;
    await store.pack(outboundId, weight);
    ElMessage.success("打包完成");
    packingWeight.value[outboundId] = 0;
    await store.fetchOutbounds();
  } catch (e: unknown) {
    const msg = (e as { response?: { status?: number } })?.response?.status === 409
      ? "该出库单已打包，请刷新"
      : "打包失败，请重试";
    ElMessage.error(msg);
  } finally {
    scanning.value = false;
  }
}
</script>

<template>
  <div style="padding: 16px">
    <el-page-header @back="router.push('/')" content="打包" />
    <el-empty v-if="packedOutbounds.length === 0" description="暂无待打包出库单" />
    <el-card
      v-for="ob in packedOutbounds"
      :key="ob.id"
      style="margin-top: 12px"
    >
      <div style="font-weight: bold">{{ ob.order_no }}</div>
      <div style="color: #909399; margin: 8px 0; font-size: 13px">{{ ob.id }}</div>
      <el-input-number
        v-model="packingWeight[ob.id]"
        :min="0"
        :step="1"
        placeholder="重量(g)"
        style="margin-right: 8px"
      />
      <el-button type="primary" :loading="scanning" @click="handlePack(ob.id)">
        确认打包
      </el-button>
    </el-card>
  </div>
</template>
