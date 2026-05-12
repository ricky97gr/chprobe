package task

import (
	"context"
	"encoding/json"
	"io"
	"sync"
	"time"

	"github.com/ricky97gr/chprobe/chprobe_client/conf"
	"github.com/ricky97gr/chprobe/chprobe_client/plugin"
	"github.com/ricky97gr/chprobe/chprobe_common/proto"
	"github.com/ricky97gr/chprobe/chprobe_common/typed"
	"github.com/ricky97gr/chprobe/chprobe_common/utils"
	"google.golang.org/grpc"
)

var (
	clientUUID string
	stream     proto.Reporter_ReportToServerClient
	streamMu   sync.RWMutex
	readyChan  = make(chan struct{})
)

func SetClientUUID(uuid string) {
	clientUUID = uuid
}

func GetStream() proto.Reporter_ReportToServerClient {
	streamMu.RLock()
	defer streamMu.RUnlock()
	return stream
}

func WaitReady() {
	<-readyChan
}

func markReady() {
	close(readyChan)
}

func StartTaskStream() {
	go func() {
		for {
			if err := runTaskStream(); err != nil {
				utils.Logger.Errorf("task stream error, reconnecting in 5s: %v\n", err)
				time.Sleep(5 * time.Second)
			}
		}
	}()
}

func runTaskStream() error {
	if clientUUID == "" {
		utils.Logger.Warnf("client uuid not set, wait for register...")
		time.Sleep(3 * time.Second)
		return nil
	}

	serverAddr := conf.GetServerAddr()
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		utils.Logger.Errorf("failed to dial server %s: %v\n", serverAddr, err)
		return err
	}
	defer conn.Close()

	client := proto.NewReporterClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	newStream, err := client.ReportToServer(ctx)
	if err != nil {
		return err
	}

	streamMu.Lock()
	stream = newStream
	streamMu.Unlock()
	utils.Logger.Infof("bidirectional task stream established with server, client: %s\n", clientUUID)
	markReady()

	go func() {
		for {
			taskMsg, err := newStream.Recv()
			if err == io.EOF {
				utils.Logger.Infof("server closed task stream\n")
				return
			}
			if err != nil {
				utils.Logger.Errorf("receive task from server error: %v\n", err)
				return
			}
			utils.Logger.Infof("received task from server, taskID: %s, type: %d\n", taskMsg.TaskID, taskMsg.TaskType)
			go handleTask(taskMsg)
		}
	}()

	<-newStream.Context().Done()
	streamMu.Lock()
	stream = nil
	streamMu.Unlock()
	return newStream.Context().Err()
}

func handleTask(taskMsg *proto.ServerMessage) {
	var task typed.Task
	if err := json.Unmarshal(taskMsg.Data, &task); err != nil {
		utils.Logger.Errorf("unmarshal task failed: %v\n", err)
		return
	}

	utils.Logger.Infof("processing task: %s, pluginID: %s\n", taskMsg.TaskID, task.PluginID)

	if err := plugin.EnsurePluginReady(&task); err != nil {
		utils.Logger.Errorf("ensure plugin %s ready failed: %v\n", task.PluginID, err)
		return
	}

	utils.Logger.Infof("plugin %s ready for execution\n", task.PluginID)
}
