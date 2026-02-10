# Chprobe 管理系统框架

这是一个基于 Vue 3 + TypeScript + Ant Design Vue 构建的现代化管理系统框架。

## 技术栈

- **Vue 3** - 渐进式 JavaScript 框架
- **TypeScript** - 类型安全的 JavaScript 超集
- **Ant Design Vue** - 企业级 UI 组件库
- **Vue Router 4** - 官方路由管理器
- **Pinia** - Vue 状态管理库
- **Axios** - HTTP 客户端
- **Vite** - 下一代前端构建工具

## 项目结构

```
chprobe_web/
├── src/
│   ├── api/              # API 接口定义
│   │   └── user.ts       # 用户相关接口
│   ├── assets/           # 静态资源
│   ├── components/       # 公共组件
│   ├── layouts/          # 布局组件
│   │   └── MainLayout.vue # 主布局（侧边栏+顶部导航）
│   ├── router/           # 路由配置
│   │   └── index.ts      # 路由定义
│   ├── stores/           # 状态管理
│   │   ├── app.ts        # 应用状态
│   │   └── user.ts       # 用户状态
│   ├── types/            # TypeScript 类型定义
│   ├── utils/            # 工具函数
│   │   └── request.ts    # Axios 请求封装
│   ├── views/            # 页面组件
│   │   ├── dashboard/    # 仪表盘
│   │   ├── user/         # 用户管理
│   │   └── settings/     # 系统设置
│   ├── App.vue           # 根组件
│   └── main.ts           # 入口文件
├── .env                  # 环境变量
├── .env.development      # 开发环境变量
├── .env.production       # 生产环境变量
├── package.json          # 项目依赖
├── tsconfig.json         # TypeScript 配置
└── vite.config.ts        # Vite 配置
```

## 快速开始

### 安装依赖

```bash
npm install
```

### 启动开发服务器

```bash
npm run dev
```

访问 http://localhost:5173 查看应用

### 构建生产版本

```bash
npm run build
```

### 预览生产构建

```bash
npm run preview
```

## 核心功能

### 1. 布局系统

- 响应式侧边栏（支持折叠）
- 顶部导航栏（用户信息、下拉菜单）
- 内容区域（路由视图）

### 2. 页面示例

#### 仪表盘
- 统计卡片展示
- 最近活动列表
- 快速操作按钮

#### 用户管理
- 用户列表展示（表格）
- 搜索和筛选功能
- 新增/编辑/删除用户

#### 系统设置
- 基本设置（系统名称、语言、主题）
- 安全设置（密码修改、登录超时）
- 系统信息展示

### 3. 状态管理

使用 Pinia 管理应用状态：

- `useUserStore`: 用户信息管理
- `useAppStore`: 应用配置管理

### 4. API 请求

已封装 Axios 请求工具，支持：
- 请求/响应拦截器
- Token 自动注入
- 统一错误处理
- 请求超时配置

### 5. 路由系统

基于 Vue Router 4，支持：
- 嵌套路由
- 路由元信息
- 路由守卫（可扩展）

## 环境变量

创建 `.env` 文件配置环境变量：

```bash
VITE_API_BASE_URL=http://localhost:8080/api
```

## 开发规范

### 组件命名

- 页面组件使用 PascalCase，如 `UserManagement.vue`
- 公共组件使用 kebab-case，如 `search-input.vue`

### 代码风格

- 使用 TypeScript 类型定义
- 组件使用 Composition API
- 遵循 Vue 3 最佳实践

### 目录组织

- 按功能模块组织目录
- 页面组件放在 `views/`
- 可复用组件放在 `components/`
- API 接口放在 `api/`

## 扩展开发

### 添加新页面

1. 在 `views/` 下创建页面组件
2. 在 `router/index.ts` 中添加路由配置
3. 在 `layouts/MainLayout.vue` 中添加菜单项

### 添加新的 API 接口

1. 在 `api/` 下创建接口文件
2. 使用封装的 `request` 工具
3. 定义 TypeScript 类型

### 添加新的 Store

1. 在 `stores/` 下创建 store 文件
2. 使用 `defineStore` 定义
3. 在组件中通过 `useXxxStore()` 使用

## 常见问题

### Q: 如何修改 API 地址？

A: 修改 `.env` 文件中的 `VITE_API_BASE_URL`

### Q: 如何添加新的菜单项？

A: 在 `layouts/MainLayout.vue` 的 `<a-menu>` 中添加 `<a-menu-item>`

### Q: 如何自定义主题？

A: 参考 Ant Design Vue 主题定制文档，在 `main.ts` 中配置

## 浏览器支持

- Chrome >= 87
- Firefox >= 78
- Safari >= 14
- Edge >= 88

## License

MIT
