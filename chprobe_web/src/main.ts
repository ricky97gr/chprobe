import { createApp, defineComponent, h } from 'vue'
import { createPinia } from 'pinia'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'
import './style.css'
import App from './App.vue'
import router from './router'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(Antd)

app.mount('#app')

// 暴露全局变量供插件使用
;(window as any).Vue = {
  createApp,
  defineComponent,
  h
}
;(window as any).VueRouter = router

export { app }
