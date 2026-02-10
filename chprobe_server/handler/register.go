package handler

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/ricky97gr/chprobe/chprobe_common/typed"
	"github.com/ricky97gr/chprobe/chprobe_common/utils"
	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
)

func HandleRegister(data []byte) (string, error) {
	var registerInfo typed.ClientRegisterInfo
	err := json.Unmarshal(data, &registerInfo)
	if err != nil {
		utils.Logger.Errorf("unmarshal register info failed, err: %v\n", err)
		return "", err
	}

	utils.Logger.Infof("receive register request from client: %+v\n", registerInfo)

	// 根据machineID生成唯一标识
	clientUUID := generateUUID(registerInfo.MachineID)
	utils.Logger.Infof("generated UUID for client: %s\n", clientUUID)

	// 检查是否已存在该客户端
	db, err := database.GetMysqlClient()
	if err != nil {
		utils.Logger.Errorf("failed to get mysql client, err: %v\n", err)
		return "", err
	}
	var existingAgent models.Agent
	result := db.Where("machine_id = ?", registerInfo.MachineID).First(&existingAgent)

	if result.Error == nil {
		// 客户端已存在，更新所有字段
		utils.Logger.Infof("agent already exists with UUID: %s, updating all fields\n", existingAgent.UUID)
		existingAgent.HostName = registerInfo.Hostname
		existingAgent.IP = registerInfo.IP
		existingAgent.ClientType = registerInfo.ClientType
		existingAgent.OsType = registerInfo.OsType
		existingAgent.Os = registerInfo.Os
		existingAgent.Arch = registerInfo.Arch
		existingAgent.KernelVersion = registerInfo.KernelVersion
		existingAgent.Version = registerInfo.Version
		existingAgent.LastHeartTime = time.Now().UnixMilli()
		existingAgent.Status = "online"
		db.Save(&existingAgent)
		return existingAgent.UUID, nil
	}

	// 新客户端，创建记录
	agent := models.Agent{
		UUID:          clientUUID,
		HostName:      registerInfo.Hostname,
		IP:            registerInfo.IP,
		MachineID:     registerInfo.MachineID,
		ClientType:    registerInfo.ClientType,
		OsType:        registerInfo.OsType,
		Os:            registerInfo.Os,
		Arch:          registerInfo.Arch,
		KernelVersion: registerInfo.KernelVersion,
		Version:       registerInfo.Version,
		RegisterTime:  time.Now().UnixMilli(),
		LastHeartTime: time.Now().UnixMilli(),
		Status:        "online",
	}

	// 创建表（如果不存在）
	db.AutoMigrate(&models.Agent{})

	// 保存到数据库
	if err := db.Create(&agent).Error; err != nil {
		utils.Logger.Errorf("save agent info to database failed, err: %v\n", err)
		return "", err
	}

	utils.Logger.Infof("agent registered successfully, UUID: %s\n", clientUUID)
	return clientUUID, nil
}

func generateUUID(machineID string) string {
	// 基于machineID生成UUID，确保同一台机器生成的UUID相同
	if machineID == "" {
		// 如果没有machineID，生成随机UUID
		return uuid.New().String()
	}

	// 使用machineID作为命名空间生成UUID
	namespace := uuid.NewMD5(uuid.Nil, []byte(machineID))
	return namespace.String()
}
