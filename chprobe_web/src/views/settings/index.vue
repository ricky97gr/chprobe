<template>
  <div>
    <a-card>
      <a-form :model="settings" :label-col="{ span: 6 }" :wrapper-col="{ span: 16 }">
        <a-form-item label="系统名称">
          <a-input v-model:value="settings.systemName" placeholder="请输入系统名称" />
        </a-form-item>
        <a-form-item label="系统描述">
          <a-textarea
            v-model:value="settings.description"
            placeholder="请输入系统描述"
            :rows="4"
          />
        </a-form-item>
        <a-form-item label="系统语言">
          <a-select v-model:value="settings.language" style="width: 200px;">
            <a-select-option value="zh-CN">简体中文</a-select-option>
            <a-select-option value="en-US">English</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="主题模式">
          <a-radio-group v-model:value="settings.theme">
            <a-radio value="light">浅色</a-radio>
            <a-radio value="dark">深色</a-radio>
            <a-radio value="auto">自动</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item :wrapper-col="{ offset: 6, span: 16 }">
          <a-button type="primary" @click="handleSave">
            <template #icon>
              <SaveOutlined />
            </template>
            保存设置
          </a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card style="margin-top: 16px;">
      <a-form :model="securitySettings" :label-col="{ span: 6 }" :wrapper-col="{ span: 16 }">
        <a-form-item label="登录密码">
          <a-input-password
            v-model:value="securitySettings.password"
            placeholder="请输入新密码"
          />
        </a-form-item>
        <a-form-item label="确认密码">
          <a-input-password
            v-model:value="securitySettings.confirmPassword"
            placeholder="请再次输入密码"
          />
        </a-form-item>
        <a-form-item label="登录超时">
          <a-input-number
            v-model:value="securitySettings.timeout"
            :min="5"
            :max="120"
            style="width: 200px;"
          />
          <span style="margin-left: 8px;">分钟</span>
        </a-form-item>
        <a-form-item label="登录验证">
          <a-switch v-model:checked="securitySettings.twoFactor" />
        </a-form-item>
        <a-form-item :wrapper-col="{ offset: 6, span: 16 }">
          <a-button type="primary" @click="handleSaveSecurity">
            <template #icon>
              <LockOutlined />
            </template>
            保存安全设置
          </a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card style="margin-top: 16px;">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="系统版本">
          v1.0.0
        </a-descriptions-item>
        <a-descriptions-item label="更新时间">
          2024-01-01
        </a-descriptions-item>
        <a-descriptions-item label="开发团队">
          Chprobe Team
        </a-descriptions-item>
        <a-descriptions-item label="联系方式">
          support@chprobe.com
        </a-descriptions-item>
      </a-descriptions>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { message } from 'ant-design-vue'
import { SaveOutlined, LockOutlined } from '@ant-design/icons-vue'

const settings = reactive({
  systemName: 'Chprobe 管理系统',
  description: '这是一个基于 Vue 3 + Ant Design Vue 的管理系统',
  language: 'zh-CN',
  theme: 'light'
})

const securitySettings = reactive({
  password: '',
  confirmPassword: '',
  timeout: 30,
  twoFactor: false
})

const handleSave = () => {
  message.success('基本设置已保存')
}

const handleSaveSecurity = () => {
  if (securitySettings.password !== securitySettings.confirmPassword) {
    message.error('两次输入的密码不一致')
    return
  }
  message.success('安全设置已保存')
}
</script>
