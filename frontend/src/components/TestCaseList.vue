<template>
  <t-table :data="testCases" :columns="columns" hover>
    <template #status="{ row }">
      <t-tag :theme="getStatusTheme(row.status)">
        {{ getStatusText(row.status) }}
      </t-tag>
    </template>
    <template #error_log="{ row }">
      <t-popup v-if="row.error_log" trigger="hover" :content="row.error_log">
        <t-button size="small" variant="text">查看错误</t-button>
      </t-popup>
      <span v-else>-</span>
    </template>
  </t-table>
</template>

<script setup>
defineProps({
  testCases: {
    type: Array,
    default: () => [],
  },
})

const columns = [
  { colKey: 'name', title: '测例名称', width: 300 },
  { colKey: 'status', title: '状态', width: 100 },
  { colKey: 'duration_ms', title: '耗时(ms)', width: 120 },
  { colKey: 'error_log', title: '错误日志', width: 150 },
]

const getStatusTheme = (status) => {
  const themes = {
    passed: 'success',
    failed: 'danger',
    skipped: 'default',
  }
  return themes[status] || 'default'
}

const getStatusText = (status) => {
  const texts = {
    passed: '通过',
    failed: '失败',
    skipped: '跳过',
  }
  return texts[status] || status
}
</script>

