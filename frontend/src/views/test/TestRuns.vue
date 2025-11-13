<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">测试运行管理</h1>
      <p class="page-description">管理和查看所有测试运行任务</p>
    </div>

    <div class="content-card">
      <t-card>
        <t-alert message="此页面将整合原有的测试列表功能，提供更强大的管理功能" theme="info" />
      </t-card>
    </div>
  </div>
</template>

<script setup>
// 占位页面
</script>

<style scoped>
.page-container {
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 32px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
}

.page-description {
  font-size: 14px;
  color: #6b7280;
}

.content-card {
  background-color: #ffffff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}
</style>