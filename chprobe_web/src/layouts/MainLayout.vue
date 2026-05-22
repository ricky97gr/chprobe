<template>
  <div class="main-layout">
    <!-- 后端服务降级时的加载层 -->
    <div v-if="isBackendDown" class="backend-down-overlay">
      <div class="backend-down-content">
        <a-spin size="large" />
        <div class="loading-text">系统重启中</div>
      </div>
    </div>

    <!-- 左侧导航栏 -->
    <div class="sider-container" :class="{ collapsed: collapsed }">
      <div class="logo">
        <span v-if="!collapsed">Chprobe 管理系统</span>
        <span v-else>CP</span>
      </div>
      <div class="menu-container">
        <a-menu
          :selected-keys="[selectedKey]"
          :open-keys="openKeys"
          :inline-collapsed="collapsed"
          mode="inline"
          theme="light"
          @click="handleMenuClick"
          @openChange="handleOpenChange"
        >
        <a-menu-item key="/dashboard">
          <template #icon>
            <DashboardOutlined />
          </template>
          <span>仪表盘</span>
        </a-menu-item>

        <!-- 插件菜单 - 作为一级菜单显示，放在仪表盘下面 -->
        <template v-for="plugin in pluginMenuItems" :key="plugin.path">
          <!-- 支持子菜单的插件 -->
          <a-sub-menu v-if="plugin.children && plugin.children.length > 0" :key="`sub-${plugin.path}`">
            <template #icon>
              <component :is="getIcon(plugin.meta?.icon as string)" />
            </template>
            <template #title>
              <span>{{ plugin.meta?.title }}</span>
            </template>
            <a-menu-item 
              v-for="child in plugin.children" 
              :key="child.path"
              @click="handlePluginMenuClick(child.path)"
            >
              <template #icon v-if="child.meta?.icon">
                <component :is="getIcon(child.meta?.icon as string)" />
              </template>
              <span>{{ child.meta?.title }}</span>
            </a-menu-item>
          </a-sub-menu>
          <!-- 没有子菜单的插件 -->
          <a-menu-item v-else :key="`item-${plugin.path}`">
            <template #icon>
              <component :is="getIcon(plugin.meta?.icon as string)" />
            </template>
            <span>{{ plugin.meta?.title }}</span>
          </a-menu-item>
        </template>

        <template v-for="menu in menuRoutes" :key="menu.path">
          <a-sub-menu v-if="menu.children && menu.children.length > 0" :key="menu.path">
            <template #icon>
              <component :is="getIcon(menu.meta?.icon as string)" />
            </template>
            <template #title>
              <span>{{ menu.meta?.title }}</span>
            </template>
            <a-menu-item v-for="child in menu.children" :key="`/${menu.path}/${child.path}`">
              {{ child.meta?.title }}
            </a-menu-item>
          </a-sub-menu>
          <a-menu-item v-else :key="`/${menu.path}`">
            <template #icon>
              <component :is="getIcon(menu.meta?.icon as string)" />
            </template>
            <span>{{ menu.meta?.title }}</span>
          </a-menu-item>
        </template>

        <a-sub-menu key="settings">
          <template #icon>
            <SettingOutlined />
          </template>
          <template #title>
            <span>系统设置</span>
          </template>
          <a-menu-item key="/settings">基本设置</a-menu-item>
          <a-menu-item key="/settings/user">用户管理</a-menu-item>
          <a-menu-item key="/settings/system-info">系统信息</a-menu-item>
          <a-menu-item key="/settings/logs">日志管理</a-menu-item>
          <a-menu-item key="/settings/auth">授权管理</a-menu-item>
        </a-sub-menu>
        </a-menu>
      </div>
      <div class="sider-footer">
        <MenuFoldOutlined
          v-if="!collapsed"
          class="trigger"
          @click="toggleCollapsed"
        />
        <MenuUnfoldOutlined
          v-else
          class="trigger"
          @click="toggleCollapsed"
        />
      </div>
    </div>

    <!-- 右侧内容区域 -->
    <div class="content-wrapper" :class="{ collapsed: collapsed }">
      <div class="header">
        <div style="flex: 1; display: flex; align-items: center;">
          <!-- 页面标题 (当没有tab时显示) -->
          <div v-if="currentPageTitle && !isAuthPage && !isLogsPage && !isUpgradePage && !isPluginPage" class="page-title-container">
            <span class="page-title">{{ currentPageTitle }}</span>
          </div>
          <!-- 授权管理页面的tab导航栏 -->
          <div v-if="isAuthPage" class="top-nav" style="height: 40px; display: flex; align-items: center;">
            <a-tabs v-model:active-key="activeTopTab" @change="handleTopTabChange" size="small" class="auth-tabs">
              <a-tab-pane key="info" tab="授权信息" />
              <a-tab-pane key="apply" tab="授权申请" />
            </a-tabs>
          </div>
          <!-- 日志管理页面的tab导航栏 -->
          <div v-else-if="isLogsPage" class="top-nav" style="height: 40px; display: flex; align-items: center;">
            <a-tabs v-model:active-key="activeTopTab" @change="handleTopTabChange" size="small" class="logs-tabs">
              <a-tab-pane key="run" tab="运行日志" />
              <a-tab-pane key="operation" tab="操作日志" />
              <a-tab-pane key="access" tab="访问日志" />
              <a-tab-pane key="audit" tab="审计日志" />
            </a-tabs>
          </div>
          <!-- 升级管理页面的tab导航栏 -->
          <div v-else-if="isUpgradePage" class="top-nav" style="height: 40px; display: flex; align-items: center;">
            <a-tabs v-model:active-key="activeTopTab" @change="handleTopTabChange" size="small" class="upgrade-tabs">
              <a-tab-pane key="client" tab="客户端升级" />
              <a-tab-pane key="server" tab="服务端升级" />
            </a-tabs>
          </div>
          <!-- 插件管理页面的tab导航栏 -->
          <div v-else-if="isPluginPage" class="top-nav" style="height: 40px; display: flex; align-items: center;">
            <a-tabs v-model:active-key="activeTopTab" @change="handleTopTabChange" size="small" class="plugin-tabs">
              <a-tab-pane key="my" tab="我的插件" />
              <a-tab-pane key="market" tab="插件市场" />
            </a-tabs>
          </div>
        </div>
        <div class="user-info">
          <a-dropdown>
            <a class="ant-dropdown-link" @click.prevent>
              <UserOutlined /> {{ currentUser.name }}
              <DownOutlined />
            </a>
            <template #overlay>
              <a-menu @click="handleMenuClick">
                <a-menu-item key="profile">
                  <UserOutlined /> 个人中心
                </a-menu-item>
                <a-menu-item key="settings">
                  <SettingOutlined /> 个人设置
                </a-menu-item>
                <a-menu-divider />
                <a-menu-item key="logout">
                  <LogoutOutlined /> 退出登录
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </div>
      <div class="content-area">
        <div style="background: #fff; margin: 0; padding: 10px; height: 100%; border-radius: 8px;">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, provide, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  DashboardOutlined,
  CloudServerOutlined,
  ContainerOutlined,
  ApiOutlined,
  AlertOutlined,
  SettingOutlined,
  DownOutlined,
  LogoutOutlined,
  UserOutlined,
  CloudUploadOutlined,
  AppstoreOutlined
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useBackendHealth } from '@/composables/useBackendHealth'
import { usePluginLoader } from '@/composables/usePluginLoader'
import { eventBus } from '@/utils/eventBus'

