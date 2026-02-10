<template>
  <div>
    <a-card>
      <template #title>
        告警策略
      </template>
      <template #extra>
        <a-space>
          <a-button @click="handleRefresh">
            <template #icon>
              <ReloadOutlined />
            </template>
            刷新
          </a-button>
          <a-button type="primary" @click="handleCreate">
            <template #icon>
              <PlusOutlined />
            </template>
            新建策略
          </a-button>
        </a-space>
      </template>

      <a-form :model="searchForm" layout="inline" style="margin-bottom: 16px;">
        <a-form-item label="策略名">
          <a-input
            v-model:value="searchForm.name"
            placeholder="请输入策略名"
            @pressEnter="handleSearch"
          />
        </a-form-item>
        <a-form-item label="状态">
          <a-select
            v-model:value="searchForm.status"
            placeholder="请选择状态"
            style="width: 120px;"
          >
            <a-select-option value="">全部</a-select-option>
            <a-select-option value="enabled">启用</a-select-option>
            <a-select-option value="disabled">禁用</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="监控类型">
          <a-select
            v-model:value="searchForm.type"
            placeholder="请选择类型"
            style="width: 150px;"
          >
            <a-select-option value="">全部类型</a-select-option>
            <a-select-option value="cpu">CPU</a-select-option>
            <a-select-option value="memory">内存</a-select-option>
            <a-select-option value="disk">磁盘</a-select-option>
            <a-select-option value="network">网络</a-select-option>
            <a-select-option value="service">服务</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">
            <template #icon>
              <SearchOutlined />
            </template>
            搜索
          </a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>

      <a-table
        :columns="columns"
        :data-source="policies"
        :pagination="pagination"
        :loading="loading"
        row-key="id"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <a @click="handleView(record)">{{ record.name }}</a>
          </template>
          <template v-if="column.key === 'status'">
            <a-switch
              :checked="record.status === 'enabled'"
              @change="(checked: boolean) => handleToggleStatus(record, checked)"
            />
          </template>
          <template v-if="column.key === 'severity'">
            <a-tag :color="getSeverityColor(record.severity)">
              {{ getSeverityText(record.severity) }}
            </a-tag>
          </template>
          <template v-if="column.key === 'action'">
            <a-dropdown :menu-items="getActionItems(record)" @menu-item-click="({ key }: { key: string }) => handleAction(key, record)">
              <a-button type="text" size="small">
                操作 <DownOutlined />
              </a-button>
            </a-dropdown>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal
      v-model:open="detailVisible"
      title="策略详情"
      width="700px"
      :footer="null"
    >
      <a-descriptions :column="2" bordered v-if="currentPolicy">
        <a-descriptions-item label="策略名">{{ currentPolicy.name }}</a-descriptions-item>
        <a-descriptions-item label="监控类型">{{ getTypeText(currentPolicy.type) }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="currentPolicy.status === 'enabled' ? 'success' : 'default'">
            {{ currentPolicy.status === 'enabled' ? '启用' : '禁用' }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="严重级别">
          <a-tag :color="getSeverityColor(currentPolicy.severity)">
            {{ getSeverityText(currentPolicy.severity) }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="监控目标">{{ currentPolicy.target }}</a-descriptions-item>
        <a-descriptions-item label="指标">{{ currentPolicy.metric }}</a-descriptions-item>
        <a-descriptions-item label="阈值">{{ currentPolicy.threshold }}</a-descriptions-item>
        <a-descriptions-item label="持续时间">{{ currentPolicy.duration }}秒</a-descriptions-item>
        <a-descriptions-item label="通知方式">{{ currentPolicy.notification.join(', ') }}</a-descriptions-item>
        <a-descriptions-item label="接收人">{{ currentPolicy.recipients.join(', ') }}</a-descriptions-item>
        <a-descriptions-item label="创建时间">{{ currentPolicy.created }}</a-descriptions-item>
        <a-descriptions-item label="更新时间">{{ currentPolicy.updated }}</a-descriptions-item>
        <a-descriptions-item label="描述" :span="2">{{ currentPolicy.description || '-' }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { message } from 'ant-design-vue'
import {
  PlusOutlined,
  ReloadOutlined,
  SearchOutlined,
  DownOutlined,
  EyeOutlined,
  EditOutlined,
  DeleteOutlined,
  CopyOutlined
} from '@ant-design/icons-vue'

interface Policy {
  id: string
  name: string
  type: 'cpu' | 'memory' | 'disk' | 'network' | 'service'
  status: 'enabled' | 'disabled'
  severity: 'critical' | 'warning' | 'info'
  target: string
  metric: string
  threshold: string
  duration: number
  notification: string[]
  recipients: string[]
  created: string
  updated: string
  description: string
}

const columns = [
  { title: '策略名', dataIndex: 'name', key: 'name' },
  { title: '监控类型', dataIndex: 'type', key: 'type', width: 120 },
  { title: '状态', dataIndex: 'status', key: 'status', width: 100 },
  { title: '严重级别', dataIndex: 'severity', key: 'severity', width: 100 },
  { title: '监控目标', dataIndex: 'target', key: 'target', width: 150 },
  { title: '阈值', dataIndex: 'threshold', key: 'threshold', width: 120 },
  { title: '通知方式', dataIndex: 'notification', key: 'notification', width: 150 },
  { title: '创建时间', dataIndex: 'created', key: 'created', width: 180 },
  { title: '操作', key: 'action', width: 100, fixed: 'right' }
]

const policies = ref<Policy[]>([
  {
    id: '1',
    name: 'CPU 使用率告警',
    type: 'cpu',
    status: 'enabled',
    severity: 'warning',
    target: '所有主机',
    metric: 'cpu_usage',
    threshold: '> 80%',
    duration: 60,
    notification: ['邮件', '钉钉'],
    recipients: ['admin@example.com', 'dev@example.com'],
    created: '2024-01-25 10:00:00',
    updated: '2024-01-25 10:00:00',
    description: '当 CPU 使用率超过 80% 持续 60 秒时触发告警'
  },
  {
    id: '2',
    name: '内存使用率告警',
    type: 'memory',
    status: 'enabled',
    severity: 'critical',
    target: '所有主机',
    metric: 'memory_usage',
    threshold: '> 90%',
    duration: 30,
    notification: ['邮件', '钉钉', '短信'],
    recipients: ['admin@example.com', 'ops@example.com'],
    created: '2024-01-24 14:30:00',
    updated: '2024-01-24 14:30:00',
    description: '当内存使用率超过 90% 持续 30 秒时触发告警'
  },
  {
    id: '3',
    name: '磁盘空间告警',
    type: 'disk',
    status: 'disabled',
    severity: 'critical',
    target: '所有主机',
    metric: 'disk_usage',
    threshold: '> 95%',
    duration: 120,
    notification: ['邮件', '钉钉'],
    recipients: ['admin@example.com'],
    created: '2024-01-26 09:00:00',
    updated: '2024-01-26 09:00:00',
    description: '当磁盘空间使用率超过 95% 持续 120 秒时触发告警'
  },
  {
    id: '4',
    name: '网络流量告警',
    type: 'network',
    status: 'enabled',
    severity: 'warning',
    target: 'server-01, server-02',
    metric: 'network_in',
    threshold: '> 100MB/s',
    duration: 60,
    notification: ['钉钉'],
    recipients: ['dev@example.com'],
    created: '2024-01-27 16:00:00',
    updated: '2024-01-27 16:00:00',
    description: '当网络入站流量超过 100MB/s 持续 60 秒时触发告警'
  },
  {
    id: '5',
    name: '服务状态告警',
    type: 'service',
    status: 'enabled',
    severity: 'critical',
    target: '所有主机',
    metric: 'service_status',
    threshold: '== down',
    duration: 10,
    notification: ['邮件', '钉钉', '短信'],
    recipients: ['admin@example.com', 'ops@example.com'],
    created: '2024-01-23 11:00:00',
    updated: '2024-01-23 11:00:00',
    description: '当服务状态变为 down 持续 10 秒时触发告警'
  },
  {
    id: '6',
    name: '容器资源告警',
    type: 'cpu',
    status: 'enabled',
    severity: 'info',
    target: '所有容器',
    metric: 'container_cpu_usage',
    threshold: '> 70%',
    duration: 60,
    notification: ['钉钉'],
    recipients: ['dev@example.com'],
    created: '2024-01-22 10:00:00',
    updated: '2024-01-22 10:00:00',
    description: '当容器 CPU 使用率超过 70% 持续 60 秒时触发告警'
  }
])

const searchForm = reactive({
  name: '',
  status: '',
  type: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 6
})

const loading = ref(false)
const detailVisible = ref(false)
const currentPolicy = ref<Policy | null>(null)

const getSeverityColor = (severity: string) => {
  const colors: Record<string, string> = {
    critical: 'error',
    warning: 'warning',
    info: 'default'
  }
  return colors[severity] || 'default'
}

const getSeverityText = (severity: string) => {
  const texts: Record<string, string> = {
    critical: '严重',
    warning: '警告',
    info: '提示'
  }
  return texts[severity] || severity
}

const getTypeText = (type: string) => {
  const texts: Record<string, string> = {
    cpu: 'CPU',
    memory: '内存',
    disk: '磁盘',
    network: '网络',
    service: '服务'
  }
  return texts[type] || type
}

const getActionItems = (record: Policy) => {
  const items: any[] = [
    { key: 'view', label: '查看', icon: EyeOutlined },
    { key: 'edit', label: '编辑', icon: EditOutlined },
    { key: 'copy', label: '复制', icon: CopyOutlined },
    { type: 'divider' },
    { key: 'delete', label: '删除', icon: DeleteOutlined, danger: true }
  ]
  
  return items
}

const handleAction = (key: string, record: Policy) => {
  switch (key) {
    case 'view':
      handleView(record)
      break
    case 'edit':
      handleEdit(record)
      break
    case 'copy':
      handleCopy(record)
      break
    case 'delete':
      handleDelete(record)
      break
  }
}

const handleSearch = () => {
  loading.value = true
  setTimeout(() => {
    loading.value = false
    message.success('搜索成功')
  }, 500)
}

const handleReset = () => {
  searchForm.name = ''
  searchForm.status = ''
  searchForm.type = ''
}

const handleRefresh = () => {
  message.success('刷新成功')
}

const handleCreate = () => {
  message.warning('新建策略功能待实现')
}

const handleView = (record: Policy) => {
  currentPolicy.value = record
  detailVisible.value = true
}

const handleEdit = (record: Policy) => {
  message.success(`编辑策略 ${record.name}`)
}

const handleCopy = (record: Policy) => {
  message.success(`复制策略 ${record.name}`)
}

const handleDelete = (record: Policy) => {
  message.warning(`删除策略 ${record.name} 功能待实现`)
}

const handleToggleStatus = (record: Policy, checked: boolean) => {
  const status = checked ? '启用' : '禁用'
  message.success(`${status} 策略 ${record.name} 成功`)
}
</script>