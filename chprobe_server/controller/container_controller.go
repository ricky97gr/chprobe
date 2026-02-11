package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
	"github.com/ricky97gr/chprobe/chprobe_server/pagination"
	"github.com/ricky97gr/chprobe/chprobe_server/response"
)

// ContainerUI 容器数据传输对象
type ContainerUI struct {
	ID             string `json:"id"`
	HostUUID       string `json:"hostUUID"`
	ContainerID    string `json:"containerID"`
	Name           string `json:"name"`
	Image          string `json:"image"`
	Command        string `json:"command"`
	State          string `json:"state"`
	Status         string `json:"status"`
	Ports          string `json:"ports"`
	Created        int64  `json:"created"`
	StartedAt      int64  `json:"startedAt"`
	FinishedAt     int64  `json:"finishedAt"`
	LastUpdateTime int64  `json:"lastUpdateTime"`
}

// 将ContainerInfo转换为ContainerUI
func toContainerUI(container models.ContainerInfo) ContainerUI {
	return ContainerUI{
		ID:             container.UUID,
		HostUUID:       container.HostUUID,
		ContainerID:    container.ContainerID,
		Name:           container.Name,
		Image:          container.Image,
		Command:        container.Command,
		State:          container.State,
		Status:         container.Status,
		Ports:          container.Ports,
		Created:        container.Created,
		StartedAt:      container.StartedAt,
		FinishedAt:     container.FinishedAt,
		LastUpdateTime: container.LastUpdateTime,
	}
}

// 将[]ContainerInfo转换为[]ContainerUI
func toContainerUIList(containers []models.ContainerInfo) []ContainerUI {
	containerUIs := make([]ContainerUI, len(containers))
	for i, container := range containers {
		containerUIs[i] = toContainerUI(container)
	}
	return containerUIs
}

// 查询容器列表
func GetContainerList(c *gin.Context) {
	// 获取分页查询参数
	pageQuery, err := pagination.GetPageQuery(c)
	if err != nil {
		response.Failed(c, response.ErrStruct, "Invalid page query parameters")
		return
	}

	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	// 查询容器总数
	var total int64
	if err := db.Model(&models.ContainerInfo{}).Count(&total).Error; err != nil {
		response.Failed(c, response.ErrDB, "Failed to count containers")
		return
	}

	// 分页查询容器列表
	var containers []models.ContainerInfo
	result := db.Scopes(pagination.ParseQuery(pageQuery)).Find(&containers)
	if result.Error != nil {
		response.Failed(c, response.ErrDB, "Failed to query containers")
		return
	}

	// 转换为DTO
	containerUIs := toContainerUIList(containers)

	// 返回响应
	response.Success(c, containerUIs, total)
}

// 查询单个容器详情
func GetContainerDetail(c *gin.Context) {
	// 获取容器ID
	id := c.Param("id")
	if id == "" {
		response.Failed(c, response.ErrStruct, "Invalid container ID")
		return
	}

	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	// 查询容器详情
	var container models.ContainerInfo
	result := db.Where("id = ?", id).First(&container)
	if result.Error != nil {
		response.Failed(c, response.ErrRecordNotFound, "Container not found")
		return
	}

	// 转换为DTO
	containerUI := toContainerUI(container)

	// 返回响应
	response.Success(c, containerUI, 1)
}
