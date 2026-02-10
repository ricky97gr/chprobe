package pluginmanager

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/hashicorp/go-plugin"
	shared "github.com/ricky97gr/chprobe/chprobe_common/plugin/shared"
)

type PluginStatus string

const (
	PluginStatusStopped PluginStatus = "stopped"
	PluginStatusRunning PluginStatus = "running"
	PluginStatusError   PluginStatus = "error"
)

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

func (m *Manager) StartPlugin(ctx context.Context, pluginID, command string, args []string, config map[string]interface{}) (*ManagedPlugin, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.plugins[pluginID]; exists {
		return nil, fmt.Errorf("plugin %s already exists", pluginID)
	}

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
	})

	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("failed to get RPC client: %w", err)
	}

	raw, err := rpcClient.Dispense("grpc_plugin")
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("failed to dispense plugin: %w", err)
	}

	pluginService, ok := raw.(shared.PluginService)
	if !ok {
		client.Kill()
		return nil, fmt.Errorf("failed to get plugin service")
	}

	info, err := pluginService.GetPluginInfo(ctx)
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("failed to get plugin info: %w", err)
	}

	routes, err := pluginService.GetRoutes(ctx)
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("failed to get plugin routes: %w", err)
	}

	health, err := pluginService.HealthCheck(ctx)
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("plugin health check failed: %w", err)
	}

	if health.Status != "healthy" {
		client.Kill()
		return nil, fmt.Errorf("plugin is not healthy: %s", health.Status)
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
	}

	m.plugins[pluginID] = managedPlugin

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
