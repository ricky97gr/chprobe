import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'
import './style.css'
import App from './App.vue'
import router from './router'
import * as Vue from 'vue'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(Antd)

app.mount('#app')

// 暴露完整的 Vue 和 VueRouter 全局变量供插件的 UMD 构建使用
// 插件构建时 external: ['vue', 'vue-router']，运行时从 window.Vue / window.VueRouter 获取
;(window as any).Vue = Vue
;(window as any).VueRouter = router

export { app }
