package report

import (
	"encoding/json"
	"time"

	"github.com/ricky97gr/chprobe/chprobe_client/conf"
	"github.com/ricky97gr/chprobe/chprobe_client/task"
	"github.com/ricky97gr/chprobe/chprobe_common/constant"
	"github.com/ricky97gr/chprobe/chprobe_common/proto"
	"github.com/ricky97gr/chprobe/chprobe_common/typed"
	"github.com/ricky97gr/chprobe/chprobe_common/utils"
)

func ReportMessageToServer(messageType int, data []byte) error {
	st := task.GetStream()
	if st == nil {
		utils.Logger.Warnf("task stream not ready, skip report\n")
		return nil
	}
	info := &proto.MessageInfo{
		Client:      conf.GetUUID(),
		MessageType: int32(messageType),
		ReportTime:  time.Now().UnixMilli(),
		Data:        data,
	}
	return st.Send(info)
}

func RegisterClient() string {
	clientInfo := typed.ClientRegisterInfo{
		ClientType:    "chprobe",
		Hostname:      getHostname(),
		IP:            getLocalIP(),
		MachineID:     getMachineID(),
		OsType:        utils.GetOsType(),
		Os:            utils.GetOs(),
		Arch:          utils.GetArch(),
		KernelVersion: utils.GetKernelVersion(),
		Version:       "0.0.1",
		Timestamp:     time.Now().UnixMilli(),
	}

	data, err := json.Marshal(clientInfo)
	if err != nil {
		utils.Logger.Errorf("marshal client register info failed, err: %v\n", err)
		return ""
	}

	st := task.GetStream()
	if st == nil {
		utils.Logger.Warnf("task stream not ready for register")
		return ""
	}

	info := &proto.MessageInfo{
		Client:      "chprobe",
		MessageType: int32(constant.MessageRegister),
		ReportTime:  time.Now().UnixMilli(),
		Data:        data,
	}
	err = st.Send(info)
	if err != nil {
		utils.Logger.Errorf("send register failed, err: %v\n", err)
		return ""
	}
	return ""
}

func getHostname() string {
	hostname, err := utils.GetHostname()
	if err != nil {
		utils.Logger.Errorf("get hostname failed, err: %v\n", err)
		return "unknown"
	}
	return hostname
}

func getLocalIP() string {
	ip, err := utils.GetLocalIP()
	if err != nil {
		utils.Logger.Errorf("get local ip failed, err: %v\n", err)
		return "unknown"
	}
	return ip
}

func getMachineID() string {
	machineID, err := utils.GetMachineID()
	if err != nil {
		utils.Logger.Errorf("get machine id failed, err: %v\n", err)
		return "unknown"
	}
	return machineID
}
