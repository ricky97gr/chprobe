# 前端打包配置修复总结

## 问题分析

### 开发模式 vs 生产模式

| 模式 | 前端地址 | 后端地址 | API配置 |
|------|---------|---------|---------|
| 开发 | http://localhost:5173 | http://10.186.39.112:8080 | 完整URL |
| 生产 | http://服务器IP:8080 | http://服务器IP:8081 | 相对路径 |

### 根本原因

1. **生产环境配置错误**：`.env.production` 配置了错误的完整URL
2. **API路径缺少前缀**：所有API调用都缺少 `/api` 前缀
3. **Nginx转发规则**：只转发 `/api/` 路径到后端

## 已修复内容

### 1. 生产环境配置 (.env.production)

```bash
# 修改前
VITE_API_BASE_URL=https://api.chprobe.com/api

# 修改后
VITE_API_BASE_URL=/api
```

### 2. API路径修复 (src/api/index.ts)

所有API调用都添加了 `/api` 前缀：

```typescript
// 修改前
export const login = async (params: LoginRequest) => {
  return post<LoginResponse>('/login', params);
};

// 修改后
export const login = async (params: LoginRequest) => {
  return post<LoginResponse>('/api/login', params);
};
```

共修复33个API接口：
- 登录相关：login, changePassword
- 主机管理：getHostList, getHostDetail
- 镜像管理：getImageList, getImageDetail
- 容器管理：getContainerList, getContainerDetail
- 日志管理：getRunningLogList, getOperationLogList, getAccessLogList, getSystemLogList, getLatestSystemLog, getUpgradeRecordList
- 用户管理：getUserList, createUser, updateUser, deleteUser, resetPassword
- 授权管理：getAuthInfo, uploadLicense, getLicenseDetail, deleteLicense
- 系统管理：getSystemInfo, getServerIPs, getDashboardStats
- 客户端管理：getAgentList, getAgentDetail, updateAgentStatus, deleteAgent
- 其他：generateInstallCommand, downloadInstaller, uploadFile

## 网络架构

```
用户访问 http://服务器IP:8080
    │
    ▼
┌─────────────────────────────────────────┐
│  Nginx (8080)                          │
│  ├─ /          → 前端静态文件          │
│  └─ /api/      → 后端 (127.0.0.1:8081) │
└─────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│  后端服务 (8081)                       │
│  所有API路由都在 /api 下                │
└─────────────────────────────────────────┘
```

## 请求流程

### 开发模式
```
前端 (localhost:5173)
  ↓
http://10.186.39.112:8080/api/login
  ↓
后端 (10.186.39.112:8080)
```

### 生产模式
```
前端 (8080)
  ↓
/api/login
  ↓
Nginx (8080) 转发
  ↓
http://127.0.0.1:8081/api/login
  ↓
后端 (8081)
```

## 环境变量说明

| 变量 | 开发环境 | 生产环境 | 说明 |
|------|---------|---------|------|
| VITE_API_BASE_URL | http://10.186.39.112:8080/api | /api | API基础URL |

## 构建命令

```bash
# 开发模式
cd chprobe_web
npm run dev

# 生产构建
npm run build

# Docker构建（包含前端和后端）
make docker-build
make docker-run
```

## 验证步骤

1. 前端页面访问：`http://服务器IP:8080`
2. 登录功能：使用正确的用户名密码登录
3. API请求：检查浏览器开发者工具 Network 面板
4. 后端日志：查看是否有API请求日志

## 注意事项

1. **端口映射**：使用 `--network host` 模式，容器直接使用宿主机网络
2. **MySQL连接**：容器内使用 `127.0.0.1` 连接宿主机MySQL
3. **Nginx配置**：确保 `/api/` 路径正确转发到后端8081端口
4. **环境变量**：生产构建时使用 `.env.production` 配置
