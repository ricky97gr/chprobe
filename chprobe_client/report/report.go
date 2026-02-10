package report

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ricky97gr/chprobe/chprobe_client/conf"
	"github.com/ricky97gr/chprobe/chprobe_common/constant"
	"github.com/ricky97gr/chprobe/chprobe_common/proto"
	"github.com/ricky97gr/chprobe/chprobe_common/typed"
	"github.com/ricky97gr/chprobe/chprobe_common/utils"
	"google.golang.org/grpc"
)

func ReportMessageToServer(messageType int, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	serverAddr := conf.GetServerAddr()
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		utils.Logger.Errorf("failed to get conn from endpoint[%s], err info: %+v\n", serverAddr, err)
		return err
	}
	defer conn.Close()
	client := proto.NewReporterClient(conn)
	info := &proto.MessageInfo{
		Client:      "chprobe",
		MessageType: int32(messageType),
		ReportTime:  time.Now().UnixMilli(),
		Data:        data,
	}
	_, err = client.ReportToServer(ctx, info)
	if err != nil {
		return err
	}
	return nil
}

func RegisterClient() string {
	// 收集客户端信息
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

	// 向服务器发送注册请求
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	serverAddr := conf.GetServerAddr()
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		utils.Logger.Errorf("failed to get conn from endpoint[%s], err info: %+v\n", serverAddr, err)
		return ""
	}
	defer conn.Close()
	client := proto.NewReporterClient(conn)
	info := &proto.MessageInfo{
		Client:      "chprobe",
		MessageType: int32(constant.MessageRegister),
		ReportTime:  time.Now().UnixMilli(),
		Data:        data,
	}

	// 接收服务器返回的UUID
	response, err := client.ReportToServer(ctx, info)
	if err != nil {
		utils.Logger.Errorf("register client failed, err: %v\n", err)
		return ""
	}

	if response != nil && response.Result != "" {
		return response.Result
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
