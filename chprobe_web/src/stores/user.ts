import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  state: () => ({
    name: '管理员',
    avatar: '',
    roles: ['admin']
  }),
  actions: {
    setUserInfo(info: { name: string; avatar?: string; roles?: string[] }) {
      this.name = info.name
      if (info.avatar) this.avatar = info.avatar
      if (info.roles) this.roles = info.roles
    },
    logout() {
      this.name = ''
      this.avatar = ''
      this.roles = []
    }
  }
})
