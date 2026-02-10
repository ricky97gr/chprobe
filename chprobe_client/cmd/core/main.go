package main

import (
	"os"

	"github.com/ricky97gr/chprobe/chprobe_client/heartbeat"
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
	go heartbeat.Run()

	select {}
}
