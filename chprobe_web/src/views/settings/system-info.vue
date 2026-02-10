<template>
  <div>
    <a-row>
      <a-col :span="12">
        <a-card>
          <template v-if="loading">
            <div style="display: flex; justify-content: center; padding: 40px;">
              <a-spin tip="加载中..." size="large" />
            </div>
          </template>
          <template v-else-if="error">
            <div style="padding: 40px;">
              <a-alert
                message="加载失败"
                description="无法获取系统信息，请稍后重试"
                type="error"
                show-icon
                action
              >
                <template #action>
                  <a-button type="primary" size="small" @click="fetchSystemInfo">
                    重新加载
                  </a-button>
                </template>
              </a-alert>
            </div>
          </template>
          <template v-else>
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px;">
              <h2 style="margin: 0;">{{ serverInfo.hostname }} ({{ serverInfo.role }})</h2>
              <a-tag :color="serverInfo.status === '正常' ? 'green' : 'red'" style="font-size: 16px; padding: 4px 12px;">
                {{ serverInfo.status }}
              </a-tag>
            </div>

            <a-card style="margin-bottom: 24px;">
              <a-descriptions :column="1" bordered size="small" style="line-height: 1.4;">
                <a-descriptions-item label="产品名称">
                  {{ serverInfo.productName }}
                </a-descriptions-item>
                <a-descriptions-item label="编译时间">
                  {{ serverInfo.buildTime }}
                </a-descriptions-item>
                <a-descriptions-item label="CommitID">
                  {{ serverInfo.commitID }}
                </a-descriptions-item>
                <a-descriptions-item label="启动时间">
                  {{ formatStartupTime(serverInfo.startupTime) }}
                </a-descriptions-item>
                <a-descriptions-item label="内核版本">
                  {{ serverInfo.kernel }}
                </a-descriptions-item>
                <a-descriptions-item label="主机名">
                  {{ serverInfo.hostname }}
                </a-descriptions-item>
                <a-descriptions-item label="IP">
                  {{ serverInfo.ip }}
                </a-descriptions-item>
                <a-descriptions-item label="序列号">
                  {{ serverInfo.serial }}
                </a-descriptions-item>
                <a-descriptions-item label="版本">
                  {{ serverInfo.version }}
                </a-descriptions-item>
                <a-descriptions-item label="角色">
                  {{ serverInfo.role }}
                </a-descriptions-item>
              </a-descriptions>
            </a-card>

            <a-row :gutter="16">
              <a-col :span="8">
                <a-card size="small">
                  <div style="display: flex; align-items: center; justify-content: space-between;">
                    <div style="display: flex; align-items: center;">
                      <SettingOutlined style="margin-right: 8px; color: #1890ff;" />
                      <span>CPU使用</span>
                    </div>
                    <span style="font-size: 24px; font-weight: bold;">{{ serverInfo.cpuUsage }}%</span>
                  </div>
                  <div style="font-size: 12px; color: #666; margin-top: 4px; margin-bottom: 8px;">{{ serverInfo.cpuConfig }}</div>
                  <a-progress :percent="parseInt(serverInfo.cpuUsage)" :stroke-color="getProgressColor(parseInt(serverInfo.cpuUsage))" style="margin-top: 8px;" />
                </a-card>
              </a-col>
              <a-col :span="8">
                <a-card size="small">
                  <div style="display: flex; align-items: center; justify-content: space-between;">
                    <div style="display: flex; align-items: center;">
                      <ApiOutlined style="margin-right: 8px; color: #52c41a;" />
                      <span>内存占用</span>
                    </div>
                    <span style="font-size: 24px; font-weight: bold;">{{ serverInfo.memoryUsage }}%</span>
                  </div>
                  <div style="font-size: 12px; color: #666; margin-top: 4px; margin-bottom: 8px;">{{ serverInfo.memoryConfig }}</div>
                  <a-progress :percent="parseInt(serverInfo.memoryUsage)" :stroke-color="getProgressColor(parseInt(serverInfo.memoryUsage))" style="margin-top: 8px;" />
                </a-card>
              </a-col>
              <a-col :span="8">
                <a-card size="small">
                  <div style="display: flex; align-items: center; justify-content: space-between;">
                    <div style="display: flex; align-items: center;">
                      <DatabaseOutlined style="margin-right: 8px; color: #faad14;" />
                      <span>磁盘剩余</span>
                    </div>
                    <span style="font-size: 24px; font-weight: bold;">{{ serverInfo.diskUsage }}%</span>
                  </div>
                  <div style="font-size: 12px; color: #666; margin-top: 4px; margin-bottom: 8px;">{{ serverInfo.diskConfig }}</div>
                  <a-progress :percent="100 - parseInt(serverInfo.diskUsage)" :stroke-color="getProgressColor(100 - parseInt(serverInfo.diskUsage))" style="margin-top: 8px;" />
                </a-card>
              </a-col>
            </a-row>
          </template>
        </a-card>
      </a-col>
      <a-col :span="12">
        <!-- 右侧空白区域 -->
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { SettingOutlined, ApiOutlined, DatabaseOutlined } from '@ant-design/icons-vue'
import { getSystemInfo } from '@/api'

// 响应式数据
const loading = ref(true)
const error = ref(false)
const serverInfo = reactive({
  hostname: '',
  role: '管理中心',
  status: '正常',
  ip: '',
  serial: '',
  version: '',
  cpuConfig: '',
  memoryConfig: '',
  diskConfig: '',
  cpuUsage: '0',
  memoryUsage: '0',
  diskUsage: '0',
  productName: '',
  buildTime: '',
  commitID: '',
  startupTime: 0,
  kernel: ''
})

// 格式化启动时间
const formatStartupTime = (timestamp: number) => {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN')
}

// 获取进度条颜色
const getProgressColor = (percent: number) => {
  if (percent < 60) return '#52c41a'
  if (percent < 80) return '#faad14'
  return '#f5222d'
}

// 获取系统信息
const fetchSystemInfo = async () => {
  loading.value = true
  error.value = false
  try {
    const response = await getSystemInfo()
    const data = response.result
    
    // 更新响应式数据
    Object.assign(serverInfo, data)
  } catch (err) {
    console.error('获取系统信息失败:', err)
    error.value = true
    message.error('获取系统信息失败')
  } finally {
    loading.value = false
  }
}

// 组件挂载时获取数据
onMounted(() => {
  fetchSystemInfo()
})
</script>

<style scoped>
/* 可根据需要添加样式 */
</style>