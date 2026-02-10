<template>
  <div class="plugin-error-container">
    <a-result
      status="error"
      title="插件加载失败"
      sub-title="抱歉，插件组件加载失败，请检查插件配置或联系管理员"
    >
      <template #extra>
        <a-button type="primary" @click="goBack">
          返回上一页
        </a-button>
        <a-button @click="goHome">
          返回首页
        </a-button>
      </template>
    </a-result>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goBack = () => {
  router.back()
}

const goHome = () => {
  router.push('/dashboard')
}
</script>

<style scoped>
.plugin-error-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}
</style>
