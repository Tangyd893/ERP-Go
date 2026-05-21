<script setup lang="ts">
</script>

<template>
  <div style="display: flex; justify-content: center; align-items: center; height: 100vh; background-color: #f5f7fa">
    <el-card style="width: 400px">
      <h2 style="text-align: center; margin-bottom: 24px">ERP-Go 管理后台</h2>
      <el-form>
        <el-form-item>
          <el-input placeholder="用户名" prefix-icon="User" />
        </el-form-item>
        <el-form-item>
          <el-input type="password" placeholder="密码" prefix-icon="Lock" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" style="width: 100%">登 录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>
