<template>
  <div>
    <a-row :gutter="[16, 16]">
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card>
          <a-statistic title="客户端总数" :value="stats.clients.total">
            <template #prefix>
              <CloudServerOutlined />
            </template>
          </a-statistic>
          <a-progress :percent="stats.clients.onlineRate" :stroke-color="'#52c41a'" style="margin-top: 16px;" />
          <div style="margin-top: 8px; font-size: 12px; color: #8c8c8c;">
            {{ stats.clients.online }} 台在线 / {{ stats.clients.offline }} 台离线
          </div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card>
          <a-statistic title="插件总数" :value="stats.plugins.total">
            <template #prefix>
              <AppstoreOutlined />
            </template>
          </a-statistic>
          <a-progress :percent="stats.plugins.installedRate" :stroke-color="'#1890ff'" style="margin-top: 16px;" />
          <div style="margin-top: 8px; font-size: 12px; color: #8c8c8c;">
            {{ stats.plugins.installed }} 个已安装 / {{ stats.plugins.available }} 个可用
          </div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card>
          <a-statistic title="授权状态" :value="stats.auth.status">
            <template #prefix>
              <SafetyOutlined />
            </template>
          </a-statistic>
          <a-progress :percent="stats.auth.progress" :stroke-color="stats.auth.color" style="margin-top: 16px;" />
          <div style="margin-top: 8px; font-size: 12px; color: #8c8c8c;">
            {{ stats.auth.message }}
          </div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card>
          <a-statistic title="用户数量" :value="stats.users.total">
            <template #prefix>
              <UserOutlined />
            </template>
          </a-statistic>
          <a-progress :percent="stats.users.activeRate" :stroke-color="'#722ed1'" style="margin-top: 16px;" />
          <div style="margin-top: 8px; font-size: 12px; color: #8c8c8c;">
            {{ stats.users.active }} 个活跃 / {{ stats.users.inactive }} 个未激活
          </div>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="[16, 16]" style="margin-top: 16px;">
      <a-col :xs="24" :lg="16">
        <a-card title="最近系统日志">
          <a-table :columns="logColumns" :data-source="recentLogs" :pagination="false" :row-key="'id'">
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'level'">
                <a-tag :color="getLogColor(record.level)">{{ record.level }}</a-tag>
              </template>
              <template v-if="column.key === 'time'">
                {{ formatTime(record.time) }}
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
      <a-col :xs="24" :lg="8">
        <a-card title="系统资源概览">
          <div style="margin-top: 16px;">
            <div style="display: flex; justify-content: space-between; margin-bottom: 8px;">
              <span>CPU 使用率</span>
              <span style="color: #52c41a;">{{ resourceUsage.cpu }}%</span>
            </div>
            <a-progress :percent="resourceUsage.cpu" :stroke-color="'#52c41a'" />
          </div>
          <div style="margin-top: 16px;">
            <div style="display: flex; justify-content: space-between; margin-bottom: 8px;">
              <span>内存使用率</span>
              <span style="color: '#faad14'">{{ resourceUsage.memory }}%</span>
            </div>
            <a-progress :percent="resourceUsage.memory" :stroke-color="'#faad14'" />
          </div>
          <div style="margin-top: 16px;">
            <div style="display: flex; justify-content: space-between; margin-bottom: 8px;">
              <span>磁盘使用率</span>
              <span style="color: '#f5222d'">{{ resourceUsage.disk }}%</span>
            </div>
            <a-progress :percent="resourceUsage.disk" :stroke-color="'#f5222d'" />
          </div>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="[16, 16]" style="margin-top: 16px;">
      <a-col :xs="24" :lg="12">
        <a-card title="客户端状态分布">
          <div style="height: 300px;">
            <a-empty description="饼图占位" />
          </div>
        </a-card>
      </a-col>
      <a-col :xs="24" :lg="12">
        <a-card title="插件使用分布">
          <div style="height: 300px;">
            <a-empty description="饼图占位" />
          </div>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { CloudServerOutlined, AppstoreOutlined, SafetyOutlined, UserOutlined } from '@ant-design/icons-vue'
import { getDashboardStats } from '@/api/index'

const stats = reactive({
  clients: {
    total: 0,
    online: 0,
    offline: 0,
    onlineRate: 0
  },
  plugins: {
    total: 0,
    installed: 0,
    available: 0,
    installedRate: 0
  },
  auth: {
    status: '未授权',
    progress: 0,
    color: '#f5222d',
    message: '请上传授权文件'
  },
  users: {
    total: 0,
    active: 0,
    inactive: 0,
    activeRate: 0
  }
})

const resourceUsage = reactive({
  cpu: 45,
  memory: 62,
  disk: 78
})

const logColumns = [
  { title: '日志级别', dataIndex: 'level', key: 'level', width: 80 },
  { title: '模块', dataIndex: 'module', key: 'module', width: 120 },
  { title: '内容', dataIndex: 'content', key: 'content' },
  { title: '时间', dataIndex: 'time', key: 'time', width: 150 }
]

const recentLogs = ref([
  { id: 1, level: 'INFO', module: '客户端', content: '客户端 client-01 成功连接', time: Date.now() - 3600000 },
  { id: 2, level: 'WARN', module: '插件', content: '插件 plugin-02 加载失败', time: Date.now() - 7200000 },
  { id: 3, level: 'ERROR', module: '授权', content: '授权验证失败，请检查授权文件', time: Date.now() - 10800000 },
  { id: 4, level: 'INFO', module: '系统', content: '系统备份完成', time: Date.now() - 14400000 },
  { id: 5, level: 'WARN', module: '客户端', content: '客户端 client-03 响应超时', time: Date.now() - 18000000 }
])

const fetchDashboardStats = async () => {
  try {
    const response = await getDashboardStats()
    const data = response.result
    
    if (data.clients) {
      stats.clients.total = data.clients.total || 0
      stats.clients.online = data.clients.online || 0
      stats.clients.offline = data.clients.offline || 0
      stats.clients.onlineRate = data.clients.onlineRate || 0
    }
    
    if (data.plugins) {
      stats.plugins.total = data.plugins.total || 0
      stats.plugins.installed = data.plugins.installed || 0
      stats.plugins.available = data.plugins.available || 0
      stats.plugins.installedRate = data.plugins.installedRate || 0
    }
    
    if (data.auth) {
      stats.auth.status = data.auth.status || '未授权'
      stats.auth.progress = data.auth.progress || 0
      stats.auth.color = data.auth.color || '#f5222d'
      stats.auth.message = data.auth.message || '请上传授权文件'
    }
    
    if (data.users) {
      stats.users.total = data.users.total || 0
      stats.users.active = data.users.active || 0
      stats.users.inactive = data.users.inactive || 0
      stats.users.activeRate = data.users.activeRate || 0
    }
  } catch (error) {
    console.error('Failed to fetch dashboard stats:', error)
  }
}

onMounted(() => {
  fetchDashboardStats()
})

const getLogColor = (level: string) => {
  const colors: Record<string, string> = {
    'ERROR': 'red',
    'WARN': 'orange',
    'INFO': 'blue',
    'DEBUG': 'gray'
  }
  return colors[level] || 'gray'
}

const formatTime = (timestamp: number) => {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  if (diff < 60000) {
    return '刚刚'
  } else if (diff < 3600000) {
    return Math.floor(diff / 60000) + ' 分钟前'
  } else if (diff < 86400000) {
    return Math.floor(diff / 3600000) + ' 小时前'
  } else {
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString().substring(0, 5)
  }
}
</script>
