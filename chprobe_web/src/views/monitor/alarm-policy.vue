<template>
  <div>
    <a-card class="alarm-policy-card">



      <a-table
        :columns="columns"
        :data-source="policies"
        :pagination="pagination"
        :loading="loading"
        row-key="id"
        class="table-sm"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'level'">
            <a-tag :color="getLevelColor(record.level)">
              {{ getLevelText(record.level) }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'status'">
            <a-switch v-model:checked="record.enabled" @change="handleStatusChange(record)" />
          </template>
          <template v-else-if="column.key === 'condition'">
            <div>{{ record.metric }} {{ record.operator }} {{ record.threshold }} {{ record.unit }}</div>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-popover
              placement="top"
              trigger="hover"
              :title="null"
            >
              <template #content>
                <div class="action-popover">
                  <div class="action-item" @click="handleEdit(record)">
                    编辑
                  </div>
                  <div class="action-item" @click="handleTest(record)">
                    测试
                  </div>
                  <div class="action-item action-danger" @click="handleDelete(record)">
                    删除
                  </div>
                </div>
              </template>
              <a-button type="text" size="small" class="btn-action">
                操作
              </a-button>
            </a-popover>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, h } from 'vue'
import { message } from 'ant-design-vue'
import {
  PlusOutlined,
  AlertOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  BellOutlined,
  EditOutlined,
  DeleteOutlined,
  PlayCircleOutlined,
  ReloadOutlined,
  DownOutlined
} from '@ant-design/icons-vue'
import type { MenuProps } from 'ant-design-vue'

interface Policy {
  id: string
  name: string
  metric: string
  operator: string
  threshold: number
  unit: string
  duration: number
  level: 'critical' | 'warning' | 'info'
  targets: string[]
  enabled: boolean
  created: string
  updated: string
}

const columns = [
  { title: '策略名称', dataIndex: 'name', key: 'name' },
  { title: '监控指标', dataIndex: 'metric', key: 'metric', width: 150 },
  { title: '告警条件', key: 'condition', width: 200 },
  { title: '持续时间', dataIndex: 'duration', key: 'duration', width: 100, render: (text: number) => `${text}s` },
  { title: '告警级别', dataIndex: 'level', key: 'level', width: 100 },
  { title: '目标主机', dataIndex: 'targets', key: 'targets', width: 200, render: (targets: string[]) => targets.join(', ') },
  { title: '状态', dataIndex: 'status', key: 'status', width: 80 },
  { title: '更新时间', dataIndex: 'updated', key: 'updated', width: 180 },
  { title: '操作', key: 'action', width: 200, fixed: 'right' }
]

const policies = ref<Policy[]>([
  {
    id: '1',
    name: 'CPU 使用率过高',
    metric: 'cpu_usage',
    operator: '>',
    threshold: 80,
    unit: '%',
    duration: 60,
    level: 'warning',
    targets: ['server-01', 'server-02', 'server-03'],
    enabled: true,
    created: '2024-01-10 10:00:00',
    updated: '2024-01-15 14:30:00'
  },
  {
    id: '2',
    name: '内存使用率过高',
    metric: 'memory_usage',
    operator: '>',
    threshold: 85,
    unit: '%',
    duration: 60,
    level: 'warning',
    targets: ['server-01', 'server-02', 'server-03'],
    enabled: true,
    created: '2024-01-11 11:00:00',
    updated: '2024-01-14 09:15:00'
  },
  {
    id: '3',
    name: '磁盘空间不足',
    metric: 'disk_usage',
    operator: '>',
    threshold: 90,
    unit: '%',
    duration: 300,
    level: 'critical',
    targets: ['server-01', 'server-02', 'server-03'],
    enabled: true,
    created: '2024-01-12 08:00:00',
    updated: '2024-01-16 10:00:00'
  },
  {
    id: '4',
    name: '网络流量异常',
    metric: 'network_traffic',
    operator: '>',
    threshold: 100,
    unit: 'MB/s',
    duration: 120,
    level: 'info',
    targets: ['server-01'],
    enabled: false,
    created: '2024-01-13 14:00:00',
    updated: '2024-01-13 14:00:00'
  },
  {
    id: '5',
    name: '容器重启频繁',
    metric: 'container_restarts',
    operator: '>',
    threshold: 5,
    unit: '次/小时',
    duration: 3600,
    level: 'warning',
    targets: ['server-02', 'server-03'],
    enabled: true,
    created: '2024-01-14 09:00:00',
    updated: '2024-01-15 16:45:00'
  },
  {
    id: '6',
    name: '进程不存在',
    metric: 'process_exists',
    operator: '==',
    threshold: 0,
    unit: '',
    duration: 60,
    level: 'critical',
    targets: ['server-01', 'server-02'],
    enabled: true,
    created: '2024-01-15 10:00:00',
    updated: '2024-01-15 10:00:00'
  }
])

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 6
})

