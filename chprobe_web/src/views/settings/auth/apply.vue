<template>
  <div style="padding: 8px;">
    <a-card style="margin-bottom: 16px;">
      <div style="text-align: center;">
        <h3 style="margin-bottom: 16px;">联系信息</h3>
        <a-row :gutter="[32, 16]" justify="center">
          <a-col>
            <div style="margin-bottom: 8px;">
              <span style="color: #8c8c8c;">邮箱：</span>
              <span>sales@chprobe.com</span>
            </div>
          </a-col>
          <a-col>
            <div style="margin-bottom: 8px;">
              <span style="color: #8c8c8c;">技术支持：</span>
              <span>support@chprobe.com</span>
            </div>
          </a-col>
        </a-row>
      </div>
    </a-card>

    <a-row :gutter="[16, 16]" justify="center">
      <a-col v-for="version in versions" :key="version.key" :xs="24" :sm="12" :md="12" :lg="5" :xl="4">
        <a-card hoverable :class="{ 'selected-card': selectedVersion === version.key }" @click="selectVersion(version.key)" class="version-card">
          <template #title>
            <div style="text-align: center;">
              <h3>{{ version.name }}</h3>
              <div :style="`font-size: 24px; font-weight: bold; color: ${version.color}; margin-top: 8px;`">
                {{ version.price }}
              </div>
              <div style="font-size: 12px; color: #8c8c8c; margin-top: 4px;">
                {{ version.position }}
              </div>
            </div>
          </template>
          <a-list size="small" :data-source="version.features" style="margin-top: 16px;">
            <template #renderItem="{ item }">
              <a-list-item>
                <CheckOutlined style="color: #52c41a; margin-right: 8px;" />
                {{ item }}
              </a-list-item>
            </template>
          </a-list>
          <div style="text-align: center; margin-top: 24px;">
            <a-button type="primary" block @click="contactUs(version.key)">联系我们</a-button>
          </div>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { CheckOutlined } from '@ant-design/icons-vue'

const selectedVersion = ref('normal')

const versions = ref([
  {
    key: 'normal',
    name: '普通版',
    price: '免费',
    position: '个人开发者',
    color: '#1890ff',
    features: [
      '基础主机监控',
      '基础容器管理',
      '最多支持1台主机',
      '社区支持'
    ]
  },
  {
    key: 'pro',
    name: 'Pro版',
    price: '¥5/月',
    position: '中小企业',
    color: '#722ed1',
    features: [
      '包含普通版所有功能',
      '高级主机监控',
      '高级容器管理',
      '高级告警策略',
      '最多支持3台主机',
      '插件市场基础功能',
      '邮件支持'
    ]
  },
  {
    key: 'plus',
    name: 'Plus版',
    price: '¥10/月',
    position: '快速成长企业',
    color: '#fa8c16',
    features: [
      '包含Pro版所有功能',
      '性能分析报告',
      '自定义仪表盘',
      '最多支持20台主机',
      '插件市场高级功能',
      'API接口访问',
      '优先技术支持'
    ]
  },
  {
    key: 'max',
    name: 'Max版',
    price: '¥100/月',
    color: '#f5222d',
    features: [
      '包含Plus版所有功能',
      '最多支持200台主机',
      '高级安全审计',
      '多租户支持',
      '私有化部署',
      '专属技术顾问',
      '7x24小时支持'
    ]
  },
  {
    key: 'enterprise',
    name: '企业版',
    price: '定制价格',
    position: '超大型企业专属定制',
    color: '#13c2c2',
    features: [
      '包含Max版所有功能',
      '完全定制化开发',
      '专属部署方案',
      'SLA服务保障',
      '现场技术支持',
      '定期培训服务',
      '专属客户经理'
    ]
  }
])

const selectVersion = (version: string) => {
  selectedVersion.value = version
}

const contactUs = (version: string) => {
  const versionData = versions.value.find(v => v.key === version)
  const email = 'sales@chprobe.com'
  const subject = encodeURIComponent(`申请${versionData?.name}授权`)
  const body = encodeURIComponent(`您好，我想申请${versionData?.name}授权。\n\n请提供相关授权信息。`)
  
  window.location.href = `mailto:${email}?subject=${subject}&body=${body}`
}
</script>

<style scoped>
.selected-card {
  border: 2px solid #1890ff;
  box-shadow: 0 0 10px rgba(24, 144, 255, 0.3);
}

.version-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.version-card :deep(.ant-card-body) {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.version-card :deep(.ant-list) {
  flex: 1;
}

:deep(.ant-card-head) {
  background: #fafafa;
}

:deep(.ant-list-item) {
  padding: 8px 0;
  border: none;
}

:deep(.ant-list-item-action) {
  margin-left: auto;
}
</style>