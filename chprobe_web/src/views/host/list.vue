<template>
  <div>
    <a-card class="host-card">
      <template #extra>
      </template>

      <a-table
        :columns="columns"
        :data-source="hosts"
        :pagination="pagination"
        :loading="loading"
        row-key="uuid"
        class="table-sm"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'hostname'">
            <a @click="handleView(record)">{{ record.hostname }}</a>
          </template>
          <template v-else-if="column.key === 'status'">
            <a-tag :color="record.status === 'online' ? 'success' : 'error'">
              {{ record.status === 'online' ? '在线' : '离线' }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'cpuUsage'">
            <a-progress
              :percent="record.cpuUsage"
              :size="'small'"
              :stroke-color="getColor(record.cpuUsage)"
            />
          </template>
          <template v-else-if="column.key === 'memoryUsage'">
            <a-progress
              :percent="record.memoryUsage"
              :size="'small'"
              :stroke-color="getColor(record.memoryUsage)"
            />
          </template>
          <template v-else-if="column.key === 'action'">
            <a-popover
              placement="top"
              trigger="hover"
              :title="null"
            >
              <template #content>
                <div class="action-popover">
                  <div class="action-item" @click="handleView(record)">
                    <EyeOutlined style="margin-right: 4px;" />
                    查看
                  </div>
                  <div class="action-item" @click="handleEdit(record)">
                    <EditOutlined style="margin-right: 4px;" />
                    编辑
                  </div>
                  <div class="action-item" @click="handleRefresh(record)">
                    <ReloadOutlined style="margin-right: 4px;" />
                    刷新
                  </div>
                  <div class="action-item action-danger" @click="handleDelete(record)">
                    <DeleteOutlined style="margin-right: 4px;" />
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

    <a-modal
      v-model:open="modalVisible"
      :title="modalTitle"
      @ok="handleModalOk"
      @cancel="handleModalCancel"
      :centered="true"
    >
      <a-form :model="hostForm" :label-col="{ span: 6 }" :wrapper-col="{ span: 16 }">
        <a-form-item label="主机名">
          <a-input v-model:value="hostForm.hostname" placeholder="请输入主机名" />
        </a-form-item>
        <a-form-item label="IP地址">
          <a-input v-model:value="hostForm.ip" placeholder="请输入IP地址" />
        </a-form-item>
        <a-form-item label="端口">
          <a-input-number v-model:value="hostForm.port" :min="1" :max="65535" style="width: 100%;" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model:value="hostForm.description" placeholder="请输入描述" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal
      v-model:open="detailVisible"
      title="主机详情"
      width="640px"
      :footer="null"
      :centered="true"
    >
      <div class="detail-content" style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px; padding: 16px;">
        <div v-for="(item, index) in detailData" :key="item.key" class="detail-item">
          <div class="detail-label">{{ item.label }}:</div>
          <div class="detail-value">
            <template v-if="item.status">
              <a-tag :color="item.status === 'online' ? 'success' : 'error'">
                {{ item.status === 'online' ? '在线' : '离线' }}
              </a-tag>
            </template>
            <template v-else-if="item.cpuUsage !== undefined">
              <a-progress :percent="item.cpuUsage" :stroke-color="getColor(item.cpuUsage)" :size="'small'" />
            </template>
            <template v-else-if="item.memoryUsage !== undefined">
              <a-progress :percent="item.memoryUsage" :stroke-color="getColor(item.memoryUsage)" :size="'small'" />
            </template>
            <template v-else>
              {{ item.value }}
            </template>
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { message } from 'ant-design-vue'
import {
  PlusOutlined,
  SearchOutlined,
  DownOutlined,
  EyeOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined
} from '@ant-design/icons-vue'

interface Host {
  id: string
  hostname: string
  ip: string
  port: number
  status: 'online' | 'offline'
  cpuUsage: number
  memoryUsage: number
  os: string
  kernel: string
  lastSeen: string
  description: string
}

const columns = [
 { title: '主机名', dataIndex: 'hostname', key: 'hostname', width: 150 },
 { title: 'IP地址', dataIndex: 'ip', key: 'ip', width: 150 },
 { title: '状态', dataIndex: 'status', key: 'status', width: 100 },
 { title: 'CPU', dataIndex: 'cpuUsage', key: 'cpuUsage', width: 120 },
 { title: '内存', dataIndex: 'memoryUsage', key: 'memoryUsage', width: 120 },
 { title: '操作系统', dataIndex: 'os', key: 'os', width: 150 },
 { title: '操作', key: 'action', width: 100, fixed: 'right' }
];
const detailColumns = [
 { title: '属性', dataIndex: 'key', key: 'key', width: 150 },
 { title: '值', dataIndex: 'value', key: 'value' }
];
const detailData = computed(() => {
 if (!currentHost.value) return [];
 return [
 { key: 'hostname', label: '主机名', value: currentHost.value.hostname },
 { key: 'ip', label: 'IP地址', value: currentHost.value.ip },
 { key: 'port', label: '端口', value: currentHost.value.port },
 { key: 'status', label: '状态', value: currentHost.value.status, status: currentHost.value.status },
 { key: 'cpuUsage', label: 'CPU使用率', value: `${currentHost.value.cpuUsage}%`, cpuUsage: currentHost.value.cpuUsage },
 { key: 'memoryUsage', label: '内存使用率', value: `${currentHost.value.memoryUsage}%`, memoryUsage: currentHost.value.memoryUsage },
 { key: 'os', label: '操作系统', value: currentHost.value.os },
 { key: 'kernel', label: '内核版本', value: currentHost.value.kernel },
 { key: 'lastSeen', label: '最后在线', value: currentHost.value.lastSeen },
 { key: 'description', label: '描述', value: currentHost.value.description || '-' }
 ];
});

const hosts = ref<Host[]>([
  {
    id: '1',
    hostname: 'server-01',
    ip: '192.168.1.101',
    port: 22,
    status: 'online',
    cpuUsage: 45,
    memoryUsage: 62,
    os: 'Ubuntu 22.04 LTS',
    kernel: '5.15.0-52-generic',
    lastSeen: '2024-01-28 10:30:00',
    description: '生产环境服务器'
  },
  {
    id: '2',
    hostname: 'server-02',
    ip: '192.168.1.102',
    port: 22,
    status: 'online',
    cpuUsage: 32,
    memoryUsage: 45,
    os: 'CentOS 7.9',
    kernel: '3.10.0-1160.el7.x86_64',
    lastSeen: '2024-01-28 10:28:00',
    description: '测试环境服务器'
  },
  {
    id: '3',
    hostname: 'server-03',
    ip: '192.168.1.103',
    port: 22,
    status: 'offline',
    cpuUsage: 0,
    memoryUsage: 0,
    os: 'Debian 11',
    kernel: '5.10.0-8-amd64',
    lastSeen: '2024-01-28 09:00:00',
    description: '开发环境服务器'
  },
  {
    id: '4',
    hostname: 'server-04',
    ip: '192.168.1.104',
    port: 22,
    status: 'online',
    cpuUsage: 78,
    memoryUsage: 85,
    os: 'Red Hat Enterprise Linux 8',
    kernel: '4.18.0-305.el8.x86_64',
    lastSeen: '2024-01-28 10:32:00',
    description: '数据库服务器'
  },
  {
    id: '5',
    hostname: 'server-05',
    ip: '192.168.1.105',
    port: 22,
    status: 'online',
    cpuUsage: 28,
    memoryUsage: 55,
    os: 'SUSE Linux Enterprise 15',
    kernel: '5.3.18-59.30-default',
    lastSeen: '2024-01-28 10:25:00',
    description: '应用服务器'
  }
])

const searchForm = reactive({
  hostname: '',
  status: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 5
})

const loading = ref(false)
const modalVisible = ref(false)
const detailVisible = ref(false)
const modalTitle = ref('添加主机')
const currentHost = ref<Host | null>(null)

const hostForm = reactive({
  id: '',
  hostname: '',
  ip: '',
  port: 22,
  description: ''
})

const getColor = (value: number) => {
  if (value >= 80) return '#f5222d'
  if (value >= 60) return '#faad14'
  return '#52c41a'
}

const getActionMenu = (record: Host) => {
  const items = [
    { key: 'view', label: '查看',  },
    { key: 'edit', label: '编辑'},
    { key: 'refresh', label: '刷新'},
    { type: 'divider' },
    { key: 'delete', label: '删除'}
  ]
  return {
    items,
    onClick: ({ key }: { key: string }) => handleAction(key, record)
  }
}

const handleAction = (key: string, record: Host) => {
  switch (key) {
    case 'view':
      handleView(record)
      break
    case 'edit':
      handleEdit(record)
      break
    case 'refresh':
      handleRefresh(record)
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
  searchForm.hostname = ''
  searchForm.status = ''
}

const handleAdd = () => {
  modalTitle.value = '添加主机'
  hostForm.id = ''
  hostForm.hostname = ''
  hostForm.ip = ''
  hostForm.port = 22
  hostForm.description = ''
  modalVisible.value = true
}

const handleView = (record: Host) => {
  currentHost.value = record
  detailVisible.value = true
}

const handleEdit = (record: Host) => {
  modalTitle.value = '编辑主机'
  hostForm.id = record.id
  hostForm.hostname = record.hostname
  hostForm.ip = record.ip
  hostForm.port = record.port
  hostForm.description = record.description
  modalVisible.value = true
}

const handleRefresh = (record: Host) => {
  message.success(`刷新 ${record.hostname} 成功`)
}

const handleDelete = (record: Host) => {
  message.warning(`删除 ${record.hostname} 功能待实现`)
}

const handleModalOk = () => {
  if (hostForm.id) {
    message.success('编辑成功')
  } else {
    message.success('添加成功')
  }
  modalVisible.value = false
}

const handleModalCancel = () => {
  modalVisible.value = false
}
</script>

<style scoped>
.host-card {
  overflow: visible;
}

.host-card :deep(.ant-card-head) {
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

.host-card :deep(.ant-card-head-title) {
  font-size: 12px !important;
  line-height: 32px !important;
  padding: 1px !important;
  margin: 0 !important;
  height: 32px !important;
}

.host-card :deep(.ant-card-body) {
  padding: 1px !important;
  margin: 0 !important;
}

.host-card :deep(.ant-card-extra) {
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
}

.detail-table :deep(.ant-table-thead > tr > th) {
  height: 32px !important;
  line-height: 32px !important;
  padding: 1px !important;
  font-size: 12px !important;
}

.detail-table :deep(.ant-table-tbody > tr > td) {
  height: 32px !important;
  line-height: 32px !important;
  padding: 1px !important;
  font-size: 12px !important;
}

.action-popover {
  padding: 2px 0;
  width: 50px;
}

.action-item {
  padding: 3px 2px;
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

.detail-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0;
  border-bottom: 1px solid #f0f0f0;
}

.detail-item:last-child {
  border-bottom: none;
}

.detail-label {
  font-size: 12px;
  color: #666;
  font-weight: normal;
  min-width: 80px;
  flex-shrink: 0;
}

.detail-value {
  font-size: 12px;
  color: #333;
  min-height: 20px;
  display: flex;
  align-items: center;
  flex: 1;
}

.detail-content {
  max-height: 400px;
  overflow-y: auto;
}
</style>