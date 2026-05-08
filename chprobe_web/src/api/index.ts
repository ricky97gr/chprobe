// API接口定义和调用示例
import { get, post, postForm, put, del } from './request';

// 登录请求参数接口
export interface LoginRequest {
  username: string;
  password: string;
}

// 用户信息接口
export interface UserInfo {
  id: string;
  username: string;
  create_time: number;
  last_login_time: number;
  status: string;
  phone: string;
  email: string;
}

// 登录响应接口
export interface LoginResponse {
  token: string;
  user: UserInfo;
}

// 主机信息接口
export interface HostInfo {
  uuid: string;
  hostName: string;
  ip: string[];
  osType: string;
  os: string;
  arch: string;
  kernelVersion: string;
  cpu: string;
  memory: string;
  storage: string;
  registerTime: number;
  lastHeartTime: number;
}

// 客户端信息接口
export interface Agent {
  uuid: string;
  hostName: string;
  ip: string;
  machineID: string;
  clientType: string;
  osType: string;
  os: string;
  arch: string;
  kernelVersion: string;
  cpu: string;
  memory: string;
  storage: string;
  version: string;
  registerTime: number;
  lastHeartTime: number;
  status: string;
}

// 分页查询参数接口
export interface PageQuery {
  page: number;
  pageSize: number;
  startTime?: number;
  endTime?: number;
  conditions?: any[];
  sorts?: any[];
}

// 登录API
export const login = async (params: LoginRequest) => {
  return post<LoginResponse>('/login', params);
};

// 修改密码API
export const changePassword = async (params: { oldPassword: string; newPassword: string }) => {
  return post<any>('/user/change-password', params);
};

// 获取主机列表API
export const getHostList = async (params: PageQuery) => {
  return get<HostInfo[]>('/host/list', params);
};

// 获取主机详情API
export const getHostDetail = async (uuid: string) => {
  return get<HostInfo>(`/host/detail/${uuid}`);
};

// 获取镜像列表API
export const getImageList = async (params: PageQuery) => {
  return get<any[]>('/image/list', params);
};

// 获取镜像详情API
export const getImageDetail = async (id: string) => {
  return get<any>(`/image/detail/${id}`);
};

// 获取容器列表API
export const getContainerList = async (params: PageQuery) => {
  return get<any[]>('/container/list', params);
};

// 获取容器详情API
export const getContainerDetail = async (id: string) => {
  return get<any>(`/container/detail/${id}`);
};

// 获取运行日志列表API
export const getRunningLogList = async (params: PageQuery) => {
  return get<any[]>('/log/running/list', params);
};

// 获取操作日志列表API
export const getOperationLogList = async (params: PageQuery) => {
  return get<any[]>('/log/operation/list', params);
};

// 获取访问日志列表API
export const getAccessLogList = async (params: PageQuery) => {
  return get<any[]>('/log/access/list', params);
};

// 系统运行日志接口
export interface SystemLog {
  uuid: string;
  level: 'info' | 'warn' | 'error' | 'debug';
  module: string;
  message: string;
  serverIp: string;
  hostname: string;
  processName: string;
  pid: number;
  traceId: string;
  createdAt: number;
}

export interface SystemLogQuery extends PageQuery {
  level?: string;
  module?: string;
  keyword?: string;
}

// 获取系统运行日志列表API
export const getSystemLogList = async (params: SystemLogQuery) => {
  return get<SystemLog[]>('/log/system/list', params);
};

// 获取最新系统运行日志API（仪表盘）
export const getLatestSystemLog = async (limit: number = 10) => {
  return get<SystemLog[]>('/log/system/latest', { limit });
};

// 获取升级记录列表API
export interface UpgradeRecord {
  uuid: string;
  version: string;
  previousVersion: string;
  upgradeType: string;
  status: string;
  upgradeTime: number;
  serverIp: string;
  hostname: string;
  operator: string;
  description: string;
  errorMessage: string;
  createdAt: number;
  updatedAt: number;
}

export const getUpgradeRecordList = async (params: PageQuery) => {
  return get<UpgradeRecord[]>('/log/upgrade/list', params);
};

// 生成安装命令API
export const generateInstallCommand = async (params: { serverIp: string; osType: string }) => {
  return get<string>('/install', params);
};

// 下载安装包API
export const downloadInstaller = async (filename: string) => {
  window.open(`/api/download/${filename}`, '_blank');
};

// 上传文件API（示例）
export const uploadFile = async (file: File) => {
  const formData = new FormData();
  formData.append('file', file);
  return postForm<any>('/upload', formData);
};

// 用户管理API
export const getUserList = async () => {
  return get<any[]>('/user/list');
};

export const createUser = async (data: any) => {
  return post<any>('/user/create', data);
};

export const updateUser = async (id: string, data: any) => {
  return put<any>(`/user/update/${id}`, data);
};

export const deleteUser = async (id: string) => {
  return del<any>(`/user/delete/${id}`);
};

export const resetPassword = async (id: string) => {
  return post<any>(`/user/reset-password/${id}`);
};

// 授权信息接口
export interface AuthInfo {
  id: string;
  type: string;
  importTime: string;
  expireTime: string;
  status: string;
}

// 授权响应接口
export interface AuthResponse {
  productSerial: string;
  authInfo: AuthInfo[];
}

// 获取授权信息API
export const getAuthInfo = async () => {
  return get<AuthResponse>('/license/info');
};

// 上传授权信息API
export const uploadLicense = async (data: FormData | { license: string }) => {
  if (data instanceof FormData) {
    return postForm<any>('/license/upload', data);
  } else {
    return post<any>('/license/upload', data);
  }
};

// 获取授权详情API
export const getLicenseDetail = async (id: string) => {
  return get<any>(`/license/detail/${id}`);
};

// 删除授权API
export const deleteLicense = async (id: string) => {
  return del<any>(`/license/delete/${id}`);
};

// 获取系统信息API
export const getSystemInfo = async () => {
  return get<any>('/system/info');
};

// 获取服务器IP列表API
export const getServerIPs = async () => {
  return get<string[]>('/system/ips');
};

// 获取仪表盘统计数据API
export const getDashboardStats = async () => {
  return get<any>('/dashboard/stats');
};

// 获取客户端列表API
export const getAgentList = async () => {
  return get<Agent[]>('/agent/list');
};

// 获取客户端详情API
export const getAgentDetail = async (uuid: string) => {
  return get<Agent>(`/agent/detail/${uuid}`);
};

// 更新客户端状态API
export const updateAgentStatus = async (uuid: string, status: string) => {
  return put<any>(`/agent/status/${uuid}`, { status });
};

// 删除客户端API
export const deleteAgent = async (uuid: string) => {
  return del<any>(`/agent/delete/${uuid}`);
};

export default {
  login,
  changePassword,
  getHostList,
  getHostDetail,
  getImageList,
  getImageDetail,
  getContainerList,
  getContainerDetail,
  getRunningLogList,
  getOperationLogList,
  getAccessLogList,
  getSystemLogList,
  getLatestSystemLog,
  getUpgradeRecordList,
  generateInstallCommand,
  downloadInstaller,
  uploadFile,
  getUserList,
  createUser,
  updateUser,
  deleteUser,
  resetPassword,
  getAuthInfo,
  uploadLicense,
  getLicenseDetail,
  deleteLicense,
  getAgentList,
  getAgentDetail,
  updateAgentStatus,
  deleteAgent
};
