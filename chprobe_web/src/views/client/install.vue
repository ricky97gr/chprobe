<template>
  <div>
    <a-card>
      <a-alert
        message="安装客户端前请确保"
        description="1. 目标主机已安装 curl 或 wget；2. 目标主机可以访问本服务器；3. 具有 root 或 sudo 权限"
        type="info"
        show-icon
        style="margin-bottom: 24px;"
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
                <a-input
                  v-model:value="formState.serverIp"
                  placeholder="请输入服务器IP地址"
                  style="width: 300px;"
                  @change="updateInstallCmds"
                />
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

      <a-card style="margin-top: 16px;">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-card size="small">
              <a-descriptions :column="1">
                <a-descriptions-item label="检查版本">
                  <code>chprobe-agent --version</code>
                </a-descriptions-item>
                <a-descriptions-item label="检查服务状态">
                  <code>systemctl is-active chprobe-agent</code>
                </a-descriptions-item>
                <a-descriptions-item label="查看日志">
                  <code>journalctl -u chprobe-agent -f</code>
                </a-descriptions-item>
              </a-descriptions>
            </a-card>
          </a-col>
          <a-col :span="12">
            <a-card size="small">
              <a-descriptions :column="1">
                <a-descriptions-item label="检查指标端点">
                  <code>curl http://localhost:9100/metrics</code>
                </a-descriptions-item>
                <a-descriptions-item label="查看主机列表">
                  <a-button type="link" @click="goToHostList">
                    前往主机管理页面
                  </a-button>
                </a-descriptions-item>
              </a-descriptions>
            </a-card>
          </a-col>
        </a-row>
      </a-card>
    </a-card>

    <a-card style="margin-top: 16px;">
      <a-collapse v-model:activeKey="faqActiveKey">
        <a-collapse-panel key="1" header="安装失败怎么办？">
          <p>1. 检查网络连接是否正常</p>
          <p>2. 查看安装日志：<code>tail -f /var/log/chprobe-install.log</code></p>
          <p>3. 手动下载安装包并安装</p>
        </a-collapse-panel>
        <a-collapse-panel key="2" header="客户端无法连接到服务器？">
          <p>1. 检查服务器地址和端口是否正确</p>
          <p>2. 检查防火墙规则</p>
          <p>3. 查看客户端日志：<code>journalctl -u chprobe-agent -f</code></p>
        </a-collapse-panel>
        <a-collapse-panel key="3" header="如何升级客户端？">
          <p>执行 <code>chprobe-agent --upgrade</code> 或重新运行安装脚本</p>
        </a-collapse-panel>
        <a-collapse-panel key="4" header="如何卸载客户端？">
          <p>执行 <code>curl -fsSL http://your-server/uninstall | bash</code></p>
        </a-collapse-panel>
        <a-collapse-panel key="5" header="Docker安装常见问题">
          <p>1. 容器无法启动：检查网络连接和端口映射</p>
          <p>2. 权限问题：确保Docker守护进程正在运行</p>
          <p>3. 数据持久化：如需持久化配置，可挂载配置文件：<code>-v /path/to/config.yml:/etc/chprobe-agent/config.yml</code></p>
          <p>4. 查看容器日志：<code>docker logs chprobe-agent</code></p>
        </a-collapse-panel>
      </a-collapse>
    </a-card>
  </div>
</template>

<script setup lang="ts">import { ref, reactive, computed } from 'vue';
import { message } from 'ant-design-vue';
import { useRouter } from 'vue-router';
import { CopyOutlined } from '@ant-design/icons-vue';
const router = useRouter();

const formState = reactive({
  serverIp: '',
  serverPort: 8080
});

const serverUrl = computed(() => {
  if (!formState.serverIp) {
    return 'http://your-server-ip:8080';
  }
  return `http://${formState.serverIp}:${formState.serverPort}`;
});

const quickInstallCmd = computed(() => {
  return `curl -fsSL ${serverUrl.value}/install | bash -s -- --server ${serverUrl.value}`;
});

const dockerInstallCmd = computed(() => {
  return `docker run -d --name chprobe-agent \
  -e SERVER_ADDRESS=${serverUrl.value} \
  -e SERVER_TOKEN="your-token-here" \
  -p 9100:9100 \
  --restart always \
  chprobe/agent:latest`;
});

const configFileContent = computed(() => {
  return `# /etc/chprobe-agent/config.yml
server:
  address: ${serverUrl.value}
  token: "your-token-here"
  interval: 10s

metrics:
  enabled: true
  port: 9100

logs:
  enabled: true
  path: /var/log/chprobe-agent
  level: info

features:
  cpu: true
  memory: true
  disk: true
  network: true
  process: true`;
});

const activeConfigKey = ref<string[]>(['config']);
const faqActiveKey = ref<string[]>([]);

const copyQuickCmd = () => {
  navigator.clipboard.writeText(quickInstallCmd.value).then(() => {
    message.success('命令已复制到剪贴板');
  });
};

const copyDockerCmd = () => {
  navigator.clipboard.writeText(dockerInstallCmd.value).then(() => {
    message.success('命令已复制到剪贴板');
  });
};

const updateInstallCmds = () => {
  // 命令会通过computed自动更新
};

const goToHostList = () => {
  router.push('/host/list');
};
</script>