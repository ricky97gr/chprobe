package task

import (
	"sync"

	"github.com/ricky97gr/chprobe/chprobe_common/proto"
	"github.com/ricky97gr/chprobe/chprobe_common/utils"
)

type ClientStream struct {
	ClientUUID string
	SendChan   chan *proto.ServerMessage
}

type StreamManager struct {
	clients map[string]*ClientStream
	mu      sync.RWMutex
}

var GlobalStreamManager = &StreamManager{
	clients: make(map[string]*ClientStream),
}

func (m *StreamManager) Register(clientUUID string, sendChan chan *proto.ServerMessage) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[clientUUID] = &ClientStream{
		ClientUUID: clientUUID,
		SendChan:   sendChan,
	}
	utils.Logger.Infof("client %s task stream registered, total clients: %d\n", clientUUID, len(m.clients))
}

func (m *StreamManager) Unregister(clientUUID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if cs, exists := m.clients[clientUUID]; exists {
		close(cs.SendChan)
		delete(m.clients, clientUUID)
		utils.Logger.Infof("client %s task stream unregistered, total clients: %d\n", clientUUID, len(m.clients))
	}
}

func (m *StreamManager) SendTask(clientUUID string, taskMsg *proto.ServerMessage) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if cs, exists := m.clients[clientUUID]; exists {
		select {
		case cs.SendChan <- taskMsg:
			utils.Logger.Infof("task %s sent to client %s\n", taskMsg.TaskID, clientUUID)
			return true
		default:
			utils.Logger.Warnf("client %s send channel full, task %s dropped\n", clientUUID, taskMsg.TaskID)
			return false
		}
	}
	utils.Logger.Warnf("client %s not connected, cannot send task\n", clientUUID)
	return false
}

func (m *StreamManager) BroadcastTask(taskMsg *proto.ServerMessage) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	count := 0
	for _, cs := range m.clients {
		select {
		case cs.SendChan <- taskMsg:
			count++
		default:
		}
	}
	utils.Logger.Infof("task %s broadcast to %d clients\n", taskMsg.TaskID, count)
	return count
}

func (m *StreamManager) GetConnectedClients() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	clients := make([]string, 0, len(m.clients))
	for uuid := range m.clients {
		clients = append(clients, uuid)
	}
	return clients
}