const router = useRouter()
const route = useRoute()

const { isBackendDown } = useBackendHealth()
const { plugins, loading: pluginLoading, fetchPlugins, pluginMenuItems } = usePluginLoader()

// 组件挂载时加载插件
onMounted(() => {
  fetchPlugins()
  
  // 监听插件状态变化事件
  eventBus.on('plugin:status:changed', handlePluginStatusChanged)
  eventBus.on('plugin:loaded', handlePluginLoaded)
})

// 组件卸载时移除事件监听器
onUnmounted(() => {
  eventBus.off('plugin:status:changed', handlePluginStatusChanged)
  eventBus.off('plugin:loaded', handlePluginLoaded)
})

// 处理插件状态变化
const handlePluginStatusChanged = async () => {
  console.log('插件状态变化，重新获取插件列表')
  await fetchPlugins()
}

// 处理插件加载完成
const handlePluginLoaded = () => {
  console.log('插件加载完成')
}

const collapsed = ref(false)
const selectedKey = computed(() => route.path)
const openKeys = ref<string[]>([])

const currentUser = ref({
  name: '管理员'
})

// 当前激活的顶部标签页
// 插件管理页面的当前标签页状态
const pluginTab = ref('my')

// 升级管理页面的当前标签页状态
const upgradeTab = ref('client')

