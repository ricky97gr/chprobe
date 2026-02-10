# Chprobe 管理系统 - 项目结构说明

## 已完成的功能

### 1. 项目基础架构

- ✅ **技术栈配置**
  - Vue 3 + TypeScript
  - Ant Design Vue (UI 组件库)
  - Vue Router 4 (路由管理)
  - Pinia (状态管理)
  - Axios (HTTP 客户端)
  - Vite (构建工具)

- ✅ **路径别名配置**
  - `@` 指向 `src` 目录
  - TypeScript 和 Vite 双配置

- ✅ **环境变量配置**
  - `.env` - 通用环境变量
  - `.env.development` - 开发环境
  - `.env.production` - 生产环境

### 2. 目录结构

```
src/
├── api/              # API 接口层
│   └── user.ts       # 用户接口示例
├── assets/           # 静态资源
├── components/       # 公共组件
├── layouts/          # 布局组件
│   └── MainLayout.vue # 主布局
├── router/           # 路由配置
│   └── index.ts      # 路由定义
├── stores/           # 状态管理
│   ├── app.ts        # 应用状态
│   └── user.ts       # 用户状态
├── types/            # TypeScript 类型
├── utils/            # 工具函数
│   └── request.ts    # Axios 封装
├── views/            # 页面组件
│   ├── dashboard/    # 仪表盘
│   ├── user/         # 用户管理
│   └── settings/     # 系统设置
├── App.vue           # 根组件
└── main.ts           # 入口文件
```

### 3. 核心组件

#### 布局组件 (MainLayout.vue)

- 侧边栏导航（可折叠）
- 顶部导航栏（用户信息、下拉菜单）
- 内容区域（路由视图）
- 响应式设计

**功能特性：**
- 三级菜单支持
- 路由高亮
- 折叠/展开动画
- 用户信息展示

#### 页面组件

**仪表盘 (dashboard/index.vue)**
- 统计卡片（总用户、活跃用户、访问量、系统状态）
- 最近活动列表
- 快速操作按钮

**用户管理 (user/index.vue)**
- 用户列表展示（Ant Design Table）
- 搜索和筛选功能
- 新增/编辑用户（Modal 弹窗）
- 删除用户操作
- 分页功能

**系统设置 (settings/index.vue)**
- 基本设置（系统名称、描述、语言、主题）
- 安全设置（密码修改、登录超时、双因素认证）
- 系统信息展示（Descriptions）

### 4. 状态管理 (Pinia)

**app store**
- 侧边栏折叠状态
- 全局加载状态

**user store**
- 用户信息（姓名、头像、角色）
- 用户操作方法（设置用户信息、退出登录）

### 5. API 请求封装

**request.ts 特性：**
- 统一的请求/响应拦截器
- Token 自动注入（Authorization header）
- 统一错误处理
- 请求超时配置（15秒）
- 401 自动跳转登录

**示例 API (user.ts)**
- getUserList - 获取用户列表
- getUserById - 获取用户详情
- createUser - 创建用户
- updateUser - 更新用户
- deleteUser - 删除用户

### 6. 路由配置

**路由结构：**
```typescript
/
├── /dashboard (仪表盘)
├── /user (用户管理)
└── /settings (系统设置)
```

**路由特性：**
- 懒加载（动态 import）
- 路由元信息（标题、图标）
- 嵌套路由支持

## 如何使用

### 启动开发服务器

```bash
cd chprobe_web
npm install
npm run dev
```

访问：http://localhost:5173

### 构建生产版本

```bash
npm run build
```

### 预览生产构建

```bash
npm run preview
```

## 扩展开发指南

### 添加新页面

1. 在 `views/` 下创建页面组件
2. 在 `router/index.ts` 添加路由配置
3. 在 `layouts/MainLayout.vue` 添加菜单项

**示例：**
```typescript
// router/index.ts
{
  path: 'new-page',
  name: 'NewPage',
  component: () => import('@/views/new-page/index.vue'),
  meta: { title: '新页面', icon: 'folder' }
}
```

### 添加新的 API 接口

1. 在 `api/` 下创建接口文件
2. 定义 TypeScript 类型
3. 使用封装的 `request` 工具

**示例：**
```typescript
// api/product.ts
import request from '@/utils/request'

export const getProductList = (params: any) => {
  return request.get('/product/list', { params })
}
```

### 添加新的 Store

1. 在 `stores/` 下创建 store 文件
2. 使用 `defineStore` 定义
3. 在组件中使用

**示例：**
```typescript
// stores/product.ts
import { defineStore } from 'pinia'

export const useProductStore = defineStore('product', {
  state: () => ({ products: [] }),
  actions: {
    setProducts(products: any[]) {
      this.products = products
    }
  }
})
```

## 注意事项

### 1. 构建警告

构建时会出现 chunk size 警告（> 500kB），这是因为 Ant Design Vue 包较大。生产环境可以通过以下方式优化：

- 使用动态导入（懒加载）
- 配置 manualChunks 代码分割
- 按需引入组件（推荐）

### 2. Node.js 版本

Vite 7.x 需要 Node.js 20.19+ 或 22.12+，当前使用 21.7.3 也能正常工作，但建议升级到推荐版本。

### 3. API 地址配置

修改 `.env` 文件中的 `VITE_API_BASE_URL` 来配置后端 API 地址。

## 下一步建议

### 功能扩展

1. **用户认证**
   - 登录/注册页面
   - Token 持久化
   - 路由守卫

2. **权限管理**
   - 角色权限控制
   - 动态菜单
   - 按钮权限

3. **数据可视化**
   - ECharts 图表集成
   - 数据统计面板

4. **国际化**
   - i18n 多语言支持
   - 语言切换功能

5. **主题定制**
   - 深色模式
   - 主题色配置

### 代码优化

1. **组件拆分**
   - 提取公共组件
   - 表单组件封装

2. **类型定义**
   - 完善 TypeScript 类型
   - 接口定义规范

3. **错误处理**
   - 全局错误捕获
   - 友好的错误提示

4. **性能优化**
   - 虚拟滚动（长列表）
   - 图片懒加载
   - 防抖节流

## 参考文档

- [Vue 3 官方文档](https://cn.vuejs.org/)
- [Ant Design Vue](https://antdv.com/)
- [Vue Router](https://router.vuejs.org/zh/)
- [Pinia](https://pinia.vuejs.org/zh/)
- [Vite](https://cn.vitejs.dev/)
- [TypeScript](https://www.typescriptlang.org/zh/)
