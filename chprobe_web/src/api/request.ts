// API请求方法封装
import axios from 'axios';
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';

// 后端响应格式接口
export interface ApiResponse<T = any> {
  code: number;
  msg: string;
  result: T;
  total: number;
}

// 错误响应格式接口
export interface ApiErrorResponse {
  code: number;
  msg: string;
}

// 创建axios实例
const service: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
});

// 请求拦截器
service.interceptors.request.use(
  (config) => {
    // 从localStorage获取token
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    console.error('请求错误:', error);
    return Promise.reject(error);
  }
);

// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const res = response.data;
    
    // 检查响应码
    if (res.code !== 200) {
      // 处理错误响应
      console.error('API错误:', res.msg);
      return Promise.reject(new Error(res.msg || 'Error'));
    }
    
    return response;
  },
  (error) => {
    console.error('网络错误:', error);
    
    // 处理网络错误
    let errorMessage = '网络错误，请稍后重试';
    if (error.response) {
      // 服务器返回错误
      const data = error.response.data as ApiErrorResponse;
      errorMessage = data.msg || errorMessage;
    } else if (error.request) {
      // 请求已发送但没有收到响应 - 可能是后端服务不通
      errorMessage = '服务器无响应，请稍后重试';
      
      // 触发后端服务降级事件
      window.dispatchEvent(new CustomEvent('backend-down'));
    }
    
    return Promise.reject(new Error(errorMessage));
  }
);

// 通用请求方法
export const request = async <T = any>(
  config: AxiosRequestConfig
): Promise<{ result: T; total: number }> => {
  try {
    const response = await service.request<ApiResponse<T>>(config);
    return {
      result: response.data.result,
      total: response.data.total
    };
  } catch (error) {
    throw error;
  }
};

// GET请求方法
export const get = async <T = any>(
  url: string,
  params?: any
): Promise<{ result: T; total: number }> => {
  return request<T>({
    method: 'GET',
    url,
    params
  });
};

// POST请求方法
export const post = async <T = any>(
  url: string,
  data?: any
): Promise<{ result: T; total: number }> => {
  return request<T>({
    method: 'POST',
    url,
    data
  });
};

// 发送表单数据的POST请求
export const postForm = async <T = any>(
  url: string,
  formData: FormData
): Promise<{ result: T; total: number }> => {
  return request<T>({
    method: 'POST',
    url,
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
};

// PUT请求方法
export const put = async <T = any>(
  url: string,
  data?: any
): Promise<{ result: T; total: number }> => {
  return request<T>({
    method: 'PUT',
    url,
    data
  });
};

// DELETE请求方法
export const del = async <T = any>(
  url: string,
  params?: any
): Promise<{ result: T; total: number }> => {
  return request<T>({
    method: 'DELETE',
    url,
    params
  });
};

export default {
  request,
  get,
  post,
  postForm,
  put,
  delete: del
};
