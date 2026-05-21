<script setup lang="ts">
import { computed } from "vue";

interface ColumnDef {
  prop: string;
  label: string;
  width?: string | number;
  minWidth?: string | number;
  fixed?: boolean | "left" | "right";
  sortable?: boolean;
  align?: "left" | "center" | "right";
  slot?: string;
}

interface Props {
  columns: ColumnDef[];
  data: Record<string, unknown>[];
  loading?: boolean;
  total?: number;
  pageSize?: number;
  page?: number;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  total: 0,
  pageSize: 20,
  page: 1,
});

const emit = defineEmits<{
  (e: "page-change", page: number): void;
  (e: "search"): void;
  (e: "refresh"): void;
  (e: "size-change", size: number): void;
}>();

const totalPages = computed(() => Math.ceil(props.total / props.pageSize));
</script>

<template>
  <div class="pro-table">
    <div class="pro-table__toolbar" v-if="$slots.toolbar">
      <slot name="toolbar" />
    </div>

    <el-table
      :data="data"
      v-loading="loading"
      stripe
      border
      style="width: 100%"
    >
      <el-table-column
        v-for="col in columns"
        :key="col.prop"
        :prop="col.prop"
        :label="col.label"
        :width="col.width"
        :min-width="col.minWidth"
        :fixed="col.fixed"
        :sortable="col.sortable"
        :align="col.align"
      >
        <template v-if="col.slot || $slots[col.prop]" #default="{ row }">
          <slot :name="col.slot || col.prop" :row="row" />
        </template>
      </el-table-column>
      <el-table-column v-if="$slots.default" label="操作" :min-width="120">
        <template #default="{ row }">
          <slot :row="row" />
        </template>
      </el-table-column>
    </el-table>

    <div class="pro-table__pagination" v-if="total > 0">
      <el-pagination
        background
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        :page-size="pageSize"
        :current-page="page"
        :page-sizes="[10, 20, 50, 100]"
        @current-change="(p: number) => emit('page-change', p)"
        @size-change="(s: number) => emit('size-change', s)"
      />
    </div>
  </div>
</template>

<style scoped>
.pro-table__pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.pro-table__toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
</style>
