package pluginmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/ricky97gr/chprobe/chprobe_common/utils"
	shared "github.com/ricky97gr/chprobe/chprobe_common/plugin/shared"
)

type PluginStatus string

const (
	PluginStatusStopped PluginStatus = "stopped"
	PluginStatusRunning PluginStatus = "running"
	PluginStatusError   PluginStatus = "error"
)

// MenuItem 菜单项
type MenuItem struct {
	ID       string      `json:"id"`
	Title    string      `json:"title"`
	Path     string      `json:"path"`
	Icon     string      `json:"icon"`
	Order    int         `json:"order"`
	Children []MenuItem  `json:"children,omitempty"`
}

// PluginMeta 插件元信息（从meta.json读取）
type PluginMeta struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	License     string   `json:"license,omitempty"`
	PluginType []string `json:"plugin_type,omitempty"`
	Main       string   `json:"main,omitempty"`
	ESModule    string   `json:"esModule,omitempty"`
	Menu       MenuItem `json:"menu,omitempty"`
	Routes     []Route  `json:"routes,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	ApiPrefix  string   `json:"apiPrefix,omitempty"`
}

// Route 路由配置
type Route struct {
	Path      string `json:"path"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	Component string `json:"component"`
}

// PluginWebConfig 插件Web配置
type PluginWebConfig struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Menu       MenuItem `json:"menu,omitempty"`
	Routes     []Route  `json:"routes,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	ApiPrefix  string   `json:"apiPrefix,omitempty"`
}

type ManagedPlugin struct {
	ID          string
	Name        string
	Version     string
	Description string
	Status      PluginStatus
	Client      *plugin.Client
	Service     shared.PluginService
	Config      map[string]interface{}
	Process     *os.Process
	StartTime   time.Time
	Routes      []*shared.RouteInfo
	Meta        *PluginMeta
	WebConfig   *PluginWebConfig
}

type Manager struct {
	plugins   map[string]*ManagedPlugin
	mu        sync.RWMutex
	pluginDir string
}

func NewManager(pluginDir string) *Manager {
	return &Manager{
		plugins:   make(map[string]*ManagedPlugin),
		pluginDir: pluginDir,
	}
}

// readPluginMeta 读取插件的meta.json文件
func (m *Manager) readPluginMeta(pluginID string) (*PluginMeta, error) {
	metaPath := filepath.Join(m.pluginDir, pluginID, "meta.json")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read meta.json: %w", err)
	}

	var meta PluginMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal meta.json: %w", err)
	}

	return &meta, nil
}

// readPluginWebConfig 读取插件的Web配置（从meta.json或web/plugin.json）
func (m *Manager) readPluginWebConfig(pluginID string) (*PluginWebConfig, error) {
	metaPath := filepath.Join(m.pluginDir, pluginID, "meta.json")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read meta.json: %w", err)
	}

	var meta PluginMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal meta.json: %w", err)
	}

	webConfig := &PluginWebConfig{
		Name:        meta.Name,
		Description: meta.Description,
		Menu:       meta.Menu,
		Routes:     meta.Routes,
		Permissions: meta.Permissions,
		ApiPrefix:  meta.ApiPrefix,
	}

	return webConfig, nil
}

func (m *Manager) StartPlugin(ctx context.Context, pluginID, command string, args []string, config map[string]interface{}) (*ManagedPlugin, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.plugins[pluginID]; exists {
		return nil, fmt.Errorf("plugin %s already exists", pluginID)
	}

	utils.Logger.Infof("Starting plugin: command=%s, args=%v\n", command, args)
	
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "BASIC_PLUGIN",
			MagicCookieValue: "hello",
		},
		Plugins: map[string]plugin.Plugin{
			"grpc_plugin": &shared.GRPCPlugin{},
		},
		Cmd:              exec.CommandContext(ctx, command, args...),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		// 设置插件启动超时时间
		StartTimeout: 10 * time.Second,
	})

	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("failed to get RPC client: %w", err)
	}
	utils.Logger.Infof("RPC client connected successfully\n")

	raw, err := rpcClient.Dispense("grpc_plugin")
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("failed to dispense plugin: %w", err)
	}
	utils.Logger.Infof("Plugin dispensed successfully\n")

	pluginService, ok := raw.(shared.PluginService)
	if !ok {
		client.Kill()
		return nil, fmt.Errorf("failed to get plugin service")
	}
	utils.Logger.Infof("Plugin service obtained successfully\n")

	utils.Logger.Infof("Calling GetPluginInfo...\n")
	info, err := pluginService.GetPluginInfo(ctx)
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("failed to get plugin info: %w", err)
	}
	utils.Logger.Infof("GetPluginInfo succeeded: id=%s, name=%s\n", info.Id, info.Name)

	utils.Logger.Infof("Calling GetRoutes...\n")
	routes, err := pluginService.GetRoutes(ctx)
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("failed to get plugin routes: %w", err)
	}
	utils.Logger.Infof("GetRoutes succeeded: %d routes\n", len(routes))

	utils.Logger.Infof("Calling HealthCheck...\n")
	health, err := pluginService.HealthCheck(ctx)
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("plugin health check failed: %w", err)
	}
	utils.Logger.Infof("HealthCheck succeeded: status=%s\n", health.Status)

	if health.Status != "healthy" {
		client.Kill()
		return nil, fmt.Errorf("plugin is not healthy: %s", health.Status)
	}

	utils.Logger.Infof("All plugin initialization steps completed successfully, plugin is running\n")

	// 读取插件配置文件
	meta, err := m.readPluginMeta(pluginID)
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("failed to read plugin meta: %w", err)
	}

	webConfig, err := m.readPluginWebConfig(pluginID)
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("failed to read plugin web config: %w", err)
	}

	managedPlugin := &ManagedPlugin{
		ID:          pluginID,
		Name:        info.Name,
		Version:     info.Version,
		Description: info.Description,
		Status:      PluginStatusRunning,
		Client:      client,
		Service:     pluginService,
		Config:      config,
		StartTime:   time.Now(),
		Routes:      routes,
		Meta:        meta,
		WebConfig:   webConfig,
	}

	m.plugins[pluginID] = managedPlugin

	utils.Logger.Infof("Plugin %s started successfully and registered\n", pluginID)

	return managedPlugin, nil
}

func (m *Manager) StopPlugin(pluginID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	plugin, exists := m.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin %s not found", pluginID)
	}

	plugin.Client.Kill()
	plugin.Status = PluginStatusStopped

	delete(m.plugins, pluginID)

	return nil
}

func (m *Manager) RestartPlugin(ctx context.Context, pluginID, command string, args []string, config map[string]interface{}) (*ManagedPlugin, error) {
	if err := m.StopPlugin(pluginID); err != nil {
		return nil, fmt.Errorf("failed to stop plugin: %w", err)
	}

	return m.StartPlugin(ctx, pluginID, command, args, config)
}

func (m *Manager) GetPlugin(pluginID string) (*ManagedPlugin, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugin, exists := m.plugins[pluginID]
	return plugin, exists
}

func (m *Manager) GetAllPlugins() []*ManagedPlugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugins := make([]*ManagedPlugin, 0, len(m.plugins))
	for _, plugin := range m.plugins {
		plugins = append(plugins, plugin)
	}

	return plugins
}

func (m *Manager) GetPluginRoutes(pluginID string) ([]*shared.RouteInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugin, exists := m.plugins[pluginID]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", pluginID)
	}

	return plugin.Routes, nil
}

func (m *Manager) QueryRoute(pluginID, path, method string) (*shared.RouteInfo, error) {
	routes, err := m.GetPluginRoutes(pluginID)
	if err != nil {
		return nil, err
	}

	for _, route := range routes {
		if route.Path == path && (method == "" || route.Method == method) {
			return route, nil
		}
	}

	return nil, fmt.Errorf("route not found")
}

func (m *Manager) HandleRequest(ctx context.Context, pluginID string, method string, path string, data map[string]string) (*shared.HandleResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugin, exists := m.plugins[pluginID]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", pluginID)
	}

	return plugin.Service.HandleRequest(ctx, method, path, data)
}

func (m *Manager) GetPluginHealth(pluginID string) (*shared.HealthStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugin, exists := m.plugins[pluginID]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", pluginID)
	}

	return plugin.Service.HealthCheck(context.Background())
}

func (m *Manager) GetPluginMetrics(pluginID string) (map[string]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugin, exists := m.plugins[pluginID]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", pluginID)
	}

	return plugin.Service.GetMetrics(context.Background())
}

func (m *Manager) GetPluginMetricsSummary(pluginID string) (map[string]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugin, exists := m.plugins[pluginID]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", pluginID)
	}

	return plugin.Service.GetMetricsSummary(context.Background())
}

func (m *Manager) ShutdownAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, plugin := range m.plugins {
		plugin.Client.Kill()
	}

	m.plugins = make(map[string]*ManagedPlugin)
}