const loading = ref(false)

const enabledCount = computed(() => policies.value.filter(p => p.enabled).length)

const criticalCount = computed(() => policies.value.filter(p => p.level === 'critical').length)

const todayAlerts = ref(28)

const getLevelColor = (level: string) => {
  if (level === 'critical') return 'red'
  if (level === 'warning') return 'orange'
  return 'blue'
}

const getLevelText = (level: string) => {
  if (level === 'critical') return '紧急'
  if (level === 'warning') return '警告'
  return '提示'
}

const handleCreate = () => {
  message.info('打开创建策略对话框')
}
const handleRefresh = () => {
 loading.value = true
 setTimeout(() => {
 loading.value = false
 message.success('刷新成功')
 }, 1000)
}
const getActionMenu = (record: Policy) => {
 return {
 items: [
 {
 key: 'edit',
 label: '编辑',
 onClick: () => handleEdit(record)
 },
 {
 key: 'test',
 label: '测试',
 onClick: () => handleTest(record)
 },
 {
 key: 'delete',
 label: '删除',
 danger: true,
 onClick: () => handleDelete(record)
 }
 ]
 }
}

const handleEdit = (record: Policy) => {
  message.info(`编辑策略: ${record.name}`)
}

const handleTest = (record: Policy) => {
  message.success(`策略 ${record.name} 测试成功`)
}

const handleDelete = (record: Policy) => {
  message.success(`策略 ${record.name} 已删除`)
}

const handleStatusChange = (record: Policy) => {
  message.success(`策略 ${record.name} 已${record.enabled ? '启用' : '禁用'}`)
}
</script>

<style scoped>
.alarm-policy-card {
  overflow: visible;
}

.alarm-policy-card :deep(.ant-card-head) {
  min-height: 32px !important;
  height: 32px !important;
  line-height: 32px !important;
  padding: 1px !important;
  margin: 0 !important;
  width: 100% !important;
  border: none !important;
  background: transparent !important;
  box-shadow: none !important;
}

.alarm-policy-card :deep(.ant-card-head-title) {
  font-size: 12px !important;
  line-height: 32px !important;
  padding: 1px !important;
  margin: 0 !important;
  height: 32px !important;
}

.alarm-policy-card :deep(.ant-card-body) {
  padding: 1px !important;
  margin: 0 !important;
}

.alarm-policy-card :deep(.ant-card-extra) {
  height: 32px !important;
  line-height: 32px !important;
  padding: 0 !important;
  margin: 0 !important;
}

.btn-sm {
  height: 24px !important;
  padding: 0 8px !important;
  font-size: 12px !important;
  line-height: 24px !important;
}

.table-sm :deep(.ant-table-thead > tr > th) {
  height: 32px !important;
  line-height: 32px !important;
  padding: 1px !important;
  font-size: 14px !important;
  text-align: center !important;
}

.table-sm :deep(.ant-table-tbody > tr > td) {
  height: 32px !important;
  line-height: 32px !important;
  padding: 1px !important;
  font-size: 14px !important;
  text-align: center !important;
}

.btn-action {
  height: 24px !important;
  padding: 0 4px !important;
  font-size: 12px !important;
  line-height: 24px !important;
  color: #1890ff !important;
}

.action-popover {
  padding: 2px 0;
  width: 70px;
}

.action-item {
  padding: 3px 4px;
  font-size: 14px;
  color: #1890ff;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
}

.action-item:hover {
  background-color: #e6f7ff;
  color: #1890ff;
}

.action-danger {
  color: #f5222d;
}

.action-danger:hover {
  background-color: #fff1f0;
  color: #f5222d;
}
</style>