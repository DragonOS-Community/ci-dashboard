<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">测试概览</h1>
      <p class="page-description">查看所有测试项目的运行状态和统计信息</p>
    </div>

    <div class="coming-soon">
      <t-icon name="chart-line" size="48" />
      <h2>功能开发中</h2>
      <p>测试概览页面正在开发中，敬请期待...</p>
    </div>
  </div>
</template>

<script setup>
// 占位页面，后续实现具体功能
</script>

<style scoped>
.page-container {
  max-width: 1200px;
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

.coming-soon {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  text-align: center;
  background-color: #ffffff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.coming-soon .t-icon {
  color: #f59e0b;
  margin-bottom: 16px;
}

.coming-soon h2 {
  font-size: 20px;
  color: #1f2937;
  margin-bottom: 8px;
}

.coming-soon p {
  font-size: 14px;
  color: #6b7280;
}
</style>
