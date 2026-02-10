package main

import (
	"os"

	"github.com/ricky97gr/chprobe/chprobe_common/utils"
	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/grpc"
	"github.com/ricky97gr/chprobe/chprobe_server/router"
	"github.com/ricky97gr/chprobe/chprobe_server/serverinfo"
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

	// 更新服务器信息
	if err := serverinfo.UpdateServerInfo(); err != nil {
		utils.Logger.Errorf("failed to update server info: %+v\n", err)
	}

	// 并行启动gRPC服务和web服务
	go grpc.Start()
	router.Start()

	select {}
}
