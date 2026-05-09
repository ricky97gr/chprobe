package main

import (
	"os"

	"github.com/ricky97gr/chprobe/chprobe_common/utils"
	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/grpc"
	"github.com/ricky97gr/chprobe/chprobe_server/heartbeat"
	"github.com/ricky97gr/chprobe/chprobe_server/router"
	"github.com/ricky97gr/chprobe/chprobe_server/serverinfo"
	syslog "github.com/ricky97gr/chprobe/chprobe_server/utils"
)

func main() {
	opt := utils.Options{
		FileName:   "",
		Level:      "info",
		ModuleName: "chprobe",
		W:          os.Stdout,
	}
	utils.Logger = utils.New(opt)

	utils.Logger.Infof("chprobe Started\n")

	// 启动数据库
	database.Start()
	syslog.LogInfo(syslog.ModuleDatabase, "数据库连接成功，表结构迁移完成")

	// 检测版本变化并记录升级信息
	syslog.CheckVersionAndRecord()

	// 更新服务器信息
	if err := serverinfo.UpdateServerInfo(); err != nil {
		utils.Logger.Errorf("failed to update server info: %+v\n", err)
		syslog.LogError(syslog.ModuleSystem, "服务器信息更新失败: "+err.Error())
	} else {
		syslog.LogInfo(syslog.ModuleSystem, "服务器信息更新完成")
	}

	// 记录服务启动日志
	syslog.LogInfo(syslog.ModuleSystem, "ChProbe 服务启动成功")

	// 启动Agent心跳监控
	heartbeat.StartMonitor()
	syslog.LogInfo(syslog.ModuleSystem, "Agent心跳监控已启动，10分钟超时")

	// 并行启动gRPC服务和web服务
	go grpc.Start()
	router.Start()

	select {}
}
