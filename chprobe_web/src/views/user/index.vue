<template>
  <div>
    <a-card>
      <!-- 头部操作区 -->
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
        <!-- 查询表单 -->
        <div class="search-form" style="display: flex; gap: 8px;">
          <a-input
            v-model:value="searchForm.username"
            placeholder="请输入用户名"
            style="width: 180px;"
            prefix-icon="UserOutlined"
          />
          <a-input
            v-model:value="searchForm.email"
            placeholder="请输入邮箱"
            style="width: 200px;"
            prefix-icon="MailOutlined"
          />
          <a-button type="primary" @click="handleSearch">
            <template #icon>
              <SearchOutlined />
            </template>
            查询
          </a-button>
          <a-button @click="handleReset">重置</a-button>
        </div>

        <!-- 新增用户按钮 -->
        <a-button type="primary" @click="handleAdd">
          <template #icon>
            <PlusOutlined />
          </template>
          新增用户
        </a-button>
      </div>

      <a-table
        :columns="columns"
        :data-source="filteredUsers"
        :pagination="pagination"
        :loading="loading"
        row-key="id"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 'active' ? 'green' : 'red'">
              {{ record.status === 'active' ? '启用' : '禁用' }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-popover
              placement="top"
              trigger="hover"
              :title="null"
            >
              <template #content>
                <div class="action-popover">
                  <div class="action-item" @click="handleEdit(record)">
                    编辑
                  </div>
                  <div class="action-item" @click="handleResetPassword(record)">
                    重置密码
                  </div>
                  <div class="action-item action-danger" @click="handleDelete(record)">
                    删除
                  </div>
                </div>
              </template>
              <a-button type="text" size="small" class="btn-action">
                操作
              </a-button>
            </a-popover>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal
      v-model:open="modalVisible"
      :title="modalTitle"
      @ok="handleModalOk"
      @cancel="handleModalCancel"
      :centered="true"
      :ok-text="'确定'"
      :cancel-text="'取消'"
      :ok-button-props="{ style: { position: 'absolute', bottom: '20px', left: '20px' } }"
      :cancel-button-props="{ style: { position: 'absolute', bottom: '20px', right: '20px' } }"
    >
      <a-form :model="userForm" :label-col="{ span: 6 }" :wrapper-col="{ span: 16 }">
        <a-form-item label="用户名">
          <a-input v-model:value="userForm.username" placeholder="请输入用户名" />
        </a-form-item>
        <a-form-item label="邮箱">
          <a-input v-model:value="userForm.email" placeholder="请输入邮箱" />
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model:checked="isActive" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { PlusOutlined, SearchOutlined, EditOutlined, DeleteOutlined, UserOutlined, MailOutlined, KeyOutlined } from '@ant-design/icons-vue'
import { getUserList, createUser, updateUser, deleteUser, resetPassword } from '@/api'

interface User {
  id: string
  username: string
  email: string
  status: string
  createTime: string
  lastLoginTime: string
  isFirstLogin: boolean
}

const loading = ref(false)
const modalVisible = ref(false)
const modalTitle = ref('新增用户')

const userForm = reactive({
  id: '',
  username: '',
  email: '',
  status: 'active'
})

const searchForm = reactive({
  username: '',
  email: ''
})

const isActive = computed({
  get: () => userForm.status === 'active',
  set: (value: boolean) => {
    userForm.status = value ? 'active' : 'inactive'
  }
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const users = ref<User[]>([])

// 过滤后的用户列表
const filteredUsers = computed(() => {
  return users.value.filter(user => {
    const matchesUsername = !searchForm.username || user.username.includes(searchForm.username)
    const matchesEmail = !searchForm.email || user.email.includes(searchForm.email)
    return matchesUsername && matchesEmail
  })
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id' },
  { title: '用户名', dataIndex: 'username', key: 'username' },
  { title: '邮箱', dataIndex: 'email', key: 'email' },
  { title: '状态', dataIndex: 'status', key: 'status' },
  { title: '创建时间', dataIndex: 'createTime', key: 'createTime' },
  { 
    title: '上次登录时间', 
    dataIndex: 'lastLoginTime', 
    key: 'lastLoginTime',
    render: (text: string) => text || '-' 
  },
  { title: '操作', key: 'action', fixed: 'right', width: 120 }
]

// 获取用户列表
const fetchUserList = async () => {
  loading.value = true
  try {
    const { result } = await getUserList()
    users.value = result
    pagination.total = result.length
  } catch (error) {
    message.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  modalTitle.value = '新增用户'
  userForm.id = ''
  userForm.username = ''
  userForm.email = ''
  userForm.status = 'active'
  modalVisible.value = true
}

const handleEdit = (record: User) => {
  modalTitle.value = '编辑用户'
  userForm.id = record.id
  userForm.username = record.username
  userForm.email = record.email
  userForm.status = record.status
  modalVisible.value = true
}

const handleDelete = async (record: User) => {
  try {
    await deleteUser(record.id)
    message.success(`删除用户: ${record.username} 成功`)
    // 重新获取用户列表
    await fetchUserList()
  } catch (error) {
    message.error('删除用户失败')
  }
}

const handleResetPassword = async (record: User) => {
  Modal.confirm({
    title: '重置密码',
    content: `确定要将用户 ${record.username} 的密码重置为默认密码 123456 吗？\n\n重置密码后，用户需要使用新密码重新登录系统。`,
    okText: '确定',
    cancelText: '取消',
    onOk: async () => {
      try {
        await resetPassword(record.id)
        message.success(`重置用户 ${record.username} 的密码成功，新密码为 123456`)
        // 重新获取用户列表
        await fetchUserList()
      } catch (error) {
        message.error('重置密码失败')
      }
    }
  })
}

const handleModalOk = async () => {
  loading.value = true
  try {
    if (userForm.id) {
      // 更新用户
      await updateUser(userForm.id, userForm)
      message.success('编辑用户成功')
    } else {
      // 新增用户
      await createUser(userForm)
      message.success('新增用户成功')
    }
    modalVisible.value = false
    // 重新获取用户列表
    await fetchUserList()
  } catch (error) {
    message.error('保存用户失败')
  } finally {
    loading.value = false
  }
}

const handleModalCancel = () => {
  modalVisible.value = false
}

const handleSearch = () => {
  // 这里可以添加实际的搜索逻辑
  message.success('查询成功')
}

const handleReset = () => {
  searchForm.username = ''
  searchForm.email = ''
}

// 组件挂载时获取用户列表
onMounted(() => {
  fetchUserList()
})
</script>

<style scoped>
.action-popover {
  padding: 2px 0;
  width: 80px;
  white-space: nowrap;
}

.action-item {
  padding: 3px 4px;
  font-size: 14px;
  color: #1890ff;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  white-space: nowrap;
}

.action-item:hover {
  background-color: #e6f7ff;
  color: #1890ff;
}

.action-danger {
  color: #f5222d;
}

.action-danger:hover {
  background-color: #fff1f0;
  color: #f5222d;
}
</style>
