<template>
  <div>
    <a-card class="image-card">



      <a-table
        :columns="columns"
        :data-source="images"
        :pagination="pagination"
        :loading="loading"
        row-key="id"
        class="table-sm"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'size'">
            <a-tag>{{ formatSize(record.size) }}</a-tag>
          </template>
          <template v-else-if="column.key === 'status'">
            <a-tag :color="record.status === 'active' ? 'success' : 'default'">
              {{ record.status === 'active' ? '活跃' : '未使用' }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'containers'">
            <a-badge :count="record.containers" :color="record.containers > 0 ? 'success' : 'default'" />
          </template>
          <template v-else-if="column.key === 'action'">
            <a-popover
              placement="top"
              trigger="hover"
              :title="null"
            >
              <template #content>
                <div class="action-popover">
                  <div class="action-item" @click="handleRun(record)">
                    运行
                  </div>
                  <div class="action-item" @click="handleTag(record)">
                    打标签
                  </div>
                  <div class="action-item" @click="handlePush(record)">
                    推送
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
  SearchOutlined,
  CloudDownloadOutlined,
  DeleteOutlined,
  AppstoreOutlined,
  CheckCircleOutlined,
  HddOutlined,
  ContainerOutlined,
  ReloadOutlined,
  DownOutlined,
  EditOutlined,
  EyeOutlined
} from '@ant-design/icons-vue'
import type { MenuProps } from 'ant-design-vue'

interface Image {
  id: string
  name: string
  tag: string
  size: number
  status: 'active' | 'inactive'
  containers: number
  created: string
  digest: string
  architecture: string
  os: string
}

const columns = [
  { title: '镜像名', dataIndex: 'name', key: 'name' },
  { title: '标签', dataIndex: 'tag', key: 'tag', width: 100 },
  { title: '大小', dataIndex: 'size', key: 'size', width: 120 },
  { title: '状态', dataIndex: 'status', key: 'status', width: 100 },
  { title: '容器数', dataIndex: 'containers', key: 'containers', width: 100 },
  { title: '架构', dataIndex: 'architecture', key: 'architecture', width: 100 },
  { title: '创建时间', dataIndex: 'created', key: 'created', width: 180 },
  { title: '操作', key: 'action', width: 280, fixed: 'right' }
]

const images = ref<Image[]>([
  {
    id: '1',
    name: 'nginx',
    tag: '1.25-alpine',
    size: 24786944,
    status: 'active',
    containers: 2,
    created: '2024-01-15 10:30:00',
    digest: 'sha256:abc123...',
    architecture: 'amd64',
    os: 'linux'
  },
  {
    id: '2',
    name: 'mysql',
    tag: '8.0',
    size: 579869440,
    status: 'active',
    containers: 1,
    created: '2024-01-14 15:20:00',
    digest: 'sha256:def456...',
    architecture: 'amd64',
    os: 'linux'
  },
  {
    id: '3',
    name: 'redis',
    tag: '7-alpine',
    size: 31457280,
    status: 'active',
    containers: 1,
    created: '2024-01-16 09:00:00',
    digest: 'sha256:ghi789...',
    architecture: 'amd64',
    os: 'linux'
  },
  {
    id: '4',
    name: 'node',
    tag: '20-alpine',
    size: 1932735488,
    status: 'inactive',
    containers: 0,
    created: '2024-01-13 14:45:00',
    digest: 'sha256:jkl012...',
    architecture: 'amd64',
    os: 'linux'
  },
  {
    id: '5',
    name: 'postgres',
    tag: '15',
    size: 3221225472,
    status: 'active',
    containers: 1,
    created: '2024-01-12 11:30:00',
    digest: 'sha256:mno345...',
    architecture: 'amd64',
    os: 'linux'
  },
  {
    id: '6',
    name: 'ubuntu',
    tag: '22.04',
    size: 775946240,
    status: 'inactive',
    containers: 0,
    created: '2024-01-10 08:00:00',
    digest: 'sha256:pqr678...',
    architecture: 'amd64',
    os: 'linux'
  }
])

const searchForm = reactive({
  name: '',
  tag: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 6
})

const loading = ref(false)

const activeCount = computed(() => images.value.filter(img => img.status === 'active').length)

const totalSize = computed(() => {
  const bytes = images.value.reduce((sum, img) => sum + img.size, 0)
  return (bytes / 1024 / 1024 / 1024).toFixed(2)
})

const totalContainers = computed(() => images.value.reduce((sum, img) => sum + img.containers, 0))

const formatSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  if (bytes < 1024 * 1024 * 1024) return (bytes / 1024 / 1024).toFixed(2) + ' MB'
  return (bytes / 1024 / 1024 / 1024).toFixed(2) + ' GB'
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
  searchForm.tag = ''
}

const handlePull = () => {
  message.info('打开拉取镜像对话框')
}

const handlePrune = () => {
  message.success('已清理无用镜像')
}

const handleRun = (record: Image) => {
  message.success(`基于镜像 ${record.name}:${record.tag} 创建容器`)
}

const handleTag = (record: Image) => {
  message.success(`镜像 ${record.name} 打标签成功`)
}

const handlePush = (record: Image) => {
  message.success(`镜像 ${record.name} 推送成功`)
}

const handleDelete = (record: Image) => {
  if (record.containers > 0) {
    message.error(`镜像 ${record.name} 有 ${record.containers} 个容器在使用，无法删除`)
    return
  }
  message.success(`镜像 ${record.name}:${record.tag} 已删除`)
}
const handleRefresh = () => {
  loading.value = true
  setTimeout(() => {
    loading.value = false
    message.success('刷新成功')
  }, 1000)
}
const getActionMenu = (record: Image) => {
  return {
    items: [
      {
        key: 'run',
        label: '运行',
        onClick: () => handleRun(record)
      },
      {
        key: 'tag',
        label: '打标签',
        onClick: () => handleTag(record)
      },
      {
        key: 'push',
        label: '推送',
        onClick: () => handlePush(record)
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
</script>

<style scoped>
.image-card {
  overflow: visible;
}

.image-card :deep(.ant-card-head) {
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

.image-card :deep(.ant-card-head-title) {
  font-size: 12px !important;
  line-height: 32px !important;
  padding: 1px !important;
  margin: 0 !important;
  height: 32px !important;
}

.image-card :deep(.ant-card-body) {
  padding: 1px !important;
  margin: 0 !important;
}

.image-card :deep(.ant-card-extra) {
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