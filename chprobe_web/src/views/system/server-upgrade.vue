<template>
  <div>
    <a-card style="margin-bottom: 24px;" title="升级记录">
      <a-table 
        :columns="historyColumns" 
        :data-source="upgradeRecords" 
        :pagination="pagination"
        :loading="loading"
        row-key="uuid"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 'success' ? 'green' : 'red'">
              {{ record.status === 'success' ? '成功' : '失败' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'upgradeType'">
            <a-tag color="blue">{{ record.upgradeType }}</a-tag>
          </template>
          <template v-if="column.key === 'upgradeTime'">
            {{ formatTime(record.upgradeTime) }}
          </template>
        </template>
      </a-table>
    </a-card>

    <a-alert
      message="升级前请确保"
      description="1. 已备份重要数据；2. 系统有足够的磁盘空间；3. 升级过程中服务可能会短暂中断"
      type="warning"
      show-icon
      style="margin-bottom: 24px;"
    />

    <a-steps :current="0" status="process" style="margin-bottom: 24px;">
      <a-step title="检查版本" />
      <a-step title="下载升级包" />
      <a-step title="执行升级" />
      <a-step title="验证升级" />
    </a-steps>

    <a-card size="small">
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="升级方式">
          <a-radio-group v-model:value="upgradeMethod">
            <a-radio value="online">在线升级</a-radio>
            <a-radio value="offline">离线升级</a-radio>
          </a-radio-group>
        </a-descriptions-item>
        <a-descriptions-item label="升级命令" v-if="upgradeMethod === 'online'">
          <a-input
            :value="onlineUpgradeCmd"
            readonly
            style="margin-bottom: 8px;"
          />
          <a-button type="primary" @click="copyOnlineCmd" block>
            <template #icon>
              <CopyOutlined />
            </template>
            复制命令
          </a-button>
        </a-descriptions-item>
        <a-descriptions-item label="升级包下载" v-else>
          <a-button type="primary" @click="downloadUpgradePackage">
            <template #icon>
              <DownloadOutlined />
            </template>
            下载升级包
          </a-button>
          <p style="margin-top: 8px; color: #999; font-size: 12px;">
            下载完成后，上传到服务器并执行安装脚本
          </p>
        </a-descriptions-item>
        <a-descriptions-item label="升级日志">
          <a-button @click="viewUpgradeLogs">
            <template #icon>
              <FileTextOutlined />
            </template>
            查看升级日志
          </a-button>
        </a-descriptions-item>
      </a-descriptions>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { CopyOutlined, DownloadOutlined, FileTextOutlined } from '@ant-design/icons-vue'
import { getUpgradeRecordList } from '@/api'
import type { PageQuery, UpgradeRecord } from '@/api'

// 服务端升级相关数据
const currentVersion = 'v1.0.0'
const currentVersionDate = '2024-01-01'
const latestVersion = 'v1.1.0'
const needUpgrade = true
const upgradeMethod = ref('online')

const loading = ref(false)
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const upgradeRecords = ref<UpgradeRecord[]>([])

const onlineUpgradeCmd = computed(() => {
  return `curl -fsSL http://your-server-ip:8080/upgrade | bash -s -- --version ${latestVersion}`
})

const historyColumns = [
  { title: '版本号', dataIndex: 'version', key: 'version', width: 120 },
  { title: '上一版本', dataIndex: 'previousVersion', key: 'previousVersion', width: 120 },
  { title: '升级类型', dataIndex: 'upgradeType', key: 'upgradeType', width: 100 },
  { title: '状态', dataIndex: 'status', key: 'status', width: 100 },
  { title: '升级时间', dataIndex: 'upgradeTime', key: 'upgradeTime', width: 180 },
  { title: '服务器IP', dataIndex: 'serverIp', key: 'serverIp', width: 140 },
  { title: '主机名', dataIndex: 'hostname', key: 'hostname', width: 140 },
  { title: '操作人', dataIndex: 'operator', key: 'operator', width: 100 },
  { title: '升级描述', dataIndex: 'description', key: 'description' }
]

const formatTime = (timestamp: number): string => {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const fetchUpgradeRecords = async () => {
  console.log('调用获取升级记录API')
  loading.value = true
  try {
    const params: PageQuery = {
      page: pagination.current,
      pageSize: pagination.pageSize
    }
    const response = await getUpgradeRecordList(params)
    console.log('API响应:', response)
    upgradeRecords.value = response.result
    pagination.total = response.total
  } catch (error) {
    console.error('获取升级记录失败:', error)
    message.error('获取升级记录失败')
  } finally {
    loading.value = false
  }
}

const handleTableChange = (pag: any) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchUpgradeRecords()
}

// 服务端升级方法
const copyOnlineCmd = () => {
  navigator.clipboard.writeText(onlineUpgradeCmd.value).then(() => {
    message.success('命令已复制到剪贴板')
  })
}

const downloadUpgradePackage = () => {
  message.info(`开始下载 ${latestVersion} 版本升级包...`)
}

const viewUpgradeLogs = () => {
  message.info('打开升级日志查看器')
}

console.log('服务端升级页面加载')
fetchUpgradeRecords()

onMounted(() => {
  console.log('服务端升级页面onMounted')
  fetchUpgradeRecords()
})
</script>

<style scoped>
/* 可根据需要添加样式 */
</style>