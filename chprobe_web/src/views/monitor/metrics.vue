<template>
  <div>
    <a-card>
      <template #title>
        监控指标
      </template>
      <template #extra>
        <a-space>
          <a-select v-model:value="selectedHost" style="width: 150px;">
            <a-select-option value="all">全部主机</a-select-option>
            <a-select-option v-for="host in hosts" :key="host.id" :value="host.id">
              {{ host.name }}
            </a-select-option>
          </a-select>
          <a-select v-model:value="timeRange" style="width: 120px;">
            <a-select-option value="1h">最近1小时</a-select-option>
            <a-select-option value="6h">最近6小时</a-select-option>
            <a-select-option value="24h">最近24小时</a-select-option>
            <a-select-option value="7d">最近7天</a-select-option>
          </a-select>
          <a-button @click="handleRefresh">
            <template #icon>
              <ReloadOutlined />
            </template>
            刷新
          </a-button>
        </a-space>
      </template>

      <a-tabs v-model:activeKey="activeTab" type="card">
        <a-tab-pane key="overview" tab="概览">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :sm="12" :lg="6">
              <a-card>
                <a-statistic title="CPU 使用率" :value="currentMetrics.cpu" suffix="%">
                  <template #prefix>
                    <CpuOutlined />
                  </template>
                  <template #valueStyle>
                    { color: getMetricColor(currentMetrics.cpu) }
                  </template>
                </a-statistic>
                <a-progress
                  :percent="currentMetrics.cpu"
                  :stroke-color="getMetricColor(currentMetrics.cpu)"
                  style="margin-top: 16px;"
                />
              </a-card>
            </a-col>
            <a-col :xs="24" :sm="12" :lg="6">
              <a-card>
                <a-statistic title="内存使用率" :value="currentMetrics.memory" suffix="%">
                  <template #prefix>
                    <DatabaseOutlined />
                  </template>
                  <template #valueStyle>
                    { color: getMetricColor(currentMetrics.memory) }
                  </template>
                </a-statistic>
                <a-progress
                  :percent="currentMetrics.memory"
                  :stroke-color="getMetricColor(currentMetrics.memory)"
                  style="margin-top: 16px;"
                />
              </a-card>
            </a-col>
            <a-col :xs="24" :sm="12" :lg="6">
              <a-card>
                <a-statistic title="磁盘使用率" :value="currentMetrics.disk" suffix="%">
                  <template #prefix>
                    <HddOutlined />
                  </template>
                  <template #valueStyle>
                    { color: getMetricColor(currentMetrics.disk) }
                  </template>
                </a-statistic>
                <a-progress
                  :percent="currentMetrics.disk"
                  :stroke-color="getMetricColor(currentMetrics.disk)"
                  style="margin-top: 16px;"
                />
              </a-card>
            </a-col>
            <a-col :xs="24" :sm="12" :lg="6">
              <a-card>
                <a-statistic title="网络流量" :value="currentMetrics.network" suffix="MB/s">
                  <template #prefix>
                    <CloudOutlined />
                  </template>
                  <template #valueStyle>
                    { color: getMetricColor(currentMetrics.network, 100) }
                  </template>
                </a-statistic>
                <a-progress
                  :percent="Math.min(currentMetrics.network, 100)"
                  :stroke-color="getMetricColor(currentMetrics.network, 100)"
                  style="margin-top: 16px;"
                />
              </a-card>
            </a-col>
          </a-row>

          <a-row :gutter="[16, 16]" style="margin-top: 16px;">
            <a-col :xs="24" :lg="12">
              <a-card title="CPU 趋势">
                <div style="height: 300px;">
                  <a-empty description="图表占位 - 集成 ECharts 后显示 CPU 使用率趋势图" />
                </div>
              </a-card>
            </a-col>
            <a-col :xs="24" :lg="12">
              <a-card title="内存趋势">
                <div style="height: 300px;">
                  <a-empty description="图表占位 - 集成 ECharts 后显示内存使用率趋势图" />
                </div>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>

        <a-tab-pane key="cpu" tab="CPU">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :lg="12">
              <a-card title="CPU 使用率">
                <div style="height: 300px;">
                  <a-empty description="图表占位" />
                </div>
              </a-card>
            </a-col>
            <a-col :xs="24" :lg="12">
              <a-card title="负载平均值">
                <div style="height: 300px;">
                  <a-empty description="图表占位" />
                </div>
              </a-card>
            </a-col>
          </a-row>
          <a-row :gutter="[16, 16]" style="margin-top: 16px;">
            <a-col :xs="24">
              <a-card title="进程 TOP 10">
                <a-table :columns="processColumns" :data-source="processes" :pagination="false" row-key="pid">
                  <template #bodyCell="{ column, record }">
                    <template v-if="column.key === 'cpu'">
                      <a-progress :percent="record.cpu" size="small" />
                    </template>
                    <template v-if="column.key === 'memory'">
                      <a-progress :percent="record.memory" size="small" />
                    </template>
                  </template>
                </a-table>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>

        <a-tab-pane key="memory" tab="内存">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :lg="12">
              <a-card title="内存使用">
                <div style="height: 300px;">
                  <a-empty description="图表占位" />
                </div>
              </a-card>
            </a-col>
            <a-col :xs="24" :lg="12">
              <a-card title="交换分区">
                <div style="height: 300px;">
                  <a-empty description="图表占位" />
                </div>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>

        <a-tab-pane key="disk" tab="磁盘">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :lg="12">
              <a-card title="磁盘使用">
                <a-table :columns="diskColumns" :data-source="disks" :pagination="false" row-key="mount">
                  <template #bodyCell="{ column, record }">
                    <template v-if="column.key === 'usage'">
                      <a-progress :percent="record.usage" size="small" />
                    </template>
                  </template>
                </a-table>
              </a-card>
            </a-col>
            <a-col :xs="24" :lg="12">
              <a-card title="I/O 性能">
                <div style="height: 300px;">
                  <a-empty description="图表占位" />
                </div>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>

        <a-tab-pane key="network" tab="网络">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :lg="12">
              <a-card title="入站流量">
                <div style="height: 300px;">
                  <a-empty description="图表占位" />
                </div>
              </a-card>
            </a-col>
            <a-col :xs="24" :lg="12">
              <a-card title="出站流量">
                <div style="height: 300px;">
                  <a-empty description="图表占位" />
                </div>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">import { ref, reactive } from 'vue';
