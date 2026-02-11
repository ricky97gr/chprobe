<template>
  <div>
    <a-card class="container-card">



      <a-table
        :columns="columns"
        :data-source="containers"
        :pagination="pagination"
        :loading="loading"
        row-key="uuid"
        class="table-sm"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <a @click="handleView(record)">{{ record.name }}</a>
          </template>
          <template v-else-if="column.key === 'status'">
            <a-tag :color="getStatusColor(record.status)">
              {{ getStatusText(record.status) }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'cpuUsage'">
            <a-progress
              :percent="record.cpuUsage"
              :size="'small'"
              :stroke-color="getColor(record.cpuUsage)"
            />
          </template>
          <template v-else-if="column.key === 'memoryUsage'">
            <a-progress
              :percent="record.memoryUsage"
              :size="'small'"
              :stroke-color="getColor(record.memoryUsage)"
            />
          </template>
          <template v-else-if="column.key === 'action'">
            <a-popover
              placement="top"
              trigger="hover"
              :title="null"
            >
              <template #content>
                <div class="action-popover">
                  <div class="action-item" @click="handleView(record)">
                    <EyeOutlined style="margin-right: 4px;" />
                    查看
                  </div>
                  <div class="action-item" @click="handleLogs(record)">
                    <FileTextOutlined style="margin-right: 4px;" />
                    日志
                  </div>
                  <template v-if="record.status === 'running'">
                    <div class="action-item" @click="handleStop(record)">
                      <StopOutlined style="margin-right: 4px;" />
                      停止
                    </div>
                  </template>
                  <template v-else>
                    <div class="action-item" @click="handleStart(record)">
                      <PlayCircleOutlined style="margin-right: 4px;" />
                      启动
                    </div>
                  </template>
                  <div class="action-item action-danger" @click="handleDelete(record)">
                    <DeleteOutlined style="margin-right: 4px;" />
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
      v-model:open="detailVisible"
      title="容器详情"
      width="500px"
      :footer="null"
      :centered="true"
    >
      <div class="detail-content" style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px; padding: 16px;">
        <div v-for="(item, index) in detailData" :key="item.key" class="detail-item">
          <div class="detail-label">{{ item.key }}:</div>
          <div class="detail-value">
            <template v-if="item.status">
              <a-tag :color="getStatusColor(item.status)">
                {{ getStatusText(item.status) }}
              </a-tag>
            </template>
            <template v-else-if="item.cpuUsage !== undefined">
              <a-progress :percent="item.cpuUsage" :stroke-color="getColor(item.cpuUsage)" :size="'small'" />
            </template>
            <template v-else-if="item.memoryUsage !== undefined">
              <a-progress :percent="item.memoryUsage" :stroke-color="getColor(item.memoryUsage)" :size="'small'" />
            </template>
            <template v-else>
              {{ item.value }}
            </template>
          </div>
        </div>
      </div>
    </a-modal>

    <a-modal
      v-model:open="logsVisible"
      title="容器日志"
      width="800px"
      :footer="null"
      :centered="true"
    >
      <a-tabs v-model:activeKey="logsActiveKey">
        <a-tab-pane key="stdout" tab="标准输出">
          <pre class="log-pre">{{ logs.stdout }}</pre>
        </a-tab-pane>
        <a-tab-pane key="stderr" tab="标准错误">
          <pre class="log-pre log-pre-error">{{ logs.stderr }}</pre>
        </a-tab-pane>
      </a-tabs>
    </a-modal>
  </div>
</template>

<script setup lang="ts">import { ref, reactive, h, computed } from 'vue';
import { message } from 'ant-design-vue';
import { PlusOutlined, ReloadOutlined, DownOutlined, EyeOutlined, PlayCircleOutlined, StopOutlined, DeleteOutlined, FileTextOutlined } from '@ant-design/icons-vue';
import type { MenuProps } from 'ant-design-vue';
interface Container {
 id: string;
 name: string;
 image: string;
 status: 'running' | 'stopped' | 'exited';
 hostId: string;
 hostName: string;
 cpuUsage: number;
 memoryUsage: number;
 memoryLimit: string;
 network: string;
 ip: string;
 ports: string[];
 created: string;
 started: string;
}
const columns = [
 { title: '容器名', dataIndex: 'name', key: 'name' },
 { title: '镜像', dataIndex: 'image', key: 'image', width: 200 },
 { title: '状态', dataIndex: 'status', key: 'status', width: 100 },
 { title: '主机', dataIndex: 'hostName', key: 'hostName', width: 120 },
 { title: 'CPU', dataIndex: 'cpuUsage', key: 'cpuUsage', width: 120 },
 { title: '内存', dataIndex: 'memoryUsage', key: 'memoryUsage', width: 120 },
 { title: 'IP 地址', dataIndex: 'ip', key: 'ip', width: 150 },
 { title: '操作', key: 'action', width: 100, fixed: 'right' }
];
const detailColumns = [
 { title: '属性', dataIndex: 'key', key: 'key', width: 150 },
 { title: '值', dataIndex: 'value', key: 'value' }
];
const containers = ref<Container[]>([
 {
 id: '1',
 name: 'nginx-prod',
 image: 'nginx:1.25',
 status: 'running',
 hostId: '1',
 hostName: 'server-01',
 cpuUsage: 45,
 memoryUsage: 32,
 memoryLimit: '512MB',
 network: 'bridge',
 ip: '172.17.0.2',
 ports: ['80:80', '443:443'],
 created: '2024-01-25 10:00:00',
 started: '2024-01-25 10:00:00'
 },
 {
 id: '2',
 name: 'mysql-db',
 image: 'mysql:8.0',
 status: 'running',
 hostId: '2',
 hostName: 'server-02',
 cpuUsage: 28,
 memoryUsage: 68,
 memoryLimit: '2GB',
 network: 'host',
 ip: '192.168.1.102',
 ports: ['3306:3306'],
 created: '2024-01-24 14:30:00',
 started: '2024-01-24 14:30:00'
 },
 {
 id: '3',
 name: 'redis-cache',
 image: 'redis:7.0',
 status: 'stopped',
 hostId: '1',
 hostName: 'server-01',
 cpuUsage: 0,
 memoryUsage: 0,
 memoryLimit: '256MB',
 network: 'bridge',
 ip: '172.17.0.3',
 ports: ['6379:6379'],
 created: '2024-01-26 09:00:00',
 started: '2024-01-26 09:00:00'
 },
 {
 id: '4',
 name: 'node-app',
 image: 'node:18-alpine',
 status: 'running',
 hostId: '3',
 hostName: 'server-03',
 cpuUsage: 72,
 memoryUsage: 45,
 memoryLimit: '1GB',
 network: 'bridge',
 ip: '172.17.0.4',
 ports: ['3000:3000'],
 created: '2024-01-27 16:00:00',
 started: '2024-01-27 16:00:00'
 },
 {
 id: '5',
 name: 'postgres-db',
 image: 'postgres:15',
 status: 'exited',
 hostId: '2',
 hostName: 'server-02',
 cpuUsage: 0,
 memoryUsage: 0,
 memoryLimit: '1.5GB',
 network: 'bridge',
 ip: '172.17.0.5',
 ports: ['5432:5432'],
 created: '2024-01-23 11:00:00',
 started: '2024-01-23 11:00:00'
 },
 {
 id: '6',
 name: 'mongodb',
 image: 'mongo:6.0',
 status: 'running',
 hostId: '3',
 hostName: 'server-03',
 cpuUsage: 35,
 memoryUsage: 58,
 memoryLimit: '3GB',
 network: 'bridge',
 ip: '172.17.0.6',
 ports: ['27017:27017'],
 created: '2024-01-22 10:00:00',
 started: '2024-01-22 10:00:00'
 }
]);
const searchForm = reactive({
 name: '',
 status: '',
 hostId: ''
});
const pagination = reactive({
 current: 1,
 pageSize: 10,
 total: 6
});
const loading = ref(false);
const detailVisible = ref(false);
const logsVisible = ref(false);
const currentContainer = ref<Container | null>(null);
const logsActiveKey = ref('stdout');
const logs = reactive({
 stdout: '172.17.0.1 - - [28/Jan/2024:10:00:00 +0000] "GET / HTTP/1.1" 200 615 "-" "Mozilla/5.0"\n172.17.0.1 - - [28/Jan/2024:10:00:01 +0000] "GET /favicon.ico HTTP/1.1" 404 153 "-" "Mozilla/5.0"\n',
 stderr: ''
});
const detailData = computed(() => {
 if (!currentContainer.value)
 return [];
 const c = currentContainer.value;
 return [
 { key: '容器名', value: c.name },
 { key: '状态', value: getStatusText(c.status), status: c.status },
 { key: '镜像', value: c.image },
 { key: '主机', value: c.hostName },
 { key: 'CPU 使用率', value: `${c.cpuUsage}%`, cpuUsage: c.cpuUsage },
 { key: '内存使用率', value: `${c.memoryUsage}%`, memoryUsage: c.memoryUsage },
 { key: '内存限制', value: c.memoryLimit },
 { key: '网络', value: c.network },
 { key: 'IP 地址', value: c.ip },
 { key: '端口', value: c.ports.join(', ') || '-' },
 { key: '创建时间', value: c.created },
 { key: '启动时间', value: c.started }
 ];
});
const getStatusColor = (status: string) => {
 const colors: Record<string, string> = {
 running: 'success',
 stopped: 'default',
 exited: 'error'
 };
 return colors[status] || 'default';
};
const getStatusText = (status: string) => {
 const texts: Record<string, string> = {
 running: '运行中',
 stopped: '已停止',
 exited: '已退出'
 };
 return texts[status] || status;
};
const getColor = (percent: number) => {
 if (percent > 80)
 return '#f5222d';
 if (percent > 60)
 return '#faad14';
 return '#52c41a';
};
const getActionMenu = (record: Container) => {
 const items: MenuProps['items'] = [
 { key: 'view', label: '查看', icon: () => h(EyeOutlined) },
 { key: 'logs', label: '日志', icon: () => h(FileTextOutlined) }
 ];
 if (record.status === 'running') {
 items.push({ key: 'stop', label: '停止', icon: () => h(StopOutlined) });
 }
 else {
 items.push({ key: 'start', label: '启动', icon: () => h(PlayCircleOutlined) });
 }
 items.push({ type: 'divider' });
 items.push({ key: 'delete', label: '删除', icon: () => h(DeleteOutlined), danger: true });
 return {
 items,
 onClick: ({ key }: { key: string }) => handleAction(key, record)
 };
};
const handleAction = (key: string, record: Container) => {
 switch (key) {
 case 'view':
 handleView(record);
 break;
 case 'logs':
 handleLogs(record);
 break;
 case 'start':
 handleStart(record);
 break;
 case 'stop':
 handleStop(record);
 break;
 case 'delete':
 handleDelete(record);
 break;
 }
};
const handleSearch = () => {
 loading.value = true;
 setTimeout(() => {
 loading.value = false;
 message.success('搜索成功');
 }, 500);
};
const handleReset = () => {
 searchForm.name = '';
 searchForm.status = '';
 searchForm.hostId = '';
};
const handleRefresh = () => {
 message.success('刷新成功');
};
const handleCreate = () => {
 message.warning('创建容器功能待实现');
};
const handleView = (record: Container) => {
 currentContainer.value = record;
 detailVisible.value = true;
};
const handleLogs = (record: Container) => {
 logsVisible.value = true;
};
const handleStart = (record: Container) => {
 message.success(`启动容器 ${record.name}`);
};
const handleStop = (record: Container) => {
 message.success(`停止容器 ${record.name}`);
};
const handleDelete = (record: Container) => {
 message.warning(`删除容器 ${record.name} 功能待实现`);
};
</script>

<style scoped>
.container-card {
  height: 32px;
  overflow: visible;
}

.container-card :deep(.ant-card-head) {
  min-height: 32px !important;
  height: 32px !important;
  line-height: 32px !important;
  padding: 1px !important;
  margin: 0 !important;
  width: 100% !important;
  border: none !important;
  background: transparent !important;
  box-shadow: none !important;
}

.container-card :deep(.ant-card-head-title) {
  font-size: 12px !important;
  line-height: 32px !important;
  padding: 1px !important;
  margin: 0 !important;
  height: 32px !important;
}

.container-card :deep(.ant-card-body) {
  padding: 1px !important;
  margin: 0 !important;
}

.container-card :deep(.ant-card-extra) {
  height: 32px !important;
  line-height: 32px !important;
  padding: 0 !important;
  margin: 0 !important;
}

.btn-sm {
  height: 24px !important;
  padding: 0 8px !important;
  font-size: 12px !important;
  line-height: 24px !important;
}



.table-sm :deep(.ant-table-thead > tr > th) {
  height: 32px !important;
  line-height: 32px !important;
  padding: 1px !important;
  font-size: 14px !important;
  text-align: center !important;
}

.table-sm :deep(.ant-table-tbody > tr > td) {
  height: 32px !important;
  line-height: 32px !important;
  padding: 1px !important;
  font-size: 14px !important;
  text-align: center !important;
}

.btn-action {
  height: 24px !important;
  padding: 0 4px !important;
  font-size: 12px !important;
  line-height: 24px !important;
}

.detail-table :deep(.ant-table-thead > tr > th) {
  height: 32px !important;
  line-height: 32px !important;
  padding: 1px !important;
  font-size: 12px !important;
}

.detail-table :deep(.ant-table-tbody > tr > td) {
  height: 32px !important;
  line-height: 32px !important;
  padding: 1px !important;
  font-size: 12px !important;
}

.log-pre {
  background: #f0f0f0;
  padding: 8px;
  overflow: auto;
  max-height: 400px;
  font-size: 12px;
  line-height: 18px;
  margin: 0;
}

.log-pre-error {
  background: #fff5f5;
  color: #f5222d;
}

.action-popover {
  padding: 2px 0;
  width: 70px;
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

.detail-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0;
  border-bottom: 1px solid #f0f0f0;
}

.detail-item:last-child {
  border-bottom: none;
}

.detail-label {
  font-size: 12px;
  color: #666;
  font-weight: normal;
  min-width: 80px;
  flex-shrink: 0;
}

.detail-value {
  font-size: 12px;
  color: #333;
  min-height: 20px;
  display: flex;
  align-items: center;
  flex: 1;
}

.detail-content {
  max-height: 400px;
  overflow-y: auto;
}
</style>