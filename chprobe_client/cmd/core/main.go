package main

import (
	"os"

	"github.com/ricky97gr/chprobe/chprobe_client/conf"
	"github.com/ricky97gr/chprobe/chprobe_client/heartbeat"
	"github.com/ricky97gr/chprobe/chprobe_client/report"
	"github.com/ricky97gr/chprobe/chprobe_client/task"
	"github.com/ricky97gr/chprobe/chprobe_common/utils"
)

func main() {
	opt := utils.Options{
		FileName:   "",
		Level:      "info",
		ModuleName: "chprobe",
		W:          os.Stdout,
	}
	utils.Logger = utils.New(opt)
	utils.Logger.Infof("reporter client started\n")

	uuid := conf.GetUUID()
	if uuid == "" {
		utils.Logger.Infof("UUID not found, will register later\n")
		uuid = "00001"
	}
	task.SetClientUUID(uuid)

	task.StartTaskStream()
	task.WaitReady()
	utils.Logger.Infof("task stream ready, start heartbeat...")

	report.RegisterClient()
	go heartbeat.Run()

	select {}
}
