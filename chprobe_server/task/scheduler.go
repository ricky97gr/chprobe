package task

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ricky97gr/chprobe/chprobe_common/proto"
	"github.com/ricky97gr/chprobe/chprobe_common/typed"
	"github.com/ricky97gr/chprobe/chprobe_common/utils"
	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
)

const (
	TaskTypeHeartbeatCheck = 999
	TaskTypePluginExecute  = 1000
)

type TaskPluginBinder struct {
	TaskType int32
	PluginID string
}

var defaultBindings = []TaskPluginBinder{
	{
		TaskType: TaskTypePluginExecute,
		PluginID: "test",
	},
}

func StartTaskScheduler() {
	utils.Logger.Infof("task scheduler started, interval: 1 hour\n")
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				broadcastBoundTasks()
			}
		}
	}()
}

func findBestPluginForTaskType(taskType int32) (*models.Plugin, error) {
	var plugins []models.Plugin
	db, err := database.GetMysqlClient()
	if err != nil {
		return nil, err
	}
	if err := db.Where("status = ?", models.PluginStatusEnabled).Find(&plugins).Error; err != nil {
		return nil, err
	}

	for _, binding := range defaultBindings {
		if binding.TaskType == taskType {
			for _, p := range plugins {
				if p.PluginID == binding.PluginID {
					return &p, nil
				}
			}
		}
	}

	if len(plugins) > 0 {
		return &plugins[0], nil
	}
	return nil, fmt.Errorf("no suitable plugin found")
}

func broadcastBoundTasks() {
	for _, binding := range defaultBindings {
		plugin, err := findBestPluginForTaskType(binding.TaskType)
		if err != nil {
			utils.Logger.Errorf("find best plugin for task type %d failed: %v\n", binding.TaskType, err)
			continue
		}

		task := typed.Task{
			TaskID:    uuid.NewString(),
			TaskName:  fmt.Sprintf("Hourly Bound Task - Plugin: %s", plugin.Name),
			Timestamp: time.Now().UnixMilli(),
			PluginID:  plugin.PluginID,
		}

		data, err := json.Marshal(task)
		if err != nil {
			utils.Logger.Errorf("marshal bound task failed: %v\n", err)
			continue
		}

		serverMsg := &proto.ServerMessage{
			TaskID:    task.TaskID,
			TaskType:  binding.TaskType,
			Timestamp: task.Timestamp,
			Data:      data,
		}

		count := GlobalStreamManager.BroadcastTask(serverMsg)
		utils.Logger.Infof("bound task %s (taskType: %d, plugin: %s) dispatched to %d online agents\n", task.TaskID, binding.TaskType, plugin.PluginID, count)
	}
}

func DispatchBoundTaskTo(clientUUID string, taskType int32) bool {
	plugin, err := findBestPluginForTaskType(taskType)
	if err != nil {
		utils.Logger.Errorf("find best plugin for dispatch failed: %v\n", err)
		return false
	}

	task := typed.Task{
		TaskID:    fmt.Sprintf("%s-manual", uuid.NewString()),
		TaskName:  "Manual Bound Task",
		Timestamp: time.Now().UnixMilli(),
		PluginID:  plugin.PluginID,
	}

	data, err := json.Marshal(task)
	if err != nil {
		utils.Logger.Errorf("marshal manual bound task failed: %v\n", err)
		return false
	}

	serverMsg := &proto.ServerMessage{
		TaskID:    task.TaskID,
		TaskType:  taskType,
		Timestamp: task.Timestamp,
		Data:      data,
	}

	return GlobalStreamManager.SendTask(clientUUID, serverMsg)
}
