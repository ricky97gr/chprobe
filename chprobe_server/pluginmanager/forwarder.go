package pluginmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	shared "github.com/ricky97gr/chprobe/chprobe_common/plugin/shared"
)

type Forwarder struct {
	manager *Manager
}

func NewForwarder(manager *Manager) *Forwarder {
	return &Forwarder{
		manager: manager,
	}
}

func (f *Forwarder) ForwardJSON(pluginID string, path string, method string, body map[string]interface{}, headers map[string]string) ([]byte, error) {
	data := make(map[string]string)
	for k, v := range body {
		switch val := v.(type) {
		case string:
			data[k] = val
		default:
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal body value: %w", err)
			}
			data[k] = string(jsonBytes)
		}
	}

	result, err := f.manager.HandleRequest(context.Background(), pluginID, method, path, data)
	if err != nil {
		return nil, fmt.Errorf("failed to handle request: %w", err)
	}

	response := map[string]interface{}{
		"success": result.Success,
		"data":    result.Data,
		"error":   result.Error,
	}

	return json.Marshal(response)
}

func (f *Forwarder) GetPluginHealth(pluginID string) error {
	_, err := f.manager.GetPluginHealth(pluginID)
	return err
}

func (f *Forwarder) MatchAndForward(pluginID string, w http.ResponseWriter, r *http.Request) bool {
	plugin, exists := f.manager.GetPlugin(pluginID)
	if !exists {
		return false
	}

	for _, route := range plugin.Routes {
		if route.Path == r.URL.Path && (route.Method == "" || route.Method == r.Method) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusBadRequest)
				return true
			}

			var data map[string]interface{}
			if len(body) > 0 {
				if err := json.Unmarshal(body, &data); err != nil {
					http.Error(w, "Failed to parse request body", http.StatusBadRequest)
					return true
				}
			}

			headers := make(map[string]string)
			for k, v := range r.Header {
				if len(v) > 0 {
					headers[k] = v[0]
				}
			}

			result, err := f.manager.HandleRequest(context.Background(), pluginID, r.Method, r.URL.Path, make(map[string]string))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return true
			}

			response := map[string]interface{}{
				"success": result.Success,
				"data":    result.Data,
				"error":   result.Error,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return true
		}
	}

	return false
}

func (f *Forwarder) GetRoutesForPlugin(pluginID string) ([]*shared.RouteInfo, error) {
	return f.manager.GetPluginRoutes(pluginID)
}

func (f *Forwarder) GetAllRoutes() map[string][]*shared.RouteInfo {
	plugins := f.manager.GetAllPlugins()
	allRoutes := make(map[string][]*shared.RouteInfo)

	for _, plugin := range plugins {
		allRoutes[plugin.ID] = plugin.Routes
	}

	return allRoutes
}
