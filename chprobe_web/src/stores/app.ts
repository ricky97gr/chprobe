import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
  state: () => ({
    collapsed: false,
    loading: false
  }),
  actions: {
    toggleCollapsed() {
      this.collapsed = !this.collapsed
    },
    setCollapsed(collapsed: boolean) {
      this.collapsed = collapsed
    },
    setLoading(loading: boolean) {
      this.loading = loading
    }
  }
})
