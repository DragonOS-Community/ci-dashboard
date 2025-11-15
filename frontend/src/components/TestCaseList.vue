<template>
  <t-table
    :data="testCases"
    :columns="columns"
    hover
    row-key="id"
    :sort="sortInfo"
    @sort-change="handleSortChange"
  >
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
import { ref, watch } from "vue";

const props = defineProps({
  testCases: {
    type: Array,
    default: () => [],
  },
  sortField: {
    type: String,
    default: "",
  },
  sortOrder: {
    type: String,
    default: "",
  },
});

const emit = defineEmits(["sort-change"]);

const sortInfo = ref([]);

// 同步外部排序状态
watch(
  () => [props.sortField, props.sortOrder],
  ([field, order]) => {
    if (field && order) {
      sortInfo.value = [
        {
          sortBy: field,
          descending: order === "desc",
        },
      ];
    } else {
      sortInfo.value = [];
    }
  },
  { immediate: true },
);

const columns = [
  {
    colKey: "name",
    title: "测例名称",
    width: 300,
    sortType: "all",
    sorter: true,
  },
  { colKey: "status", title: "状态", width: 100 },
  {
    colKey: "duration_ms",
    title: "耗时(ms)",
    width: 120,
    sortType: "all",
    sorter: true,
  },
  { colKey: "error_log", title: "错误日志", width: 150 },
];

const handleSortChange = (sortInfo) => {
  emit("sort-change", sortInfo);
};

const getStatusTheme = (status) => {
  const themes = {
    passed: "success",
    failed: "danger",
    skipped: "default",
  };
  return themes[status] || "default";
};

const getStatusText = (status) => {
  const texts = {
    passed: "通过",
    failed: "失败",
    skipped: "跳过",
  };
  return texts[status] || status;
};
</script>
