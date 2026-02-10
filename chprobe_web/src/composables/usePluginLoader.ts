import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { get, post } from '@/api/request'
import { eventBus } from '@/utils/eventBus'

interface PluginRoute {
  path: string
  name: string
  component: any
  meta: {
    title: string
    icon?: string
    requiresAuth?: boolean
  }
  children?: PluginRoute[]
}

interface PluginInfo {
  pluginId: string
  name: string
  version: string
  description: string
  author: string
  status: string
  isRunning?: boolean
  routes?: PluginRoute[]
}

export function usePluginLoader() {
  const router = useRouter()
  const plugins = ref<PluginInfo[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 获取已安装的插件列表
  const fetchPlugins = async () => {
    console.log('开始获取插件列表')
    loading.value = true
    error.value = null
    try {
      const response = await get('/plugin-manager/list')
      console.log('获取插件列表成功:', response)
      const pluginList = response.result || []
      console.log('插件列表:', pluginList)
      
      // 转换数据格式，将PluginID转换为pluginId
      plugins.value = pluginList.map((plugin: any) => ({
        pluginId: plugin.PluginID || plugin.pluginId,
        name: plugin.Name || plugin.name,
        version: plugin.Version || plugin.version,
        description: plugin.Description || plugin.description,
        author: plugin.Author || plugin.author,
        status: plugin.Status || plugin.status,
        isRunning: plugin.IsRunning || plugin.isRunning,
        routes: plugin.routes
      }))
      console.log('转换后的插件列表:', plugins.value)
      
      // 为每个插件添加路由
      await addPluginRoutes(plugins.value)
      
      // 发送插件加载完成事件
      eventBus.emit('plugin:loaded', 'all')
    } catch (err: any) {
      error.value = err.message || '获取插件列表失败'
      console.error('获取插件列表失败:', err)
    } finally {
      loading.value = false
      console.log('获取插件列表完成')
    }
  }

  // 添加插件路由
  const addPluginRoutes = async (pluginList: PluginInfo[]) => {
    for (const plugin of pluginList) {
      // 只为已启用的插件添加路由
      if (plugin.status === 'enabled' || plugin.isRunning) {
        await addPluginRoute(plugin)
      }
    }
  }

  // 添加单个插件的路由
  const addPluginRoute = async (plugin: PluginInfo) => {
    try {
      // 检查路由是否已存在
      const routeName = `Plugin-${plugin.pluginId}`
      const existingRoute = router.hasRoute(routeName)
      if (existingRoute) {
        console.log(`路由 ${routeName} 已存在，跳过添加`)
        return
      }

      // 动态加载组件
      const component = await loadPluginComponent(plugin.pluginId)
      
      // 添加路由
      router.addRoute('Layout', {
        path: `plugin/${plugin.pluginId}`,
        name: routeName,
        component,
        meta: {
          title: plugin.name,
          icon: 'CloudServerOutlined',
          requiresAuth: true
        }
      })
      
      console.log(`成功添加路由: ${routeName}`)
    } catch (err) {
      console.error(`添加插件路由失败:`, err)
    }
  }

  // 动态加载插件组件
  const loadPluginComponent = async (pluginId: string) => {
    try {
      // 对于host-monitor插件，使用本地的host/list.vue组件
      if (pluginId === 'host-monitor') {
        return () => import('@/views/host/list.vue')
      }
      // 对于其他插件，使用默认的错误组件
      return () => import('@/views/error/PluginError.vue')
    } catch (err) {
      console.error('加载插件组件失败:', err)
      // 返回一个默认的错误组件
      return () => import('@/views/error/PluginError.vue')
    }
  }

  // 移除插件路由
  const removePluginRoutes = async (pluginId: string) => {
    const plugin = plugins.value.find(p => p.pluginId === pluginId)
    if (!plugin || !plugin.routes) return

    for (const route of plugin.routes) {
      if (router.hasRoute(route.name)) {
        router.removeRoute(route.name)
        console.log(`成功移除路由: ${route.name}`)
      }
    }
  }

  // 获取插件菜单项
  const pluginMenuItems = computed(() => {
    return plugins.value
      .filter(plugin => plugin.isRunning === true)
      .map(plugin => ({
        path: `/plugin/${plugin.pluginId}`,
        name: `Plugin-${plugin.pluginId}`,
        meta: {
          title: plugin.name,
          icon: 'AppstoreOutlined',
          requiresAuth: true
        },
        children: plugin.routes?.map(route => ({
          ...route,
          path: `${plugin.pluginId}/${route.path}`,
          name: `Plugin-${plugin.pluginId}-${route.name}`
        })) || []
      }))
  })

  // 启动插件
  const startPlugin = async (pluginId: string) => {
    try {
      await post('/plugin-manager/start', {
        pluginId,
        command: `/opt/chprobe/plugins/${pluginId}/plugin`,
        args: [],
        config: {}
      })
      
      // 发送插件状态变更事件
      eventBus.emit('plugin:status:changed', { pluginId, status: 'running' })
      
      // 重新加载整个页面
      window.location.reload()
      return true
    } catch (err: any) {
      console.error('启动插件失败:', err)
      return false
    }
  }

  // 停止插件
  const stopPlugin = async (pluginId: string) => {
    try {
      await post('/plugin-manager/stop', { pluginId })
      
      // 发送插件状态变更事件
      eventBus.emit('plugin:status:changed', { pluginId, status: 'stopped' })
      
      // 重新加载整个页面
      window.location.reload()
      return true
    } catch (err: any) {
      console.error('停止插件失败:', err)
      return false
    }
  }

  return {
    plugins,
    loading,
    error,
    fetchPlugins,
    addPluginRoutes,
    removePluginRoutes,
    pluginMenuItems,
    startPlugin,
    stopPlugin
  }
}
