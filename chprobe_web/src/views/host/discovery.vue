<template>
  <div>
    <a-row :gutter="[16, 16]">
      <a-col :xs="24" :lg="8">
        <a-card title="发现配置">
          <a-form :model="discoveryForm" :label-col="{ span: 8 }" :wrapper-col="{ span: 16 }">
            <a-form-item label="IP范围">
              <a-input v-model:value="discoveryForm.ipRange" placeholder="例如: 192.168.1.1-100" />
            </a-form-item>
            <a-form-item label="端口范围">
              <a-input v-model:value="discoveryForm.portRange" placeholder="例如: 22,80,443" />
            </a-form-item>
            <a-form-item label="超时时间">
              <a-input-number v-model:value="discoveryForm.timeout" :min="1" :max="60" style="width: 100%;" />
              <span style="margin-left: 8px;">秒</span>
            </a-form-item>
            <a-form-item label="并发数">
              <a-input-number v-model:value="discoveryForm.concurrent" :min="1" :max="100" style="width: 100%;" />
            </a-form-item>
            <a-form-item label="扫描协议">
              <a-checkbox-group v-model:value="discoveryForm.protocols">
                <a-checkbox value="ssh">SSH</a-checkbox>
                <a-checkbox value="http">HTTP</a-checkbox>
                <a-checkbox value="https">HTTPS</a-checkbox>
              </a-checkbox-group>
            </a-form-item>
            <a-form-item :wrapper-col="{ offset: 8, span: 16 }">
              <a-button type="primary" @click="handleStart" :loading="scanning" block>
                <template #icon>
                  <PlayCircleOutlined v-if="!scanning" />
                  <LoadingOutlined v-else />
                </template>
                {{ scanning ? '扫描中...' : '开始扫描' }}
              </a-button>
              <a-button @click="handleStop" v-if="scanning" style="margin-top: 8px;" block>
                停止扫描
              </a-button>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>

      <a-col :xs="24" :lg="16">
        <a-card title="发现结果" :bordered="false">
          <a-alert
            v-if="!scanning && discoveredHosts.length === 0"
            message="暂无发现结果"
            description="配置扫描参数后点击开始扫描"
            type="info"
            show-icon
            style="margin-bottom: 16px;"
          />

          <a-table
            v-else
            :columns="resultColumns"
            :data-source="discoveredHosts"
            :loading="scanning"
            row-key="ip"
            :pagination="{ pageSize: 10 }"
          >
            <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'status'">
              <a-tag :color="record.status === 'online' ? 'success' : 'default'">
                {{ record.status === 'online' ? '在线' : '离线' }}
              </a-tag>
            </template>
            <template v-else-if="column.key === 'services'">
              <a-tag v-for="service in record.services" :key="service" color="blue">
                {{ service }}
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
                    <div class="action-item" @click="handleAddHost(record)">
                      添加到主机列表
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

          <a-progress
            v-if="scanning"
            :percent="progress"
            :status="progress === 100 ? 'success' : 'active'"
            style="margin-top: 16px;"
          />
        </a-card>
      </a-col>
    </a-row>

    <a-card title="快速扫描" style="margin-top: 16px;">
      <a-space wrap>
        <a-button @click="handleQuickScan('local')">
          <template #icon>
            <EnvironmentOutlined />
          </template>
          本地网络
        </a-button>
        <a-button @click="handleQuickScan('common')">
          <template #icon>
            <CloudOutlined />
          </template>
          常用端口
        </a-button>
        <a-button @click="handleQuickScan('custom')">
          <template #icon>
            <SettingOutlined />
          </template>
          自定义
        </a-button>
      </a-space>
    </a-card>
  </div>
</template>

<script setup lang="ts">import { ref, reactive } from 'vue';
import { message } from 'ant-design-vue';
import { PlayCircleOutlined, LoadingOutlined, EnvironmentOutlined, CloudOutlined, SettingOutlined } from '@ant-design/icons-vue';
interface DiscoveredHost {
 ip: string;
 hostname: string;
 status: 'online' | 'offline';
 services: string[];
 latency: number;
 os: string;
}
const resultColumns = [
 { title: 'IP地址', dataIndex: 'ip', key: 'ip' },
 { title: '主机名', dataIndex: 'hostname', key: 'hostname' },
 { title: '状态', dataIndex: 'status', key: 'status', width: 80 },
 { title: '服务', dataIndex: 'services', key: 'services' },
 { title: '延迟', dataIndex: 'latency', key: 'latency', width: 80, render: (text: number) => `${text}ms` },
 { title: '操作系统', dataIndex: 'os', key: 'os', width: 120 },
 { title: '操作', key: 'action', width: 150 }
];
const discoveryForm = reactive({
 ipRange: '192.168.1.1-100',
 portRange: '22,80,443,8080',
 timeout: 5,
 concurrent: 10,
 protocols: ['ssh', 'http']
});
const scanning = ref(false);
const progress = ref(0);
const discoveredHosts = ref<DiscoveredHost[]>([]);
const handleStart = () => {
 if (!discoveryForm.ipRange) {
 message.error('请输入IP范围');
 return;
 }
 scanning.value = true;
 progress.value = 0;
 discoveredHosts.value = [];
 const interval = setInterval(() => {
 progress.value += Math.floor(Math.random() * 10) + 5;
 if (progress.value >= 100) {
 progress.value = 100;
 clearInterval(interval);
 scanning.value = false;
 discoveredHosts.value = [
 { ip: '192.168.1.10', hostname: 'router.local', status: 'online', services: ['HTTP', 'HTTPS'], latency: 1, os: 'Linux' },
 { ip: '192.168.1.20', hostname: 'server-01.local', status: 'online', services: ['SSH', 'HTTP'], latency: 2, os: 'Ubuntu 22.04' },
 { ip: '192.168.1.30', hostname: 'server-02.local', status: 'online', services: ['SSH'], latency: 3, os: 'CentOS 7.9' },
 { ip: '192.168.1.40', hostname: 'workstation-01.local', status: 'online', services: ['HTTP', 'HTTPS'], latency: 5, os: 'Windows 11' },
 { ip: '192.168.1.50', hostname: 'storage.local', status: 'online', services: ['HTTP', 'HTTPS'], latency: 8, os: 'FreeBSD' },
 { ip: '192.168.1.60', hostname: '', status: 'offline', services: [], latency: 0, os: 'Unknown' }
 ];
 message.success(`发现完成，共找到 ${discoveredHosts.value.length} 台主机`);
 }
 }, 300);
};
const handleStop = () => {
 message.warning('扫描已停止');
};
const handleAddHost = (record: DiscoveredHost) => {
 message.success(`已将主机 ${record.ip} 添加到主机列表`);
};
const handleQuickScan = (type: string) => {
  if (type === 'local') {
    discoveryForm.ipRange = '192.168.1.1-254';
    discoveryForm.portRange = '22,80,443';
  }
  else if (type === 'common') {
    discoveryForm.portRange = '21,22,23,80,443,3306,3389,8080';
  }
  message.info(`已设置${type === 'local' ? '本地网络' : type === 'common' ? '常用端口' : '自定义'}扫描配置`);
};
</script>

<style scoped>
.btn-action {
  color: #1890ff !important;
}

.action-popover {
  padding: 2px 0;
  width: 120px;
}

.action-item {
  padding: 3px 4px;
  font-size: 12px;
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
</style>