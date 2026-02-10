<template>
  <div>
    <a-card>
      <a-table :columns="clientColumns" :data-source="agentList" :loading="loading" row-key="uuid">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'registerTime'">
            {{ formatTime(record.registerTime) }}
          </template>
          <template v-if="column.key === 'lastHeartTime'">
            {{ formatTime(record.lastHeartTime) }}
          </template>
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 'online' ? 'green' : 'red'">
              {{ record.status === 'online' ? '在线' : '离线' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'action'">
            <a-space>
              <a-button 
                type="primary" 
                size="small" 
                @click="upgradeClient(record.uuid)"
                :disabled="record.status !== 'online'"
              >
                升级
              </a-button>
              <a-button 
                danger 
                size="small" 
                @click="uninstallClient(record.uuid)"
                :disabled="record.status !== 'online'"
              >
                卸载
              </a-button>
            </a-space>
          </template>
        </template>
      </a-table>

      <div style="margin-top: 16px; text-align: right;">
        <a-button 
          type="primary" 
          @click="batchUpgradeClient"
          :disabled="needUpgradeClientCount === 0"
        >
          批量升级所有在线客户端
        </a-button>
      </div>
    </a-card>

    <a-alert
      message="客户端安装指南"
      description="1. 下载客户端安装包；2. 运行安装程序；3. 按照提示完成安装；4. 启动客户端并连接到服务器。"
      type="info"
      show-icon
      style="margin-top: 24px; margin-bottom: 24px;"
    />

    <a-steps :current="0" status="process" style="margin-bottom: 24px;">
      <a-step title="选择服务器IP" />
      <a-step title="执行安装命令" />
      <a-step title="验证安装" />
    </a-steps>

    <a-row :gutter="16" style="margin-bottom: 24px;">
      <a-col :span="24">
        <a-card size="small">
          <a-form :model="formState" layout="inline">
            <a-form-item label="服务器IP">
              <a-select
                v-model:value="formState.serverIp"
                placeholder="请选择服务器IP"
                style="width: 300px;"
                show-search
                :filter-option="filterOption"
                @focus="fetchServerIPs"
                @change="updateInstallCmds"
              >
                <a-select-option v-for="ip in serverIPs" :key="ip" :value="ip">
                  {{ ip }}
                </a-select-option>
              </a-select>
            </a-form-item>
            <a-form-item label="端口">
              <a-input-number
                v-model:value="formState.serverPort"
                :min="1"
                :max="65535"
                style="width: 100px;"
                @change="updateInstallCmds"
              />
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16">
      <a-col :xs="24" :lg="12">
        <a-card size="small">
          <a-descriptions :column="1" bordered>
            <a-descriptions-item label="适用场景">
              快速部署，自动配置
            </a-descriptions-item>
            <a-descriptions-item label="支持系统">
              <a-tag>Ubuntu</a-tag>
              <a-tag>CentOS</a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="安装命令">
              <a-input
                :value="quickInstallCmd"
                readonly
                style="margin-bottom: 8px;"
              />
              <a-button type="primary" @click="copyQuickCmd" block>
                <template #icon>
                  <CopyOutlined />
                </template>
                复制命令
              </a-button>
            </a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>

      <a-col :xs="24" :lg="12">
        <a-card size="small">
          <a-descriptions :column="1" bordered>
            <a-descriptions-item label="适用场景">
              容器化部署，隔离环境
            </a-descriptions-item>
            <a-descriptions-item label="支持系统">
              <a-tag>任何支持Docker的系统</a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="安装命令">
              <a-input
                :value="dockerInstallCmd"
                readonly
                style="margin-bottom: 8px;"
              />
              <a-button type="primary" @click="copyDockerCmd" block>
                <template #icon>
                  <CopyOutlined />
                </template>
                复制命令
              </a-button>
            </a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>
    </a-row>

    <a-card style="margin-top: 16px;">
      <a-collapse v-model:activeKey="activeConfigKey">
        <a-collapse-panel key="config" header="配置文件说明">
          <pre style="background: #f0f0f0; padding: 16px;">{{ configFileContent }}</pre>
        </a-collapse-panel>
        <a-collapse-panel key="service" header="服务管理">
          <a-descriptions :column="1" bordered>
            <a-descriptions-item label="启动服务">
              <code>systemctl start chprobe-agent</code>
            </a-descriptions-item>
            <a-descriptions-item label="停止服务">
              <code>systemctl stop chprobe-agent</code>
            </a-descriptions-item>
            <a-descriptions-item label="重启服务">
              <code>systemctl restart chprobe-agent</code>
            </a-descriptions-item>
            <a-descriptions-item label="查看状态">
              <code>systemctl status chprobe-agent</code>
            </a-descriptions-item>
            <a-descriptions-item label="开机自启">
              <code>systemctl enable chprobe-agent</code>
            </a-descriptions-item>
          </a-descriptions>
        </a-collapse-panel>
        <a-collapse-panel key="docker" header="Docker管理">
          <a-descriptions :column="1" bordered>
            <a-descriptions-item label="启动容器">
              <code>docker start chprobe-agent</code>
            </a-descriptions-item>
            <a-descriptions-item label="停止容器">
              <code>docker stop chprobe-agent</code>
            </a-descriptions-item>
            <a-descriptions-item label="重启容器">
              <code>docker restart chprobe-agent</code>
            </a-descriptions-item>
            <a-descriptions-item label="查看状态">
              <code>docker ps -a | grep chprobe-agent</code>
            </a-descriptions-item>
            <a-descriptions-item label="查看日志">
              <code>docker logs -f chprobe-agent</code>
            </a-descriptions-item>
            <a-descriptions-item label="删除容器">
              <code>docker rm -f chprobe-agent</code>
            </a-descriptions-item>
          </a-descriptions>
        </a-collapse-panel>
        <a-collapse-panel key="firewall" header="防火墙配置">
          <a-alert
            message="如果主机启用了防火墙，请开放以下端口"
            description="TCP 端口 9100（用于指标采集）和 UDP 端口 9100（可选，用于日志采集）"
            type="warning"
            show-icon
          />
          <a-space direction="vertical" style="margin-top: 16px; width: 100%;">
            <div>
              <strong>UFW (Ubuntu/Debian):</strong>
              <pre style="background: #f0f0f0; padding: 8px; margin-top: 8px;">ufw allow 9100/tcp
ufw reload</pre>
            </div>
            <div>
              <strong>Firewalld (CentOS/RHEL):</strong>
              <pre style="background: #f0f0f0; padding: 8px; margin-top: 8px;">firewall-cmd --permanent --add-port=9100/tcp
firewall-cmd --reload</pre>
            </div>
          </a-space>
        </a-collapse-panel>
      </a-collapse>
    </a-card>

    <a-collapse style="margin-top: 24px; margin-bottom: 24px;">
      <a-collapse-panel key="faq1" header="安装失败怎么办？">
        <p>1. 检查网络连接是否正常</p>
        <p>2. 查看安装日志：<code>tail -f /var/log/chprobe-install.log</code></p>
        <p>3. 手动下载安装包并安装</p>
      </a-collapse-panel>

      <a-collapse-panel key="faq2" header="客户端无法连接到服务器？">
        <p>1. 检查服务器地址和端口是否正确</p>
        <p>2. 检查防火墙规则</p>
        <p>3. 查看客户端日志：<code>journalctl -u chprobe-agent -f</code></p>
      </a-collapse-panel>

      <a-collapse-panel key="faq3" header="如何升级客户端？">
        <p>执行 <code>chprobe-agent --upgrade</code> 或重新运行安装脚本</p>
      </a-collapse-panel>

      <a-collapse-panel key="faq4" header="如何卸载客户端？">
        <p>执行 <code>curl -fsSL http://your-server/uninstall | bash</code></p>
      </a-collapse-panel>

      <a-collapse-panel key="faq5" header="Docker安装常见问题">
        <p>1. 容器无法启动：检查网络连接和端口映射</p>
        <p>2. 权限问题：确保Docker守护进程正在运行</p>
        <p>3. 数据持久化：如需持久化配置，可挂载配置文件：<code>-v /path/to/config.yml:/etc/chprobe-agent/config.yml</code></p>
        <p>4. 查看容器日志：<code>docker logs chprobe-agent</code></p>
      </a-collapse-panel>
    </a-collapse>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { useRouter } from 'vue-router'
import { CopyOutlined } from '@ant-design/icons-vue'
import { getAgentList, deleteAgent, getServerIPs, type Agent } from '@/api/index'

const router = useRouter()

const loading = ref(false)
const agentList = ref<Agent[]>([])

const currentVersion = 'v1.0.0'
const latestVersion = 'v1.1.0'
const needUpgradeClientCount = computed(() => agentList.value.filter(a => a.status === 'online').length)

const clientColumns = [
  { title: '主机名', dataIndex: 'hostName', key: 'hostName', width: 150 },
  { title: 'IP地址', dataIndex: 'ip', key: 'ip', width: 120 },
  { title: '客户端类型', dataIndex: 'clientType', key: 'clientType', width: 100 },
  { title: '操作系统', dataIndex: 'os', key: 'os', width: 120 },
  { title: '系统类型', dataIndex: 'osType', key: 'osType', width: 80 },
  { title: '架构', dataIndex: 'arch', key: 'arch', width: 80 },
  { title: '内核版本', dataIndex: 'kernelVersion', key: 'kernelVersion', width: 120 },
  { title: '客户端版本', dataIndex: 'version', key: 'version', width: 100 },
  { title: '状态', dataIndex: 'status', key: 'status', width: 80 },
  { title: '注册时间', dataIndex: 'registerTime', key: 'registerTime', width: 180 },
  { title: '最后心跳', dataIndex: 'lastHeartTime', key: 'lastHeartTime', width: 180 },
  { 
    title: '操作', 
    dataIndex: 'action', 
    key: 'action', 
    width: 150,
    fixed: 'right'
  }
]

const formatTime = (timestamp: number) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const fetchAgentList = async () => {
  loading.value = true
  try {
    console.log('Fetching agent list...')
    const response = await getAgentList()
    console.log('Agent list response:', response)
    agentList.value = response.result
    console.log('Agent list:', agentList.value)
  } catch (error: any) {
    console.error('Failed to fetch agent list:', error)
    message.error(error.message || '获取客户端列表失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchAgentList()
})

const upgradeClient = async (uuid: string) => {
  message.info('升级功能开发中...')
}

const uninstallClient = async (uuid: string) => {
  try {
    await deleteAgent({ uuid })
    message.success('客户端卸载成功')
    fetchAgentList()
  } catch (error: any) {
    message.error(error.message || '卸载失败')
  }
}

const batchUpgradeClient = () => {
  message.info('批量升级功能开发中...')
}

const goToHostList = () => {
  router.push('/host/list')
}

const formState = reactive({
  serverIp: '',
  serverPort: 8080
})

const serverIPs = ref<string[]>([])

const activeConfigKey = ref<string[]>([])

const fetchServerIPs = async () => {
  try {
    const response = await getServerIPs()
    serverIPs.value = response.result
  } catch (error: any) {
    console.error('Failed to fetch server IPs:', error)
  }
}

const filterOption = (input: string, option: any) => {
  return option.value.toLowerCase().indexOf(input.toLowerCase()) >= 0
}

const serverUrl = computed(() => {
  if (!formState.serverIp) {
    return 'your-server'
  }
  return `${formState.serverIp}:${formState.serverPort}`
})

const quickInstallCmd = computed(() => {
  return `curl -fsSL http://${serverUrl.value}/install | bash`
})

const dockerInstallCmd = computed(() => {
  return `docker run -d --name chprobe-agent --restart=always -p 9100:9100 -e SERVER_URL=http://${serverUrl.value} chprobe-agent:latest`
})

const configFileContent = `server:
  url: http://${serverUrl.value}
  port: 9100

logging:
  level: info
  file: /var/log/chprobe-agent.log

metrics:
  enabled: true
  port: 9100`

const copyQuickCmd = () => {
  navigator.clipboard.writeText(quickInstallCmd.value)
  message.success('命令已复制到剪贴板')
}

const copyDockerCmd = () => {
  navigator.clipboard.writeText(dockerInstallCmd.value)
  message.success('命令已复制到剪贴板')
}

const updateInstallCmds = () => {
}
</script>
