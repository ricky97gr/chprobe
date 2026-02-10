<template>
  <div>
    <a-card style="margin-bottom: 24px;">
      <a-table :columns="historyColumns" :data-source="upgradeHistory" row-key="id">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 'success' ? 'green' : 'red'">
              {{ record.status === 'success' ? '成功' : '失败' }}
            </a-tag>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-alert
      message="升级前请确保"
      description="1. 已备份重要数据；2. 系统有足够的磁盘空间；3. 升级过程中服务可能会短暂中断"
      type="warning"
      show-icon
      style="margin-bottom: 24px;"
    />

    <a-steps :current="0" status="process" style="margin-bottom: 24px;">
      <a-step title="检查版本" />
      <a-step title="下载升级包" />
      <a-step title="执行升级" />
      <a-step title="验证升级" />
    </a-steps>

    <a-card size="small">
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="升级方式">
          <a-radio-group v-model:value="upgradeMethod">
            <a-radio value="online">在线升级</a-radio>
            <a-radio value="offline">离线升级</a-radio>
          </a-radio-group>
        </a-descriptions-item>
        <a-descriptions-item label="升级命令" v-if="upgradeMethod === 'online'">
          <a-input
            :value="onlineUpgradeCmd"
            readonly
            style="margin-bottom: 8px;"
          />
          <a-button type="primary" @click="copyOnlineCmd" block>
            <template #icon>
              <CopyOutlined />
            </template>
            复制命令
          </a-button>
        </a-descriptions-item>
        <a-descriptions-item label="升级包下载" v-else>
          <a-button type="primary" @click="downloadUpgradePackage">
            <template #icon>
              <DownloadOutlined />
            </template>
            下载升级包
          </a-button>
          <p style="margin-top: 8px; color: #999; font-size: 12px;">
            下载完成后，上传到服务器并执行安装脚本
          </p>
        </a-descriptions-item>
        <a-descriptions-item label="升级日志">
          <a-button @click="viewUpgradeLogs">
            <template #icon>
              <FileTextOutlined />
            </template>
            查看升级日志
          </a-button>
        </a-descriptions-item>
      </a-descriptions>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { message } from 'ant-design-vue'
import { CopyOutlined, DownloadOutlined, FileTextOutlined } from '@ant-design/icons-vue'

// 服务端升级相关数据
const currentVersion = 'v1.0.0'
const currentVersionDate = '2024-01-01'
const latestVersion = 'v1.1.0'
const needUpgrade = true
const upgradeMethod = ref('online')

const onlineUpgradeCmd = computed(() => {
  return `curl -fsSL http://your-server-ip:8080/upgrade | bash -s -- --version ${latestVersion}`
})

const historyColumns = [
  { title: '版本', dataIndex: 'version', key: 'version' },
  { title: '升级时间', dataIndex: 'upgradeTime', key: 'upgradeTime' },
  { title: '状态', dataIndex: 'status', key: 'status' },
  { title: '操作人', dataIndex: 'operator', key: 'operator' },
  { title: '备注', dataIndex: 'remark', key: 'remark' }
]

const upgradeHistory = [
  {
    id: '1',
    version: 'v1.0.0',
    upgradeTime: '2024-01-01 10:00:00',
    status: 'success',
    operator: 'admin',
    remark: '初始安装'
  }
]

// 服务端升级方法
const copyOnlineCmd = () => {
  navigator.clipboard.writeText(onlineUpgradeCmd.value).then(() => {
    message.success('命令已复制到剪贴板')
  })
}

const downloadUpgradePackage = () => {
  message.info(`开始下载 ${latestVersion} 版本升级包...`)
}

const viewUpgradeLogs = () => {
  message.info('打开升级日志查看器')
}
</script>

<style scoped>
/* 可根据需要添加样式 */
</style>