const activeTopTab = computed(() => {
  // 授权管理页面
  if (route.path === '/settings/auth') {
    return 'info'
  } else if (route.path === '/settings/auth/apply') {
    return 'apply'
  }
  // 日志管理页面
  else if (route.path === '/settings/logs') {
    return 'run'
  } else if (route.path === '/settings/logs/operation') {
    return 'operation'
  } else if (route.path === '/settings/logs/access') {
    return 'access'
  } else if (route.path === '/settings/logs/audit') {
    return 'audit'
  }
  // 升级管理页面
  else if (route.path === '/upgrade-management') {
    return upgradeTab.value
  }
  // 插件管理页面
  else if (route.path === '/plugin-management') {
    return pluginTab.value // 返回当前插件管理页面的标签页
  }
  // 系统信息页面
  else if (route.path === '/settings/system-info') {
    return 'system-info'
  }
  // 基本设置和用户管理页面
  else if (route.path === '/settings' || route.path === '/settings/user') {
    return 'basic'
  }
  return 'basic'
})

// 提供activeTopTab给子组件
provide('activeTopTab', activeTopTab)

// 判断是否在授权管理页面
const isAuthPage = computed(() => {
  return route.path.startsWith('/settings/auth')
})

// 判断是否在日志管理页面
const isLogsPage = computed(() => {
  return route.path.startsWith('/settings/logs')
})

// 判断是否在升级管理页面
const isUpgradePage = computed(() => {
  return route.path === '/upgrade-management'
})

// 判断是否在插件管理页面
const isPluginPage = computed(() => {
  return route.path.startsWith('/plugin-management')
})

// 当前页面标题
const currentPageTitle = computed(() => {
  const path = route.path
  if (path === '/dashboard') return '仪表盘'
  if (path === '/container/list') return '容器列表'
  if (path === '/container/image') return '镜像列表'
  if (path === '/upgrade-management') return '升级管理'
  if (path === '/plugin-management') return '插件管理'
  if (path === '/monitor/alarm-policy') return '告警策略'
  if (path === '/monitor/alarm-event') return '告警事件'
  if (path === '/monitor/metrics') return '监控指标'
  if (path === '/settings') return '基本设置'
  if (path === '/settings/user') return '用户管理'
  if (path === '/settings/system-info') return '系统信息'
  // 插件页面标题
  if (path.startsWith('/plugin/')) {
    const plugin = plugins.value.find(p => path.includes(p.pluginId))
    return plugin?.name || '插件页面'
  }
  return ''
})

// 图标映射
const iconMap: Record<string, any> = {
  DashboardOutlined,
  CloudServerOutlined,
  ContainerOutlined,
  ApiOutlined,
  AlertOutlined,
  SettingOutlined,
  CloudUploadOutlined,
  AppstoreOutlined
}

// 获取图标组件
const getIcon = (iconName: string) => {
  return iconMap[iconName] || null
}

// 获取路由配置中的菜单项
const menuRoutes = computed(() => {
  const routes = router.getRoutes()
  const layoutRoute = routes.find(r => r.name === 'Layout')
  if (!layoutRoute || !layoutRoute.children) return []
  
  return layoutRoute.children
    .filter(route => !route.meta?.hidden && route.path !== 'dashboard' && route.path !== 'settings')
    .map(route => ({
      ...route,
      children: route.children?.filter(child => !child.meta?.hidden)
    }))
    .filter(route => !route.children || route.children.length > 0)
})

const toggleCollapsed = () => {
  collapsed.value = !collapsed.value
}

const onCollapse = (value: boolean) => {
  collapsed.value = value
}

const handleMenuClick = ({ key }: { key: string }) => {
    console.log('Menu item clicked:', key)
    if (key === 'logout') {
      // 退出登录处理
      localStorage.removeItem('token')
      message.success('退出登录成功')
      router.push('/login')
    } else {
      console.log('Navigating to:', key)
      router.push(key)
    }
  }

  // 处理插件子菜单点击
  const handlePluginMenuClick = (path: string) => {
    console.log('Plugin menu item clicked:', path)
    router.push(path)
  }

