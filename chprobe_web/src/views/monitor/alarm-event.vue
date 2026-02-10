<template>
  <div>
    <a-card class="alarm-event-card">



      <a-table
        :columns="columns"
        :data-source="events"
        :pagination="pagination"
        :loading="loading"
        row-key="id"
        class="table-sm"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'level'">
            <a-tag :color="getLevelColor(record.level)">
              {{ getLevelText(record.level) }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'status'">
            <a-tag :color="getStatusColor(record.status)">
              {{ getStatusText(record.status) }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'value'">
            <a-progress
              :percent="record.value"
              :size="'small'"
              :stroke-color="getLevelColor(record.level)"
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
                  <div class="action-item" @click="handleDetail(record)">
                    详情
                  </div>
                  <div class="action-item" @click="handleAcknowledge(record)" :style="record.status !== 'firing' ? { cursor: 'not-allowed', color: '#999' } : {}">
                    处理
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
  </div>
</template>

<script setup lang="ts">import { ref, reactive, computed, h } from 'vue';
import { message } from 'ant-design-vue';
import { ReloadOutlined, DeleteOutlined, SearchOutlined, BellOutlined, ExclamationCircleOutlined, CheckCircleOutlined, InboxOutlined, StopOutlined, DownOutlined } from '@ant-design/icons-vue';
import type { MenuProps } from 'ant-design-vue';
interface Event {
 id: string;
 name: string;
 level: 'critical' | 'warning' | 'info';
 status: 'firing' | 'resolved' | 'acknowledged';
 host: string;
 metric: string;
 value: number;
 threshold: number;
 unit: string;
 started: string;
 ended: string;
 duration: string;
 message: string;
}
const columns = [
 { title: '告警名称', dataIndex: 'name', key: 'name' },
 { title: '告警级别', dataIndex: 'level', key: 'level', width: 100 },
 { title: '状态', dataIndex: 'status', key: 'status', width: 100 },
 { title: '目标主机', dataIndex: 'host', key: 'host', width: 150 },
 { title: '监控指标', dataIndex: 'metric', key: 'metric', width: 150 },
 { title: '当前值', dataIndex: 'value', key: 'value', width: 120 },
 { title: '触发时间', dataIndex: 'started', key: 'started', width: 180 },
 { title: '操作', key: 'action', width: 200, fixed: 'right' }
];
const events = ref<Event[]>([
 {
 id: '1',
 name: 'CPU 使用率过高',
 level: 'warning',
 status: 'firing',
 host: 'server-01',
 metric: 'cpu_usage',
 value: 85,
 threshold: 80,
 unit: '%',
 started: '2024-01-16 14:30:00',
 ended: '',
 duration: '15m 30s',
 message: 'CPU 使用率持续超过 80%'
 },
 {
 id: '2',
 name: '磁盘空间不足',
 level: 'critical',
 status: 'firing',
 host: 'server-02',
 metric: 'disk_usage',
 value: 92,
 threshold: 90,
 unit: '%',
 started: '2024-01-16 14:20:00',
 ended: '',
 duration: '25m 30s',
 message: '磁盘使用率达到 92%，请及时清理'
 },
 {
 id: '3',
 name: '内存使用率过高',
 level: 'warning',
 status: 'resolved',
 host: 'server-03',
 metric: 'memory_usage',
 value: 75,
 threshold: 85,
 unit: '%',
 started: '2024-01-16 13:00:00',
 ended: '2024-01-16 13:30:00',
 duration: '30m',
 message: '内存使用率已恢复正常'
 },
 {
 id: '4',
 name: '网络流量异常',
 level: 'info',
 status: 'acknowledged',
 host: 'server-01',
 metric: 'network_traffic',
 value: 120,
 threshold: 100,
 unit: 'MB/s',
 started: '2024-01-16 12:00:00',
 ended: '2024-01-16 12:30:00',
 duration: '30m',
 message: '网络流量突增'
 },
 {
 id: '5',
 name: '容器重启频繁',
 level: 'warning',
 status: 'firing',
 host: 'server-02',
 metric: 'container_restarts',
 value: 8,
 threshold: 5,
 unit: '次/小时',
 started: '2024-01-16 14:00:00',
 ended: '',
 duration: '45m',
 message: '容器在过去 1 小时内重启 8 次'
 },
 {
 id: '6',
 name: '进程不存在',
 level: 'critical',
 status: 'resolved',
 host: 'server-01',
 metric: 'process_exists',
 value: 0,
 threshold: 0,
 unit: '',
 started: '2024-01-16 10:00:00',
 ended: '2024-01-16 10:15:00',
 duration: '15m',
 message: '关键进程已恢复运行'
 }
]);
const searchForm = reactive({
 level: '',
 status: '',
 timeRange: []
});
const pagination = reactive({
 current: 1,
 pageSize: 10,
 total: 6
});
const loading = ref(false);
const detailVisible = ref(false);
const currentEvent = ref<Event | null>(null);
const todayCount = ref(28);
const firingCount = computed(() => events.value.filter(e => e.status === 'firing').length);
const resolvedCount = computed(() => events.value.filter(e => e.status === 'resolved').length);
const acknowledgedCount = computed(() => events.value.filter(e => e.status === 'acknowledged').length);
const criticalCount = computed(() => events.value.filter(e => e.level === 'critical' && e.status === 'firing').length);
const getLevelColor = (level: string) => {
 if (level === 'critical')
 return 'red';
 if (level === 'warning')
 return 'orange';
 return 'blue';
};
const getLevelText = (level: string) => {
 if (level === 'critical')
 return '紧急';
 if (level === 'warning')
 return '警告';
 return '提示';
};
const getStatusColor = (status: string) => {
 if (status === 'firing')
 return 'red';
 if (status === 'resolved')
 return 'green';
 return 'gray';
};
const getStatusText = (status: string) => {
 if (status === 'firing')
 return '告警中';
 if (status === 'resolved')
 return '已恢复';
 return '已处理';
};
const handleSearch = () => {
 loading.value = true;
 setTimeout(() => {
 loading.value = false;
 message.success('搜索成功');
 }, 500);
};
const handleReset = () => {
 searchForm.level = '';
 searchForm.status = '';
 searchForm.timeRange = [];
};
const handleRefresh = () => {
 loading.value = true;
 setTimeout(() => {
 loading.value = false;
 message.success('刷新成功');
 }, 500);
};
const getActionMenu = (record: Event) => {
 return {
 items: [
 {
 key: 'detail',
 label: '详情',
 onClick: () => handleDetail(record)
 },
 {
 key: 'acknowledge',
 label: '处理',
 onClick: () => handleAcknowledge(record),
 disabled: record.status !== 'firing'
 },
 {
 key: 'delete',
 label: '删除',
 danger: true,
 onClick: () => handleDelete(record)
 }
 ]
 }
};
const handleClear = () => {
 message.success('已清除已处理的告警');
};
const handleDetail = (record: Event) => {
 currentEvent.value = record;
 detailVisible.value = true;
};
const handleAcknowledge = (record: Event) => {
 record.status = 'acknowledged';
 message.success(`告警 ${record.name} 已标记为已处理`);
 if (detailVisible.value) {
 detailVisible.value = false;
 }
};
const handleDelete = (record: Event) => {
 message.success(`告警 ${record.name} 已删除`);
};
</script>

<style scoped>
.alarm-event-card {
  overflow: visible;
}

.alarm-event-card :deep(.ant-card-head) {
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

.alarm-event-card :deep(.ant-card-head-title) {
  font-size: 12px !important;
  line-height: 32px !important;
  padding: 1px !important;
  margin: 0 !important;
  height: 32px !important;
}

.alarm-event-card :deep(.ant-card-body) {
  padding: 1px !important;
  margin: 0 !important;
}

.alarm-event-card :deep(.ant-card-extra) {
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
  color: #1890ff !important;
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
</style>