import { message } from 'ant-design-vue';
import { ReloadOutlined, DatabaseOutlined, HddOutlined, CloudOutlined } from '@ant-design/icons-vue';
interface Host {
 id: string;
 name: string;
}
interface Process {
 pid: number;
 name: string;
 cpu: number;
 memory: number;
 user: string;
}
interface Disk {
 mount: string;
 total: string;
 used: string;
 free: string;
 usage: number;
}
const hosts = ref<Host[]>([
 { id: '1', name: 'server-01' },
 { id: '2', name: 'server-02' },
 { id: '3', name: 'server-03' }
]);
const selectedHost = ref('all');
const timeRange = ref('1h');
const activeTab = ref('overview');
const currentMetrics = reactive({
 cpu: 45,
 memory: 62,
 disk: 78,
 network: 45
});
const processColumns = [
 { title: 'PID', dataIndex: 'pid', key: 'pid', width: 80 },
 { title: '进程名', dataIndex: 'name', key: 'name' },
 { title: 'CPU', dataIndex: 'cpu', key: 'cpu', width: 150 },
 { title: '内存', dataIndex: 'memory', key: 'memory', width: 150 },
 { title: '用户', dataIndex: 'user', key: 'user', width: 100 }
];
const processes = ref<Process[]>([
 { pid: 1234, name: 'nginx', cpu: 25, memory: 15, user: 'root' },
 { pid: 5678, name: 'mysql', cpu: 18, memory: 35, user: 'mysql' },
 { pid: 9012, name: 'redis', cpu: 5, memory: 8, user: 'redis' },
 { pid: 3456, name: 'node', cpu: 12, memory: 22, user: 'www-data' },
 { pid: 7890, name: 'docker', cpu: 8, memory: 12, user: 'root' },
 { pid: 2345, name: 'sshd', cpu: 3, memory: 5, user: 'root' },
 { pid: 6789, name: 'systemd', cpu: 2, memory: 3, user: 'root' },
 { pid: 1111, name: 'postgres', cpu: 15, memory: 28, user: 'postgres' },
 { pid: 2222, name: 'rabbitmq', cpu: 10, memory: 20, user: 'rabbitmq' },
 { pid: 3333, name: 'java', cpu: 20, memory: 35, user: 'tomcat' }
]);
const diskColumns = [
 { title: '挂载点', dataIndex: 'mount', key: 'mount' },
 { title: '总大小', dataIndex: 'total', key: 'total', width: 120 },
 { title: '已用', dataIndex: 'used', key: 'used', width: 120 },
 { title: '可用', dataIndex: 'free', key: 'free', width: 120 },
 { title: '使用率', dataIndex: 'usage', key: 'usage', width: 150 }
];
const disks = ref<Disk[]>([
 { mount: '/', total: '50G', used: '35G', free: '15G', usage: 70 },
 { mount: '/home', total: '100G', used: '45G', free: '55G', usage: 45 },
 { mount: '/var', total: '20G', used: '12G', free: '8G', usage: 60 },
 { mount: '/tmp', total: '10G', used: '2G', free: '8G', usage: 20 }
]);
const getMetricColor = (value: number, threshold = 80) => {
 if (value >= threshold)
 return '#f5222d';
 if (value >= threshold * 0.75)
 return '#faad14';
 return '#52c41a';
};
const handleRefresh = () => {
 currentMetrics.cpu = Math.floor(Math.random() * 50) + 30;
 currentMetrics.memory = Math.floor(Math.random() * 40) + 40;
 currentMetrics.disk = Math.floor(Math.random() * 30) + 60;
 currentMetrics.network = Math.floor(Math.random() * 50) + 20;
 message.success('刷新成功');
};
</script>