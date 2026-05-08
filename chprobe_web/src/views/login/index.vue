<template>
  <div class="login-container">
    <div class="login-form-wrapper">
      <h2 class="login-title">ChProbe</h2>
      <p class="login-subtitle">企业级容器与主机监控平台</p>
      
      <a-form
        :model="formState"
        name="login"
        @finish="handleSubmit"
        class="login-form"
      >
        <a-form-item
          name="username"
          :rules="[
            { required: true, message: '请输入用户名' },
            { min: 2, max: 20, message: '用户名长度应在2-20个字符之间' }
          ]"
          class="login-form-item"
        >
          <a-input
            v-model:value="formState.username"
            placeholder="用户名"
            prefix-icon="UserOutlined"
            size="large"
            class="login-input"
          />
        </a-form-item>

        <a-form-item
          name="password"
          :rules="[
            { required: true, message: '请输入密码' },
            { min: 6, message: '密码长度至少为6个字符' }
          ]"
          class="login-form-item"
        >
          <a-input-password
            v-model:value="formState.password"
            placeholder="密码"
            prefix-icon="LockOutlined"
            size="large"
            class="login-input"
          />
        </a-form-item>

        <a-form-item class="login-form-item">
          <a-checkbox v-model:checked="formState.remember">记住我</a-checkbox>
          <a href="#" class="login-forgot">忘记密码?</a>
        </a-form-item>

        <a-form-item class="login-form-item">
          <a-button
            type="primary"
            html-type="submit"
            class="login-button"
            size="large"
            block
            :loading="loading"
          >
            登录
          </a-button>
        </a-form-item>
      </a-form>
    </div>

    <!-- 强制修改密码模态框 -->
    <a-modal
      v-model:open="changePasswordModalVisible"
      title="强制修改密码"
      @ok="handleChangePassword"
      @cancel="handleCancelChangePassword"
      :centered="true"
      :ok-text="'确定'"
      :cancel-text="'取消'"
      :ok-button-props="{ disabled: !validatePasswordForm() }"
    >
      <a-form :model="passwordForm" class="password-form">
        <a-form-item
          name="oldPassword"
          label="原密码"
          :rules="[
            { required: true, message: '请输入原密码' }
          ]"
        >
          <a-input-password v-model:value="passwordForm.oldPassword" placeholder="请输入原密码" />
        </a-form-item>
        <a-form-item
          name="newPassword"
          label="新密码"
          :rules="[
            { required: true, message: '请输入新密码' },
            { min: 6, message: '密码长度至少为6个字符' }
          ]"
        >
          <a-input-password v-model:value="passwordForm.newPassword" placeholder="请输入新密码" />
        </a-form-item>
        <a-form-item
          name="confirmPassword"
          label="确认密码"
          :rules="confirmPasswordRules"
        >
          <a-input-password v-model:value="passwordForm.confirmPassword" placeholder="请再次输入新密码" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'
import { message, Modal } from 'ant-design-vue'
import { useRouter } from 'vue-router'
import { login, changePassword } from '@/api'

const router = useRouter()
const loading = ref(false)
const changePasswordModalVisible = ref(false)
const formState = reactive({
  username: '',
  password: '',
  remember: false
})
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const confirmPasswordRules = [
  { required: true, message: '请再次输入新密码' },
  {
    validator: (_rule: any, value: string) => {
      if (!value || passwordForm.newPassword === value) {
        return Promise.resolve()
      }
      return Promise.reject(new Error('两次输入的密码不一致'))
    }
  }
]

const validatePasswordForm = () => {
  return passwordForm.oldPassword && 
         passwordForm.newPassword && 
         passwordForm.confirmPassword && 
         passwordForm.newPassword === passwordForm.confirmPassword &&
         passwordForm.newPassword !== passwordForm.oldPassword
}

const handleSubmit = async () => {
  loading.value = true
  
  try {
    // 发送真实登录请求
    const { result } = await login({
      username: formState.username,
      password: formState.password
    })
    
    // 登录成功
    message.success('登录成功，正在跳转...', 3)
    
    // 设置token到localStorage
    localStorage.setItem('token', result.token)
    
    // 如果需要记住登录状态，可以存储更多信息
    if (formState.remember) {
      localStorage.setItem('username', formState.username)
    }
    
    // 检查用户是否首次登录
    if ((result.user as any).isFirstLogin) {
      // 显示强制修改密码模态框
      changePasswordModalVisible.value = true
    } else {
      // 跳转到首页
      setTimeout(() => {
        router.push('/')
      }, 500)
    }
  } catch (error: any) {
    // 登录失败
    message.error(error.message || '登录失败，请稍后重试', 3)
  } finally {
    loading.value = false
  }
}

const handleChangePassword = async () => {
  if (!validatePasswordForm()) {
    return
  }
  
  loading.value = true
  
  try {
    // 发送修改密码请求
    await changePassword({
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword
    })
    
    // 修改密码成功
    message.success('密码修改成功，正在跳转...', 3)
    
    // 隐藏模态框
    changePasswordModalVisible.value = false
    
    // 跳转到首页
    setTimeout(() => {
      router.push('/')
    }, 500)
  } catch (error: any) {
    // 修改密码失败
    message.error(error.message || '密码修改失败，请稍后重试', 3)
  } finally {
    loading.value = false
  }
}

const handleCancelChangePassword = () => {
  // 取消修改密码，退出登录
  Modal.confirm({
    title: '确认退出',
    content: '您需要修改密码才能继续使用系统，确定要退出登录吗？',
    okText: '确定',
    cancelText: '取消',
    onOk: () => {
      // 清除token
      localStorage.removeItem('token')
      // 刷新页面
      window.location.reload()
    }
  })
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f0f2f5;
  padding: 0 20px;
}

.login-form-wrapper {
  background: #fff;
  border-radius: 6px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  padding: 20px;
  width: 100%;
  max-width: 320px;
}

.login-header {
  text-align: center;
  margin-bottom: 20px;
}

.login-title {
  font-size: 18px;
  font-weight: 600;
  color: #1890ff;
  margin-bottom: 4px;
}

.login-subtitle {
  font-size: 11px;
  color: rgba(0, 0, 0, 0.45);
  margin: 0;
}

.login-form {
  width: 100%;
}

.login-form-item {
  margin-bottom: 14px;
}

.login-input {
  height: 36px;
  font-size: 13px;
}

.login-button {
  height: 36px;
  font-size: 13px;
  font-weight: 500;
  margin-top: 2px;
}

.login-forgot {
  float: right;
  font-size: 10px;
  color: #1890ff;
  text-decoration: none;
}

.login-forgot:hover {
  text-decoration: underline;
}

/* 调整提示信息字体大小 */
:deep(.ant-form-item-explain) {
  font-size: 10px !important;
}

:deep(.ant-message-notice-content) {
  font-size: 10px !important;
}
</style>