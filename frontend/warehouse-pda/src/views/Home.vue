<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { useWarehouseStore } from "@/stores/warehouse";

const router = useRouter();
const store = useWarehouseStore();
const loadError = ref("");

const pickingCount = computed(
  () => store.outbounds.filter((o) => o.status === "picking" || o.status === "created").length
);
const checkingCount = computed(
  () => store.outbounds.filter((o) => o.status === "picked" || o.status === "checking").length
);
const packingCount = computed(
  () => store.outbounds.filter((o) => ["checked"].includes(o.status)).length
);
const weighedCount = computed(
  () => store.outbounds.filter((o) => o.status === "packed").length
);
/** 待出库 = shipped 之前的所有单（含 weighed 但未 confirmShip） */
const pendingShipCount = computed(
  () => store.outbounds.filter((o) => o.status !== "shipped" && o.status !== "cancelled").length
);

onMounted(async () => {
  try {
    await store.fetchOutbounds();
  } catch (e: unknown) {
    const err = e as { response?: { status?: number } };
    if (err?.response?.status === 401) {
      router.push("/login");
      return;
    }
    loadError.value = "加载出库任务失败，请下拉刷新重试";
  }
});
</script>

<template>
  <div class="home-page">
    <h3 class="home-title">WMS 仓库作业</h3>

    <!-- 加载态 -->
    <div v-if="store.loading" class="home-loading">
      <el-skeleton :rows="4" animated />
    </div>

    <!-- 错误态 -->
    <el-empty v-else-if="loadError" :description="loadError">
      <el-button type="primary" @click="loadError = ''; store.fetchOutbounds()">重试</el-button>
    </el-empty>

    <!-- 空态 -->
    <el-empty v-else-if="store.outbounds.length === 0" description="暂无出库任务">
      <el-button type="primary" @click="store.fetchOutbounds()">刷新</el-button>
    </el-empty>

    <!-- 正常数据 -->
    <el-space v-else direction="vertical" class="home-cards" :size="12">
      <el-card shadow="hover" class="pda-touch-target home-card" @click="router.push('/pick')">
        <div class="home-card__label">拣货任务</div>
        <div class="home-card__count home-card__count--pick">{{ pickingCount }}</div>
      </el-card>
      <el-card shadow="hover" class="pda-touch-target home-card" @click="router.push('/check')">
        <div class="home-card__label">复核任务</div>
        <div class="home-card__count home-card__count--check">{{ checkingCount }}</div>
      </el-card>
      <el-card shadow="hover" class="pda-touch-target home-card" @click="router.push('/pack')">
        <div class="home-card__label">打包</div>
        <div class="home-card__count home-card__count--pack">{{ packingCount }}</div>
      </el-card>
      <el-card shadow="hover" class="pda-touch-target home-card" @click="router.push('/weigh')">
        <div class="home-card__label">称重</div>
        <div class="home-card__count home-card__count--weigh">{{ weighedCount }}</div>
      </el-card>
      <el-card shadow="hover" class="pda-touch-target home-card" @click="router.push('/ship')">
        <div class="home-card__label">待出库</div>
        <div class="home-card__count home-card__count--ship">{{ pendingShipCount }}</div>
      </el-card>
    </el-space>
  </div>
</template>

<style scoped>
.home-page {
  /* container */
}

.home-title {
  margin: 0 0 12px;
  font-size: 18px;
}

.home-loading {
  text-align: center;
  padding: 48px 0;
}

.home-cards {
  width: 100%;
}

.home-card {
  cursor: pointer;
}

.home-card__label {
  text-align: center;
  font-size: 16px;
}

.home-card__count {
  text-align: center;
  font-size: 28px;
  font-weight: bold;
}

.home-card__count--pick  { color: var(--pda-primary); }
.home-card__count--check { color: var(--pda-warning); }
.home-card__count--pack  { color: var(--pda-text-secondary); }
.home-card__count--weigh { color: var(--pda-danger); }
.home-card__count--ship  { color: var(--pda-success); }
</style>
