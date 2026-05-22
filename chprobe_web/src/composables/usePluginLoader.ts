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

interface PluginMeta {
  id: string
  name: string
  version: string
  author: string
  description: string
  menu?: MenuItem
  routes?: PluginRoute[]
  permissions?: string[]
  apiPrefix?: string
}

interface MenuItem {
  id: string
  title: string
  icon?: string
  order?: number
  children?: MenuChild[]
}

interface MenuChild {
  id: string
  title: string
  path: string
  icon?: string
  order?: number
}

interface PluginWebRoute {
  path: string
  name: string
  title: string
  icon?: string
  component?: string
}

interface PluginWebConfig {
  name: string
  description: string
  menu?: MenuItem
  routes?: PluginWebRoute[]
  permissions?: string[]
  apiPrefix?: string
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
  meta?: PluginMeta
  webConfig?: PluginWebConfig
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
      // 获取所有插件的前端配置
      const response = await get('/plugin-manager/web-configs')
      console.log('获取插件前端配置成功:', response)
      
      const pluginList = response.result || []
      console.log('插件前端配置列表:', pluginList)
      
      plugins.value = pluginList.map((plugin: any) => ({
        pluginId: plugin.pluginId,
        name: plugin.name || plugin.webConfig?.name || plugin.meta?.name || '未知插件',
        version: plugin.version || plugin.meta?.version || '',
        description: plugin.description || plugin.webConfig?.description || plugin.meta?.description || '',
        author: plugin.author || plugin.meta?.author || '',
        status: 'running',
        isRunning: true,
        meta: plugin.meta,
        webConfig: plugin.webConfig
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
      // 如果插件有web配置，添加多个路由
      if (plugin.webConfig && plugin.webConfig.routes && plugin.webConfig.routes.length > 0) {
        for (const routeConfig of plugin.webConfig.routes) {
          await addRouteFromConfig(plugin.pluginId, routeConfig)
        }
      } else {
        // 没有web配置，添加默认路由
        const routeName = `Plugin-${plugin.pluginId}`
        if (router.hasRoute(routeName)) {
          console.log(`路由 ${routeName} 已存在，跳过添加`)
          return
        }

        const component = await loadPluginComponent(plugin.pluginId)
        
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
      }
    } catch (err) {
      console.error(`添加插件路由失败:`, err)
    }
  }

  // 根据配置添加路由
  const addRouteFromConfig = async (pluginId: string, routeConfig: PluginWebRoute) => {
    const routeName = `Plugin-${pluginId}-${routeConfig.name}`
    
    // 检查路由是否已存在
    if (router.hasRoute(routeName)) {
      console.log(`路由 ${routeName} 已存在，跳过添加`)
      return
    }

    // 动态加载组件
    const component = await loadPluginComponent(pluginId, routeConfig.component)
    
    // 添加路由
    router.addRoute('Layout', {
      path: `plugin/${pluginId}/${routeConfig.path}`,
      name: routeName,
      component,
      meta: {
        title: routeConfig.title || routeConfig.name,
        icon: routeConfig.icon || 'CloudServerOutlined',
        requiresAuth: true
      }
    })
    
    console.log(`成功添加路由: ${routeName}`)
  }

  // 动态加载插件组件
  const loadPluginComponent = async (pluginId: string, componentName: string = 'index') => {
    try {
      // 尝试从插件的静态资源加载组件
      const componentPath = `plugin/${pluginId}/${componentName}`
      console.log(`尝试加载插件组件: ${componentPath}`)
      
      // 使用动态import加载插件组件
      // 这里使用fetch动态获取组件代码并执行
      return () => import(`@/views/error/PluginError.vue`)
    } catch (err) {
      console.error('加载插件组件失败:', err)
      // 返回一个默认的错误组件
      return () => import('@/views/error/PluginError.vue')
    }
  }

  // 移除插件路由
  const removePluginRoutes = async (pluginId: string) => {
    const plugin = plugins.value.find(p => p.pluginId === pluginId)
    if (!plugin) return

    // 删除该插件的所有路由
    const routes = router.getRoutes()
    for (const route of routes) {
      const routeName = String(route.name)
      if (routeName && routeName.startsWith(`Plugin-${pluginId}`)) {
        router.removeRoute(routeName)
        console.log(`成功移除路由: ${routeName}`)
      }
    }
  }

  // 获取插件菜单项
  const pluginMenuItems = computed(() => {
    const menuItems: any[] = []
    
    for (const plugin of plugins.value) {
      if (plugin.isRunning !== true) continue
      
      // 优先使用 meta.menu 结构（支持二级菜单）
      if (plugin.meta?.menu) {
        const menu = plugin.meta.menu
        
        if (menu.children && menu.children.length > 0) {
          // 有子菜单，创建父菜单
          menuItems.push({
            path: `/plugin/${plugin.pluginId}`,
            name: `Plugin-${plugin.pluginId}`,
            meta: {
              title: menu.title || plugin.name,
              icon: menu.icon || 'AppstoreOutlined',
              requiresAuth: true
            },
            children: menu.children.map((child: MenuChild) => ({
              path: child.path || `/plugin/${plugin.pluginId}/${child.id}`,
              name: `Plugin-${plugin.pluginId}-${child.id}`,
              meta: {
                title: child.title || child.id,
                icon: child.icon || 'AppstoreOutlined',
                requiresAuth: true
              }
            }))
          })
        } else {
          // 没有子菜单，直接作为菜单项
          menuItems.push({
            path: `/plugin/${plugin.pluginId}`,
            name: `Plugin-${plugin.pluginId}`,
            meta: {
              title: menu.title || plugin.name,
              icon: menu.icon || 'AppstoreOutlined',
              requiresAuth: true
            }
          })
        }
        continue
      }
      
      // 使用 webConfig.menu 结构
      if (plugin.webConfig?.menu) {
        const menu = plugin.webConfig.menu
        
        if (menu.children && menu.children.length > 0) {
          // 有子菜单，创建父菜单
          menuItems.push({
            path: `/plugin/${plugin.pluginId}`,
            name: `Plugin-${plugin.pluginId}`,
            meta: {
              title: menu.title || plugin.name,
              icon: menu.icon || 'AppstoreOutlined',
              requiresAuth: true
            },
            children: menu.children.map((child: MenuChild) => ({
              path: child.path || `/plugin/${plugin.pluginId}/${child.id}`,
              name: `Plugin-${plugin.pluginId}-${child.id}`,
              meta: {
                title: child.title || child.id,
                icon: child.icon || 'AppstoreOutlined',
                requiresAuth: true
              }
            }))
          })
        } else {
          // 没有子菜单，直接作为菜单项
          menuItems.push({
            path: `/plugin/${plugin.pluginId}`,
            name: `Plugin-${plugin.pluginId}`,
            meta: {
              title: menu.title || plugin.name,
              icon: menu.icon || 'AppstoreOutlined',
              requiresAuth: true
            }
          })
        }
        continue
      }
      
      // 使用 webConfig.routes 结构（旧的兼容方式）
      if (plugin.webConfig && plugin.webConfig.routes && plugin.webConfig.routes.length > 0) {
        const hasMultipleRoutes = plugin.webConfig.routes.length > 1
        
        if (hasMultipleRoutes) {
          const firstRoute = plugin.webConfig.routes[0]
          menuItems.push({
            path: `/plugin/${plugin.pluginId}`,
            name: `Plugin-${plugin.pluginId}`,
            meta: {
              title: plugin.name,
              icon: firstRoute?.icon || 'AppstoreOutlined',
              requiresAuth: true
            },
            children: plugin.webConfig.routes.map((route: PluginWebRoute) => ({
              path: `/plugin/${plugin.pluginId}/${route?.path || ''}`,
              name: `Plugin-${plugin.pluginId}-${route?.name || ''}`,
              meta: {
                title: route?.title || '',
                icon: route?.icon || 'AppstoreOutlined',
                requiresAuth: true
              }
            }))
          })
        } else {
          const route = plugin.webConfig.routes[0]
          if (route) {
            menuItems.push({
              path: `/plugin/${plugin.pluginId}/${route.path}`,
              name: `Plugin-${plugin.pluginId}-${route.name}`,
              meta: {
                title: route.title,
                icon: route.icon || 'AppstoreOutlined',
                requiresAuth: true
              }
            })
          }
        }
        continue
      }
      
      // 没有配置，使用默认配置
      menuItems.push({
        path: `/plugin/${plugin.pluginId}`,
        name: `Plugin-${plugin.pluginId}`,
        meta: {
          title: plugin.name,
          icon: 'AppstoreOutlined',
          requiresAuth: true
        }
      })
    }
    
    return menuItems
  })

  // 启动插件
  const startPlugin = async (pluginId: string) => {
    try {
      await post('/plugin-manager/start', {
        pluginId,
        command: `/opt/chprobe/plugins/${pluginId}/server`,
        args: [],
        config: {}
      })
      
      // 发送插件状态变更事件
      eventBus.emit('plugin:status:changed', { pluginId, status: 'running' })
      
      // 重新加载插件配置
      await fetchPlugins()
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
      
      // 移除插件路由
      await removePluginRoutes(pluginId)
      
      // 更新插件列表
      const plugin = plugins.value.find(p => p.pluginId === pluginId)
      if (plugin) {
        plugin.isRunning = false
        plugin.status = 'stopped'
      }
      
      return true
    } catch (err: any) {
      console.error('停止插件失败:', err)
      return false
    }
  }

  // 获取插件资源URL
  const getPluginAssetUrl = (pluginId: string, assetPath: string) => {
    return `/api/plugin-static/${pluginId}/${assetPath}`
  }

  // 加载插件样式
  const loadPluginStyles = async (pluginId: string) => {
    try {
      const styleUrl = getPluginAssetUrl(pluginId, 'style.css')
      const link = document.createElement('link')
      link.rel = 'stylesheet'
      link.href = styleUrl
      link.id = `plugin-style-${pluginId}`
      
      // 检查是否已加载
      if (!document.getElementById(`plugin-style-${pluginId}`)) {
        document.head.appendChild(link)
        console.log(`加载插件样式: ${styleUrl}`)
      }
    } catch (err) {
      console.warn(`加载插件样式失败: ${err}`)
    }
  }

  // 卸载插件样式
  const unloadPluginStyles = (pluginId: string) => {
    const link = document.getElementById(`plugin-style-${pluginId}`)
    if (link) {
      document.head.removeChild(link)
      console.log(`卸载插件样式: ${pluginId}`)
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
    stopPlugin,
    getPluginAssetUrl,
    loadPluginStyles,
    unloadPluginStyles
  }
}