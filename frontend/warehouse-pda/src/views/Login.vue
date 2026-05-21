<script setup lang="ts">
</script>

<template>
  <div style="padding: 20px; max-width: 400px; margin: 0 auto">
    <h2 style="text-align: center; margin-bottom: 24px">WMS PDA</h2>
    <el-form>
      <el-form-item>
        <el-input placeholder="工号" size="large" />
      </el-form-item>
      <el-form-item>
        <el-input type="password" placeholder="密码" size="large" show-password />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" size="large" style="width: 100%">登 录</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>
