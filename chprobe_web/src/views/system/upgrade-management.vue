<template>
  <div>
    <ClientManagement v-if="activeTab === 'client'" />
    <ServerUpgrade v-if="activeTab === 'server'" />
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import ClientManagement from './client-management.vue'
import ServerUpgrade from './server-upgrade.vue'

const activeTab = inject('activeTopTab', 'client')
</script>
