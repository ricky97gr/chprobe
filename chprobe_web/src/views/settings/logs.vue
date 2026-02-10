<template>
  <div>
    <a-card :padding="1" :margin="1">
      <a-table :columns="logColumns" :data-source="logs" :pagination="pagination" row-key="id" style="width: 100%;">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'level'">
            <a-tag :color="getLevelColor(record.level)">
              {{ record.level.toUpperCase() }}
            </a-tag>
          </template>
          <template v-if="column.key === 'time'">
            {{ record.time }}
          </template>
        </template>
      </a-table>

      <div style="margin-top: 16px; text-align: right;">
        <a-button @click="exportLogs">
          <template #icon>
            <DownloadOutlined />
          </template>
          导出日志
        </a-button>
        <a-button type="danger" @click="clearLogs" style="margin-left: 8px;">
          <template #icon>
            <DeleteOutlined />
          </template>
          清空日志
        </a-button>
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { DownloadOutlined, DeleteOutlined } from '@ant-design/icons-vue'

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const logColumns = [
  { title: '时间', dataIndex: 'time', key: 'time', width: 200 },
  { title: '级别', dataIndex: 'level', key: 'level', width: 100 },
  { title: '模块', dataIndex: 'module', key: 'module', width: 120 },
  { title: '内容', dataIndex: 'content', key: 'content' }
]

const logs = [
  {
    id: '1',
    time: '2024-01-28 10:00:00',
    level: 'info',
    module: 'system',
    content: '系统启动成功',
    ip: '127.0.0.1'
  },
  {
    id: '2',
    time: '2024-01-28 10:01:30',
    level: 'warn',
    module: 'client',
    content: '客户端连接超时',
    ip: '192.168.1.101'
  },
  {
    id: '3',
    time: '2024-01-28 10:02:15',
    level: 'error',
    module: 'database',
    content: '数据库连接失败',
    ip: '127.0.0.1'
  }
]

const getLevelColor = (level: string) => {
  if (level === 'error') return 'red'
  if (level === 'warn') return 'orange'
  if (level === 'info') return 'blue'
  return 'default'
}

const exportLogs = () => {
  message.info('开始导出日志...')
}

const clearLogs = () => {
  message.warning('确定要清空所有日志吗？')
}
</script>

<style scoped>
/* 可根据需要添加样式 */
</style>