// 处理顶部标签页切换
const handleTopTabChange = (key: string) => {
  // 根据标签页切换路由
  // 授权管理页面
  if (key === 'info') {
    router.push('/settings/auth')
  } else if (key === 'apply') {
    router.push('/settings/auth/apply')
  }
  // 日志管理页面
  else if (key === 'run') {
    router.push('/settings/logs')
  } else if (key === 'operation') {
    router.push('/settings/logs/operation')
  } else if (key === 'access') {
    router.push('/settings/logs/access')
  } else if (key === 'audit') {
    router.push('/settings/logs/audit')
  }
  // 升级管理页面
  else if (key === 'client' || key === 'server') {
    // 更新升级管理页面的标签页状态
    upgradeTab.value = key
  }
  // 插件管理页面
  else if (key === 'my' || key === 'market') {
    // 更新插件管理页面的标签页状态
    pluginTab.value = key
  }
  // 系统信息页面和基本设置页面
  else if (key === 'system-info' || key === 'basic') {
    // 这些页面没有标签页，不需要路由切换
  }
}

const handleOpenChange = (keys: string[]) => {
  openKeys.value = keys
}
</script>

<style scoped>
.main-layout {
  height: 100vh;
  overflow: hidden;
  display: flex;
  width: 100vw;
  margin: 0;
  padding: 0;
}

/* 左侧固定导航栏 */
.sider-container {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  width: 200px;
  z-index: 10;
  overflow: hidden;
  transition: width 0.3s;
  background: #fff;
  border-right: 1px solid #f0f0f0;
}

.sider-container.collapsed {
  width: 80px;
}

/* 右侧内容区域 */
.content-wrapper {
  flex: 1;
  margin-left: 200px;
  transition: margin-left 0.3s;
  height: 100vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.content-wrapper.collapsed {
  margin-left: 80px;
}

/* 顶部导航栏 */
.header {
  padding: 0 16px;
  background: #fff;
  display: flex;
  align-items: center;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  height: 60px;
  line-height: 60px;
  flex-shrink: 0;
}

/* 内容区域 */
.content-area {
  flex: 1;
  margin: 12px;
  padding: 0;
  background: #f0f2f5;
  overflow-y: auto;
  overflow-x: hidden;
  border-radius: 8px;
}

.content-area::-webkit-scrollbar {
  width: 8px;
}

.content-area::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.content-area::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

.content-area::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

.logo {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: bold;
  color: #1890ff;
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
}

.trigger {
  font-size: 14px;
  padding: 0 16px;
  cursor: pointer;
  transition: color 0.3s;
}

.trigger:hover {
  color: #1890ff;
}

.menu-container {
  height: calc(100vh - 112px);
  overflow-y: auto;
  overflow-x: hidden;
}

.sider-footer {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  padding: 12px 0;
  background: #fff;
  border-top: 1px solid #f0f0f0;
}

.menu-container::-webkit-scrollbar {
  width: 6px;
}

.menu-container::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.menu-container::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.menu-container::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

.user-info {
  display: flex;
  align-items: center;
}

.page-title-container {
  margin-right: 24px;
}

.page-title {
  font-size: 16px;
  font-weight: bold;
  color: #333;
}

.top-nav {
  display: flex;
  align-items: center;
}

:deep(.ant-menu-item),
:deep(.ant-menu-submenu-title) {
  font-size: 14px;
  padding-left: 16px !important;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

:deep(.ant-menu-submenu-title > span) {
  flex: 1;
  text-align: left;
}

:deep(.ant-menu-submenu-title > .anticon-down) {
  margin-left: auto;
}

/* 二级菜单文字与一级菜单文字对齐 */
:deep(.ant-menu-submenu > .ant-menu) {
  padding-left: 16px !important;
}

:deep(.ant-menu-submenu > .ant-menu > .ant-menu-item) {
  padding-left: 32px !important;
  margin-left: 0 !important;
}

:deep(.ant-menu-item-group > .ant-menu-item) {
  padding-left: 32px !important;
}

:deep(.ant-tabs-tab) {
  display: flex;
  align-items: center;
  height: 40px;
  padding: 0 10px;
  margin-right: 4px;
  min-width: auto !important;
}

:deep(.ant-tabs-nav) {
  margin: 0 !important;
  padding: 0 4px;
  min-width: auto !important;
  max-width: none !important;
}

:deep(.ant-tabs-nav-list) {
  min-width: auto !important;
  max-width: none !important;
}

:deep(.ant-tabs-nav-wrap) {
  overflow: visible !important;
  min-width: auto !important;
  max-width: none !important;
}

:deep(.ant-tabs) {
  min-width: auto !important;
  max-width: none !important;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.backend-down-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.95);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.backend-down-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
}

.loading-text {
  font-size: 18px;
  color: #333;
  font-weight: 500;
}
</style>
