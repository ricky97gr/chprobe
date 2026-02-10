import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/login/index.vue'),
      meta: { title: '登录', requiresAuth: false }
    },
    {
      path: '/',
      name: 'Layout',
      component: () => import('@/layouts/MainLayout.vue'),
      redirect: '/dashboard',
      meta: { requiresAuth: true },
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('@/views/dashboard/index.vue'),
          meta: { title: '仪表盘', icon: 'DashboardOutlined' }
        },
        {
          path: 'container',
          name: 'Container',
          meta: { title: '容器管理', icon: 'ContainerOutlined', hidden: true },
          children: [
            {
              path: 'list',
              name: 'ContainerList',
              component: () => import('@/views/container/list.vue'),
              meta: { title: '容器列表' }
            },
            {
              path: 'image',
              name: 'ContainerImage',
              component: () => import('@/views/container/image.vue'),
              meta: { title: '镜像列表' }
            }
          ]
        },
        {
          path: 'upgrade-management',
          name: 'UpgradeManagement',
          component: () => import('@/views/system/upgrade-management.vue'),
          meta: { title: '升级管理', icon: 'CloudUploadOutlined' }
        },
        {
          path: 'plugin-management',
          name: 'PluginManagement',
          component: () => import('@/views/system/plugin-management.vue'),
          meta: { title: '插件管理', icon: 'AppstoreOutlined' }
        },
        {
          path: 'monitor',
          name: 'Monitor',
          meta: { title: '系统监控', icon: 'AlertOutlined', hidden: true },
          children: [
            {
              path: 'alarm-policy',
              name: 'AlarmPolicy',
              component: () => import('@/views/monitor/alarm-policy.vue'),
              meta: { title: '告警策略' }
            },
            {
              path: 'alarm-event',
              name: 'AlarmEvent',
              component: () => import('@/views/monitor/alarm-event.vue'),
              meta: { title: '告警事件' }
            },
            {
              path: 'metrics',
              name: 'MonitorMetrics',
              component: () => import('@/views/monitor/metrics.vue'),
              meta: { title: '监控指标' }
            }
          ]
        },
        {
          path: 'settings',
          name: 'Settings',
          meta: { title: '系统设置', icon: 'SettingOutlined' },
          children: [
            {
              path: '',
              name: 'BasicSettings',
              component: () => import('@/views/settings/index.vue'),
              meta: { title: '基本设置' }
            },
            {
              path: 'user',
              name: 'UserManagement',
              component: () => import('@/views/user/index.vue'),
              meta: { title: '用户管理' }
            },
            {
              path: 'logs',
              name: 'LogManagement',
              meta: { title: '日志管理' },
              children: [
                {
                  path: '',
                  name: 'RunLogs',
                  component: () => import('@/views/settings/logs.vue'),
                  meta: { title: '运行日志' }
                },
                {
                  path: 'operation',
                  name: 'OperationLogs',
                  component: () => import('@/views/settings/logs/operation.vue'),
                  meta: { title: '操作日志' }
                },
                {
                  path: 'access',
                  name: 'AccessLogs',
                  component: () => import('@/views/settings/logs/access.vue'),
                  meta: { title: '访问日志' }
                },
                {
                  path: 'audit',
                  name: 'AuditLogs',
                  component: () => import('@/views/settings/logs/audit.vue'),
                  meta: { title: '审计日志' }
                }
              ]
            },
            {
              path: 'auth',
              name: 'AuthManagement',
              meta: { title: '授权管理' },
              children: [
                {
                  path: '',
                  name: 'AuthInfo',
                  component: () => import('@/views/settings/auth.vue'),
                  meta: { title: '授权信息' }
                },
                {
                  path: 'apply',
                  name: 'AuthApply',
                  component: () => import('@/views/settings/auth/apply.vue'),
                  meta: { title: '授权申请' }
                }
              ]
            },
            {
              path: 'system-info',
              name: 'SystemInfo',
              component: () => import('@/views/settings/system-info.vue'),
              meta: { title: '系统信息' }
            }
          ]
        }
      ]
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - ChProbe`
  }
  
  // 检查是否需要认证
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth !== false)
  
  // 模拟检查用户是否已登录
  const isLoggedIn = localStorage.getItem('token') !== null
  
  if (requiresAuth && !isLoggedIn) {
    // 未登录，重定向到登录页面
    next({ name: 'Login' })
  } else if (to.name === 'Login' && isLoggedIn) {
    // 已登录，重定向到首页
    next({ name: 'Dashboard' })
  } else {
    next()
  }
})

export default router
