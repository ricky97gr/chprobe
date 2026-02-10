package heartbeat

import (
	"encoding/json"
	"time"

	"github.com/ricky97gr/chprobe/chprobe_client/conf"
	"github.com/ricky97gr/chprobe/chprobe_client/report"
	"github.com/ricky97gr/chprobe/chprobe_common/constant"
	"github.com/ricky97gr/chprobe/chprobe_common/typed"
	"github.com/ricky97gr/chprobe/chprobe_common/utils"
)

const HeartbeatInterval = 60

func Run() {
	// 检查并注册UUID
	uuid := conf.GetUUID()
	if uuid == "" {
		utils.Logger.Infof("UUID not found, registering to server...\n")
		uuid = report.RegisterClient()
		if uuid != "" {
			conf.SetUUID(uuid)
			utils.Logger.Infof("Client registered successfully, UUID: %s\n", uuid)
		} else {
			utils.Logger.Errorf("Failed to register client, using default UUID\n")
			uuid = "00001"
		}
	}

	timer := time.NewTicker(HeartbeatInterval * time.Second)
	heartbeatInfo := typed.HeartBeatInfo{
		UUID:       uuid,
		ReportTime: time.Now().UnixMilli(),
	}
	for range timer.C {
		utils.Logger.Debugf("report heartbeat to server\n")
		heartbeatInfo.ReportTime = time.Now().UnixMilli()
		data, err := json.Marshal(heartbeatInfo)
		if err != nil {
			utils.Logger.Errorf("marshal heartbeat info to json failed, err: %v\n", err)
			continue
		}
		report.ReportMessageToServer(constant.MessageHeartbeat, data)
	}
}
