<template>
  <div>
    <a-card>
      <div style="margin-bottom: 24px; display: flex; justify-content: space-between; align-items: center;">
        <a-descriptions :column="1" bordered style="flex: 1;">
          <a-descriptions-item label="产品序列号" :labelStyle="{ width: '180px' }">
            <div style="display: flex; align-items: center;">
              <span>{{ productSerial }}</span>
              <a-button 
                type="link" 
                size="small" 
                style="margin-left: 12px;" 
                @click="copySerial"
              >
                复制
              </a-button>
            </div>
          </a-descriptions-item>
        </a-descriptions>
        <a-button type="primary" @click="showUploadModal = true">
          上传授权
        </a-button>
      </div>

      <a-table 
        :columns="authColumns" 
        :data-source="authInfo" 
        :loading="loading"
        row-key="id" 
        style="width: 100%;"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 'valid' ? 'green' : 'red'">
              {{ record.status === 'valid' ? '有效' : '无效' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'action'">
            <a-space>
              <a-button type="link" size="small" @click="viewDetail(record)">
                查看详情
              </a-button>
              <a-popconfirm
                title="确定要删除此授权吗？"
                ok-text="确定"
                cancel-text="取消"
                @confirm="handleDeleteLicense(record.id)"
              >
                <a-button type="link" size="small" danger>
                  删除
                </a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal
      v-model:open="showUploadModal"
      title="上传授权"
      :confirm-loading="uploading"
      @ok="handleUpload"
      @cancel="handleCancel"
    >
      <a-tabs v-model:activeKey="uploadTab">
        <a-tab-pane key="string" tab="授权字符串">
          <a-form layout="vertical">
            <a-form-item label="授权字符串">
              <a-textarea
                v-model:value="licenseString"
                placeholder="请输入授权字符串"
                :rows="6"
              />
            </a-form-item>
          </a-form>
        </a-tab-pane>
        <a-tab-pane key="file" tab="授权文件">
          <a-upload
            :file-list="fileList"
            :before-upload="beforeUpload"
            @remove="handleRemove"
          >
            <a-button>
              <upload-outlined />
              选择文件
            </a-button>
          </a-upload>
        </a-tab-pane>
      </a-tabs>
    </a-modal>

    <a-modal
      v-model:open="showDetailModal"
      title="授权详情"
      :footer="null"
      width="600px"
    >
      <a-descriptions :column="1" bordered v-if="currentLicense">
        <a-descriptions-item label="授权ID">
          {{ currentLicense.id }}
        </a-descriptions-item>
        <a-descriptions-item label="授权序列号">
          {{ currentLicense.serial }}
        </a-descriptions-item>
        <a-descriptions-item label="授权类型">
          {{ currentLicense.type }}
        </a-descriptions-item>
        <a-descriptions-item label="导入时间">
          {{ currentLicense.importTime }}
        </a-descriptions-item>
        <a-descriptions-item label="过期时间">
          {{ currentLicense.expireTime }}
        </a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="currentLicense.status === 'valid' ? 'green' : 'red'">
            {{ currentLicense.status === 'valid' ? '有效' : '无效' }}
          </a-tag>
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { UploadOutlined } from '@ant-design/icons-vue'
import { getAuthInfo, uploadLicense, getLicenseDetail, deleteLicense } from '@/api'

// 加载状态
const loading = ref(false)

// 产品序列号
const productSerial = ref('')

// 授权信息
const authInfo = ref([])

// 上传模态框
const showUploadModal = ref(false)
const uploading = ref(false)
const uploadTab = ref('string')

// 授权字符串
const licenseString = ref('')

// 文件列表
const fileList = ref([])

// 详情模态框
const showDetailModal = ref(false)
const currentLicense = ref(null)

// 授权表格列
const authColumns = [
  { title: '授权类型', dataIndex: 'type', key: 'type' },
  { title: '导入时间', dataIndex: 'importTime', key: 'importTime' },
  { title: '过期时间', dataIndex: 'expireTime', key: 'expireTime' },
  { title: '状态', dataIndex: 'status', key: 'status' },
  { title: '操作', key: 'action', width: 150 }
]

// 复制产品序列号
const copySerial = () => {
  if (!productSerial.value) {
    message.warning('产品序列号为空')
    return
  }
  navigator.clipboard.writeText(productSerial.value)
    .then(() => {
      message.success('复制成功')
    })
    .catch(() => {
      message.error('复制失败')
    })
}

// 获取授权信息
const fetchAuthInfo = async () => {
  loading.value = true
  try {
    const response = await getAuthInfo()
    productSerial.value = response.result.productSerial
    authInfo.value = response.result.authInfo
  } catch (error) {
    message.error('获取授权信息失败')
    console.error('获取授权信息失败:', error)
    
    productSerial.value = ''
    authInfo.value = []
  } finally {
    loading.value = false
  }
}

// 上传前处理
const beforeUpload = (file: File) => {
  fileList.value = [file]
  return false
}

// 移除文件
const handleRemove = () => {
  fileList.value = []
}

// 处理上传
const handleUpload = async () => {
  uploading.value = true
  try {
    if (uploadTab.value === 'string') {
      if (!licenseString.value.trim()) {
        message.warning('请输入授权字符串')
        return
      }
      await uploadLicense({ license: licenseString.value })
    } else {
      if (fileList.value.length === 0) {
        message.warning('请选择授权文件')
        return
      }
      const formData = new FormData()
      formData.append('file', fileList.value[0])
      await uploadLicense(formData)
    }
    
    message.success('上传成功')
    showUploadModal.value = false
    licenseString.value = ''
    fileList.value = []
    await fetchAuthInfo()
  } catch (error) {
    message.error('上传失败')
    console.error('上传失败:', error)
  } finally {
    uploading.value = false
  }
}

// 取消上传
const handleCancel = () => {
  showUploadModal.value = false
  licenseString.value = ''
  fileList.value = []
}

// 查看详情
const viewDetail = async (record: any) => {
  try {
    const response = await getLicenseDetail(record.id)
    currentLicense.value = response.result
    showDetailModal.value = true
  } catch (error) {
    message.error('获取授权详情失败')
    console.error('获取授权详情失败:', error)
  }
}

// 删除授权
const handleDeleteLicense = async (id: string) => {
  try {
    await deleteLicense(id)
    message.success('删除成功')
    await fetchAuthInfo()
  } catch (error) {
    message.error('删除失败')
    console.error('删除失败:', error)
  }
}

// 组件挂载时获取授权信息
onMounted(() => {
  fetchAuthInfo()
})
</script>

<style scoped>
/* 可根据需要添加样式 */
pre {
  background: #f5f5f5;
  padding: 12px;
  border-radius: 4px;
  font-size: 12px;
}
</style>