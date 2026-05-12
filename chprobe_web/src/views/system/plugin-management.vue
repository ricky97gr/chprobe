<template>
  <div>
    <a-card>
      <!-- 我的插件页面 -->
      <div v-if="activeTab === 'my'">
        <a-alert
          message="我的插件"
          description="已安装的插件列表，支持启动、停止、重启、查看路由和健康检查操作。"
          type="info"
          show-icon
          style="margin-bottom: 24px;"
        />

        <a-card>
          <a-table :columns="myPluginColumns" :data-source="myPlugins" row-key="pluginId" :loading="loading">
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <a-tag :color="getStatusColor(record.status)">
                  {{ getStatusText(record.status) }}
                </a-tag>
              </template>
              <template v-if="column.key === 'version'">
                <span>{{ record.version }}</span>
              </template>
              <template v-if="column.key === 'action'">
                <a-dropdown>
                  <a-button type="text" size="small">
                    操作 <DownOutlined />
                  </a-button>
                  <template #overlay>
                    <a-menu>
                      <a-menu-item key="start" @click="startPlugin(record.pluginId)" :disabled="record.status === 'running' || record.loading">
                        启动
                      </a-menu-item>
                      <a-menu-item key="stop" @click="stopPlugin(record.pluginId)" :disabled="record.status !== 'running' || record.loading">
                        停止
                      </a-menu-item>
                      <a-menu-item key="restart" @click="restartPlugin(record.pluginId)" :disabled="record.loading">
                        重启
                      </a-menu-item>
                      <a-menu-item key="routes" @click="viewRoutes(record.pluginId)">
                        路由
                      </a-menu-item>
                      <a-menu-item key="health" @click="healthCheck(record.pluginId)">
                        健康
                      </a-menu-item>
                      <a-menu-divider />
                      <a-menu-item key="delete" @click="deletePlugin(record.pluginId)" :disabled="record.loading" danger>
                        删除
                      </a-menu-item>
                    </a-menu>
                  </template>
                </a-dropdown>
              </template>
            </template>
            <template #description="{ text }">
              <a-tooltip placement="top" :title="text">
                <span style="display: inline-block; width: 100%; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ text }}</span>
              </a-tooltip>
            </template>
          </a-table>
        </a-card>
      </div>

      <!-- 插件市场页面 -->
      <div v-else-if="activeTab === 'market'">
        <a-alert
          message="插件市场"
          description="可下载和安装的插件列表，支持搜索和筛选功能。"
          type="info"
          show-icon
          style="margin-bottom: 24px;"
        />

        <div style="margin-bottom: 16px; display: flex; justify-content: space-between; align-items: center;">
          <div style="display: flex; gap: 12px; align-items: center;">
            <a-input-search 
              v-model:value="searchKeyword" 
              placeholder="搜索插件"
              style="width: 300px;"
            />
            <a-select 
              v-model:value="categoryFilter" 
              placeholder="选择分类"
              style="width: 150px;"
            >
              <a-select-option value="">全部</a-select-option>
              <a-select-option value="monitor">监控</a-select-option>
              <a-select-option value="security">安全</a-select-option>
              <a-select-option value="integration">集成</a-select-option>
              <a-select-option value="tool">工具</a-select-option>
            </a-select>
          </div>
          <div style="display: flex; gap: 12px; align-items: center;">
            <a-input 
              v-model:value="marketIp" 
              placeholder="请输入插件市场IP地址"
              style="width: 300px;"
            />
            <a-button type="primary" @click="saveMarketIp">
              确定
            </a-button>
            <a-button type="primary" @click="syncMarketPlugins" :loading="syncing">
              <template #icon>
                <SyncOutlined :spin="syncing" />
              </template>
              同步
            </a-button>
          </div>
        </div>

        <a-row :gutter="[24, 24]">
          <a-col :xs="24" :sm="12" :md="8" :lg="8" v-for="plugin in processedMarketPlugins" :key="plugin.id">
            <a-card hoverable style="padding: 16px; min-width: 300px; min-height: 280px;">
              <template #head>
                  <div style="display: flex; justify-content: flex-end; align-items: center;">
                    <a-tag v-if="plugin.installed" color="green">已下载</a-tag>
                    <a-tag v-else color="gray">未下载</a-tag>
                  </div>
                </template>
              <div style="height: 200px; display: flex; flex-direction: column; justify-content: space-between;">
                <div>
                  <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
                    <h3 style="margin: 0;">{{ plugin.name }}</h3>
                    <div style="display: flex; align-items: center;">
                      <a-tag v-if="plugin.versionType" color="purple">{{ plugin.versionType }}</a-tag>
                      <a-tooltip v-if="plugin.tips" placement="top" :title="plugin.tips" style="margin-left: 8px;">
                        <span style="display: inline-block; width: 16px; height: 16px; border-radius: 50%; background-color: #faad14; color: white; text-align: center; line-height: 16px; font-size: 12px; cursor: help;">i</span>
                      </a-tooltip>
                    </div>
                  </div>
                  <a-tooltip :title="plugin.description" placement="top">
                    <p style="margin-bottom: 16px; color: #666; line-height: 1.5; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; text-overflow: ellipsis; white-space: normal;">{{ plugin.description }}</p>
                  </a-tooltip>
                </div>
                <div>
                  <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 16px;">
                    <div style="display: flex; flex-direction: column; gap: 8px;">
                      <a-tag color="blue">插件版本: v{{ plugin.version }}</a-tag>
                      <a-tag color="orange">服务端最低版本: v{{ plugin.miniServerVersion }}</a-tag>
                      <a-tag color="green">客户端最低版本: v{{ plugin.miniClientVersion }}</a-tag>
                    </div>
                    <div style="text-align: right;">
                      <p style="margin: 0 0 4px 0; font-size: 12px; color: #999;">作者: {{ plugin.author }}</p>
                      <p style="margin: 0; font-size: 12px; color: #999;">上传时间: {{ formatTime(plugin.uploadTime) }}</p>
                    </div>
                  </div>
                  <!-- 下载进度条 -->
                  <div v-if="(downloadProgress as any)[plugin.id] !== undefined && (downloadProgress as any)[plugin.id] < 100" style="margin-bottom: 16px;">
                    <a-progress :percent="(downloadProgress as any)[plugin.id]" status="active" />
                  </div>
                </div>
              </div>
              <div style="margin-top: 16px;">
                <a-button 
                  type="primary" 
                  block 
                  @click="installPlugin(plugin.uuid)"
                  :disabled="plugin.installed || downloadProgress[plugin.uuid] !== undefined"
                >
                  {{ plugin.installed ? '已下载' : downloadProgress[plugin.uuid] !== undefined ? '下载中' : '下载' }}
                </a-button>
              </div>
            </a-card>
          </a-col>
        </a-row>
      </div>
    </a-card>

    <!-- 路由查看对话框 -->
    <a-modal
      v-model:open="routesModalVisible"
      title="插件路由"
      :footer="null"
      width="800px"
    >
      <a-table :columns="routeColumns" :data-source="currentRoutes" row-key="path" :pagination="false">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'method'">
            <a-tag :color="getMethodColor(record.method)">{{ record.method }}</a-tag>
          </template>
        </template>
      </a-table>
    </a-modal>

    <!-- 健康检查对话框 -->
    <a-modal
      v-model:open="healthModalVisible"
      title="插件健康检查"
      :footer="null"
      width="600px"
    >
      <a-descriptions bordered :column="1" v-if="currentHealth">
        <a-descriptions-item label="状态">
          <a-tag :color="currentHealth.status === 'healthy' ? 'green' : 'red'">
            {{ currentHealth.status }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="插件名称" v-if="currentHealth.name">
          {{ currentHealth.name }}
        </a-descriptions-item>
        <a-descriptions-item label="版本" v-if="currentHealth.version">
          {{ currentHealth.version }}
        </a-descriptions-item>
        <a-descriptions-item label="描述" v-if="currentHealth.description">
          {{ currentHealth.description }}
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, inject, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { useRoute } from 'vue-router'
import { SyncOutlined, DownOutlined } from '@ant-design/icons-vue'
import * as api from '@/api/request'
import { usePluginLoader } from '@/composables/usePluginLoader'

const route = useRoute()
const { fetchPlugins } = usePluginLoader()

// 从父组件注入activeTopTab
const activeTopTab = inject('activeTopTab', computed(() => 'my'))

// 当前标签页
const activeTab = ref('my')

// 初始化时根据activeTopTab设置当前标签页
const initActiveTab = () => {
  // 根据activeTopTab的值设置当前标签页
  activeTab.value = activeTopTab.value
}

// 初始化
initActiveTab()

// 监听activeTopTab变化
watch(activeTopTab, (newValue) => {
  // 当顶部标签栏变化时，更新当前标签页
  activeTab.value = newValue
})

// 暴露activeTab给父组件
defineExpose({
  activeTab
})


// 搜索和筛选
const searchKeyword = ref('')
const categoryFilter = ref('')

// 插件市场IP配置
const marketIp = ref('http://192.168.0.201:8081')
const syncing = ref(false)

// 保存插件市场IP配置
const saveMarketIp = () => {
  if (!marketIp.value.trim()) {
    message.warning('请输入插件市场IP地址')
    return
  }
  // 保存IP到localStorage
  localStorage.setItem('pluginMarketIp', marketIp.value)
  message.success('插件市场IP配置成功')
}

// 同步插件市场插件列表
const syncMarketPlugins = async () => {
  let ip = marketIp.value || localStorage.getItem('pluginMarketIp') || 'http://192.168.0.201:8081'
  
  // 如果IP包含http://或https://，移除它
  ip = ip.replace(/^https?:\/\//, '')
  
  if (!ip) {
    message.warning('请先配置插件市场IP地址')
    return
  }

  syncing.value = true
  try {
    const baseUrl = `http://${ip}/api/public/plugins`
    const response = await fetch(`${baseUrl}?page=1&pageSize=10&code=plugin1&versionType=normal`)
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    const data = await response.json()
    
    if (data && data.result && Array.isArray(data.result)) {
      // 更新市场插件列表
      marketPlugins.value = data.result.map((item: any, index: number) => ({
        id: String(item.id || index + 1),
        pluginId: item.code || item.pluginId || '',
        uuid: item.uuid || '',
        name: item.name || '',
        version: item.version || '',
        versionType: item.versionType || '',
        tips: item.tips || '',
        miniClientVersion: item.miniClientVersion || '',
        miniServerVersion: item.miniServerVersion || '',
        installed: false,
        description: item.description || '',
        author: item.author || 'ChProbe Team',
        downloads: item.downloadCount || 0,
        uploadTime: item.uploadedAt || item.createTime || '',
        downloadUrl: item.downloadUrl || ''
      }))
      message.success('同步成功')
    } else {
      message.warning('同步失败：返回数据格式不正确')
    }
  } catch (error: any) {
    message.error('同步失败: ' + error.message)
    console.error('同步插件市场失败:', error)
  } finally {
    syncing.value = false
  }
}

// 我的插件相关数据
const myPluginColumns = [
  { title: '插件ID', dataIndex: 'pluginId', key: 'pluginId' },
  { title: '插件名称', dataIndex: 'name', key: 'name' },
  { title: '版本', dataIndex: 'version', key: 'version' },
  { title: '状态', dataIndex: 'status', key: 'status' },
  { 
    title: '描述', 
    dataIndex: 'description', 
    key: 'description',
    width: 300,
    ellipsis: true,
    slots: {
      customRender: 'description'
    }
  },
  { title: '作者', dataIndex: 'author', key: 'author' },
  { title: '安装时间', dataIndex: 'installTime', key: 'installTime' },
  {
    title: '操作', 
    dataIndex: 'action', 
    key: 'action', 
    fixed: 'right',
    width: 100
  }
]

// 我的插件数据，从后端请求
const myPlugins = ref<any[]>([])
const loading = ref(false)

// 市场插件数据
const marketPlugins = ref<any[]>([])

// 下载进度
const downloadProgress = ref<Record<string, number>>({})

// 路由查看相关
const routesModalVisible = ref(false)
const currentRoutes = ref<any[]>([])
const routeColumns = [
  { title: '路径', dataIndex: 'path', key: 'path' },
  { title: '方法', dataIndex: 'method', key: 'method' },
  { title: '描述', dataIndex: 'description', key: 'description' },
  { title: '处理器', dataIndex: 'handler', key: 'handler' }
]

// 健康检查相关
const healthModalVisible = ref(false)
const currentHealth = ref<any>(null)

// 格式化时间
const formatTime = (time: string) => {
  if (!time) return ''
  try {
    const date = new Date(time)
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  } catch {
    return time
  }
}

// 计算插件市场中每个插件的安装状态
const processedMarketPlugins = computed(() => {
  return marketPlugins.value.map(plugin => {
    const installedPlugin = myPlugins.value.find(p => p.pluginId === plugin.pluginId)
    return {
      ...plugin,
      installed: !!installedPlugin,
      status: installedPlugin?.status || 'not_installed'
    }
  })
})

// 从后端获取我的插件列表
const getMyPlugins = async () => {
  try {
    loading.value = true
    const { result } = await api.get('/plugin-manager/list')
    myPlugins.value = result.map((plugin: any) => ({
      ...plugin,
      loading: false,
      status: plugin.isRunning ? 'running' : 'stopped'
    }))
  } catch (error: any) {
    message.error('获取插件列表失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 启动插件
const startPlugin = async (pluginId: string) => {
  try {
    const plugin = myPlugins.value.find(p => p.pluginId === pluginId)
    if (plugin) {
      plugin.loading = true
    }
    
    await api.post('/plugin-manager/start', {
      pluginId: pluginId,
      command: `/opt/chprobe/plugins/${pluginId}/plugin`,
      args: [],
      config: {}
    })
    
    message.success('插件启动成功')
    // 重新加载整个页面
    window.location.reload()
  } catch (error: any) {
    message.error('插件启动失败: ' + error.message)
    const plugin = myPlugins.value.find(p => p.pluginId === pluginId)
    if (plugin) {
      plugin.loading = false
    }
  }
}

// 停止插件
const stopPlugin = async (pluginId: string) => {
  try {
    const plugin = myPlugins.value.find(p => p.pluginId === pluginId)
    if (plugin) {
      plugin.loading = true
    }
    
    await api.post('/plugin-manager/stop', {
      pluginId: pluginId
    })
    
    message.success('插件停止成功')
    // 重新加载整个页面
    window.location.reload()
  } catch (error: any) {
    message.error('插件停止失败: ' + error.message)
    const plugin = myPlugins.value.find(p => p.pluginId === pluginId)
    if (plugin) {
      plugin.loading = false
    }
  }
}

// 删除插件
const deletePlugin = async (pluginId: string) => {
  try {
    const plugin = myPlugins.value.find(p => p.pluginId === pluginId)
    if (plugin) {
      plugin.loading = true
    }
    
    await api.del(`/plugin/uninstall/${pluginId}`)
    
    message.success('插件删除成功')
    // 重新加载整个页面
    window.location.reload()
  } catch (error: any) {
    message.error('插件删除失败: ' + error.message)
    const plugin = myPlugins.value.find(p => p.pluginId === pluginId)
    if (plugin) {
      plugin.loading = false
    }
  }
}

// 重启插件
const restartPlugin = async (pluginId: string) => {
  try {
    const plugin = myPlugins.value.find(p => p.pluginId === pluginId)
    if (plugin) {
      plugin.loading = true
    }
    
    await api.post('/plugin-manager/restart', {
      pluginId: pluginId,
      command: `/opt/chprobe/plugins/${pluginId}/plugin`,
      args: [],
      config: {}
    })
    
    message.success('插件重启成功')
    await getMyPlugins()
  } catch (error: any) {
    message.error('插件重启失败: ' + error.message)
  } finally {
    const plugin = myPlugins.value.find(p => p.pluginId === pluginId)
    if (plugin) {
      plugin.loading = false
    }
  }
}

// 查看路由
const viewRoutes = async (pluginId: string) => {
  try {
    const { result } = await api.get(`/plugin-manager/routes?pluginId=${pluginId}`)
    currentRoutes.value = result
    routesModalVisible.value = true
  } catch (error: any) {
    message.error('获取插件路由失败: ' + error.message)
  }
}

// 健康检查
const healthCheck = async (pluginId: string) => {
  try {
    const { result } = await api.get(`/plugin-manager/health?pluginId=${pluginId}`)
    currentHealth.value = result
    healthModalVisible.value = true
  } catch (error: any) {
    message.error('健康检查失败: ' + error.message)
  }
}

// 根据插件状态返回对应的颜色
const getStatusColor = (status: string) => {
  switch (status) {
    case 'running':
      return 'green';
    case 'stopped':
      return 'red';
    case 'disabled':
      return 'gray';
    default:
      return 'gray';
  }
}

// 根据插件状态返回对应的文本
const getStatusText = (status: string) => {
  switch (status) {
    case 'running':
      return '运行中';
    case 'stopped':
      return '已停止';
    case 'disabled':
      return '已禁用';
    default:
      return '未知状态';
  }
}

// 根据HTTP方法返回对应的颜色
const getMethodColor = (method: string) => {
  switch (method.toUpperCase()) {
    case 'GET':
      return 'green';
    case 'POST':
      return 'blue';
    case 'PUT':
      return 'orange';
    case 'DELETE':
      return 'red';
    default:
      return 'gray';
  }
}

// 安装插件
const installPlugin = async (pluginId: string) => {
  const plugin = marketPlugins.value.find(p => p.uuid === pluginId)
  if (!plugin) {
    message.error('插件不存在')
    return
  }

  try {
    // 开始下载，设置下载进度为0
    downloadProgress.value[pluginId] = 0

    // 通过后端API创建下载任务
    const createTaskResponse = await api.post('/plugin/download/task', {
      uuid: plugin.uuid,
      pluginId: plugin.uuid,
      pluginName: plugin.name,
      version: plugin.version,
      author: plugin.author,
      description: plugin.description,
      downloadUrl: plugin.downloadUrl
    })

    const taskId = createTaskResponse.result.taskId
    if (!taskId) {
      throw new Error('创建下载任务失败: 未获取到任务ID')
    }

    // 3秒检查一次任务进度（通过后端API）
    const checkStatusInterval = setInterval(async () => {
      try {
        const statusResponse = await api.get('/plugin/download/status/' + taskId)
        
        // 更新下载进度
        downloadProgress.value[pluginId] = parseFloat(statusResponse.result.progress.toFixed(2))

        // 检查下载是否完成
        if (statusResponse.result.status === 'completed') {
          clearInterval(checkStatusInterval)
          // 下载完成，设置进度为100
          downloadProgress.value[pluginId] = 100
          
          // 延迟1秒后清除进度，显示安装成功
          setTimeout(() => {
            delete downloadProgress.value[pluginId]
            message.success('插件下载成功')
            // 重新获取我的插件列表
            getMyPlugins()
          }, 1000)
        } else if (statusResponse.result.status === 'failed') {
          clearInterval(checkStatusInterval)
          delete downloadProgress.value[pluginId]
          message.error('插件下载失败: ' + (statusResponse.result.error || '未知错误'))
        }
      } catch (error: any) {
        console.error('查询下载状态失败:', error)
        clearInterval(checkStatusInterval)
        delete downloadProgress.value[pluginId]
        message.error('查询下载状态失败: ' + error.message)
      }
    }, 3000) // 3秒检查一次任务进度
  } catch (error: any) {
    message.error('插件下载失败: ' + error.message)
    // 清除下载进度
    delete downloadProgress.value[pluginId]
  }
}

// 初始化时获取我的插件列表
onMounted(() => {
  getMyPlugins()
  // 从localStorage加载已保存的插件市场IP
  const savedIp = localStorage.getItem('pluginMarketIp')
  if (savedIp) {
    marketIp.value = savedIp
  }
})

// 监听标签页变化，当切换到我的插件或插件市场标签时重新获取数据
watch(activeTab, (newValue) => {
  if (newValue === 'my' || newValue === 'market') {
    getMyPlugins()
  }
})

</script>

<style scoped>
/* 可根据需要添加样式 */
</style>
