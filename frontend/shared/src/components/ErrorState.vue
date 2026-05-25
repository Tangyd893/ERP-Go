<script setup lang="ts">
defineProps<{
  type?: "403" | "404" | "500" | "empty" | "timeout";
  message?: string;
  retry?: () => void;
}>();

const defaults: Record<string, string> = {
  "403": "您没有权限访问此功能，请联系管理员",
  "404": "未找到对应数据",
  "500": "服务器繁忙，请稍后重试",
  empty: "暂无数据",
  timeout: "请求超时，请检查网络后重试",
};
</script>

<template>
  <div style="padding: 40px 16px; text-align: center">
    <el-result
      :icon="type === 'empty' ? 'info' : type === '403' ? 'warning' : 'error'"
      :title="type === '403' ? '权限不足' : type === '404' ? '未找到' : type === 'timeout' ? '请求超时' : type === 'empty' ? '暂无数据' : '出错了'"
      :sub-title="message || defaults[type || '500'] || defaults['500']"
    >
      <template v-if="retry" #extra>
        <el-button type="primary" @click="retry">重试</el-button>
      </template>
    </el-result>
  </div>
</template>
