package grpc

import (
	"fmt"
	"io"
	"net"

	"github.com/ricky97gr/chprobe/chprobe_common/constant"
	"github.com/ricky97gr/chprobe/chprobe_common/proto"
	"github.com/ricky97gr/chprobe/chprobe_common/utils"
	"github.com/ricky97gr/chprobe/chprobe_server/handler"
	"github.com/ricky97gr/chprobe/chprobe_server/task"
	"google.golang.org/grpc"
)

type reportServer struct {
	proto.UnimplementedReporterServer
}

func (s *reportServer) ReportToServer(stream grpc.BidiStreamingServer[proto.MessageInfo, proto.ServerMessage]) error {
	var clientUUID string
	sendChan := make(chan *proto.ServerMessage, 100)
	defer func() {
		if clientUUID != "" {
			task.GlobalStreamManager.Unregister(clientUUID)
		}
		close(sendChan)
	}()

	go func() {
		for {
			taskMsg, ok := <-sendChan
			if !ok {
				return
			}
			if err := stream.Send(taskMsg); err != nil {
				utils.Logger.Errorf("send task to client %s failed: %v\n", clientUUID, err)
				return
			}
		}
	}()

	for {
		agentMsg, err := stream.Recv()
		if err == io.EOF {
			utils.Logger.Infof("client task stream closed by client\n")
			return nil
		}
		if err != nil {
			utils.Logger.Errorf("receive from client task stream error: %v\n", err)
			return err
		}

		if clientUUID == "" {
			clientUUID = agentMsg.Client
			task.GlobalStreamManager.Register(clientUUID, sendChan)
			utils.Logger.Infof("client %s task stream established\n", clientUUID)
		}

		switch agentMsg.MessageType {
		case int32(constant.MessageHeartbeat):
			_ = handler.HandleHeartbeat(agentMsg.Client)
		case int32(constant.MessageRegister):
			utils.Logger.Infof("receive register from client: %s\n", agentMsg.Client)
		}

		utils.Logger.Infof("receive from client %s, type=%d\n", clientUUID, agentMsg.MessageType)
	}
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
