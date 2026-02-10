import mitt from 'mitt'

// 定义事件类型
type Events = {
  'plugin:loaded': string // 插件加载完成
  'plugin:unloaded': string // 插件卸载完成
  'plugin:status:changed': { pluginId: string; status: string } // 插件状态变更
  'plugin:message': { pluginId: string; message: any } // 插件消息
}

// 创建事件总线实例
export const eventBus = mitt<Events>()

export default eventBus