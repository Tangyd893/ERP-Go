<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const authStore = useAuthStore();

const form = ref({ username: "admin", password: "admin123" });
const loading = ref(false);

async function handleLogin() {
  if (!form.value.username || !form.value.password) {
    ElMessage.warning("请输入工号和密码");
    return;
  }
  loading.value = true;
  try {
    await authStore.login(form.value.username, form.value.password, "default");
    ElMessage.success("登录成功");
    router.push("/");
  } catch {
    ElMessage.error("登录失败，请检查账号密码或服务是否启动");
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div style="padding: 20px; max-width: 400px; margin: 0 auto">
    <h2 style="text-align: center; margin-bottom: 24px">WMS PDA</h2>
    <el-form @keyup.enter="handleLogin">
      <el-form-item>
        <el-input v-model="form.username" placeholder="工号" size="large" />
      </el-form-item>
      <el-form-item>
        <el-input v-model="form.password" type="password" placeholder="密码" size="large" show-password />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" size="large" style="width: 100%" :loading="loading" @click="handleLogin">
          登 录
        </el-button>
      </el-form-item>
    </el-form>
    <p style="color: #909399; font-size: 12px; text-align: center">开发默认: admin / admin123</p>
  </div>
</template>
