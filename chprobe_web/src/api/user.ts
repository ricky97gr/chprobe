import request from '@/utils/request'

export interface User {
  id: string
  username: string
  email: string
  status: string
  createTime: string
}

export const getUserList = (params: {
  page: number
  pageSize: number
  username?: string
  status?: string
}) => {
  return request.get<User[]>('/user/list', { params })
}

export const getUserById = (id: string) => {
  return request.get<User>(`/user/${id}`)
}

export const createUser = (data: Partial<User>) => {
  return request.post('/user/create', data)
}

export const updateUser = (id: string, data: Partial<User>) => {
  return request.put(`/user/${id}`, data)
}

export const deleteUser = (id: string) => {
  return request.delete(`/user/${id}`)
}
