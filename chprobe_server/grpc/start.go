package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/ricky97gr/chprobe/chprobe_common/constant"
	"github.com/ricky97gr/chprobe/chprobe_common/proto"
	"github.com/ricky97gr/chprobe/chprobe_common/utils"
	"github.com/ricky97gr/chprobe/chprobe_server/handler"
	"google.golang.org/grpc"
)

type reportServer struct {
	proto.UnimplementedReporterServer
}

func (s *reportServer) ReportToServer(ctx context.Context, msg *proto.MessageInfo) (*proto.Response, error) {
	utils.Logger.Infof("receive client %s, message: %+v", msg.Client, msg)
	switch msg.MessageType {
	case int32(constant.MessageHeartbeat):
		// 处理心跳消息
		err := handler.HandleHeartbeat(msg.Client)
		if err != nil {
			utils.Logger.Errorf("handle heartbeat failed, err: %v\n", err)
			return &proto.Response{
				Success: false,
				Result:  "",
			}, nil
		}
	case int32(constant.MessageRegister):
		// 处理注册消息
		uuid, err := handler.HandleRegister(msg.Data)
		if err != nil {
			utils.Logger.Errorf("handle register failed, err: %v\n", err)
			return &proto.Response{
				Success: false,
				Result:  "",
			}, nil
		}
		return &proto.Response{
			Success: true,
			Result:  uuid,
		}, nil
	default:
		utils.Logger.Warnf("unknown message type: %d", msg.MessageType)
	}
	return &proto.Response{
		Success: true,
	}, nil
}

func Start() {
	utils.Logger.Infof("reporter server started on port %d\n", 32000)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 32000))
	if err != nil {
		utils.Logger.Errorf("failed to listen %d port: %+v\n", 32000, err)
		return
	}
	s := grpc.NewServer()
	proto.RegisterReporterServer(s, &reportServer{})
	err = s.Serve(listener)
	if err != nil {
		utils.Logger.Errorf("failed to start grpc server, err: %+v\n", err)
		return
	}
}
